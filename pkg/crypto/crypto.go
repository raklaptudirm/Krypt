// Copyright © 2021 Rak Laptudirm <raklaptudirm@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crypto

import (
	"fmt"
	"sync"
	"time"

	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"math/rand"

	"github.com/raklaptudirm/krypt/pkg/dir"
	"golang.org/x/crypto/pbkdf2"
)

// EncryptWithKey is a wrapper on Encrypt which automatically reads the
// key from the key root file.
func EncryptWithKey(src []byte) (enc []byte, err error) {
	key, err := dir.Key()
	if err != nil {
		return
	}

	enc, err = Encrypt(src, key)
	return
}

// DecryptWithKey is a wrapper on Decrypt which automatically reads the
// key from the key root file.
func DecryptWithKey(ct []byte) (clt []byte, err error) {
	key, err := dir.Key()
	if err != nil {
		return
	}

	clt, err = Decrypt(ct, key)
	return
}

// Sha256 wraps the sha256.Sum256 method to return a []byte instead of
// a [32]byte array.
func Sha256(data []byte) []byte {
	checksum := sha256.Sum256(data)
	return checksum[:]
}

// Encrypt encrypts src with key using the AES encryption algorithm. It
// automatically generates a random salt and appends to the front of the
// ciphertext.
func Encrypt(src []byte, key []byte) (enc []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return enc, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := RandBytes(aesgcm.NonceSize())    // random iv
	enc = aesgcm.Seal(nonce, nonce, src, nil) // append to iv
	return
}

var ErrNoNonce = fmt.Errorf("ciphertext smaller than nonce")

// Decrypt decrypts ct with key with the AES encryption algorithm. It
// automatically extracts the salt from the ct.
func Decrypt(ct []byte, key []byte) (clt []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return clt, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonceSize := aesgcm.NonceSize()
	if len(ct) < nonceSize {
		// ciphertext can't be smaller than iv
		err = ErrNoNonce
		return
	}

	nonce, ct := ct[:nonceSize], ct[nonceSize:] // extract iv
	clt, err = aesgcm.Open(nil, nonce, ct, nil)
	return
}

// Pbkdf2 generates an AES algorithm key from pw, using the SHA-256 hash
// algorithm and the provided salt.
func Pbkdf2(pw []byte, salt []byte) (key []byte) {
	iter := 4096 // no of pbkdf2 iterations
	klen := 32   // length of key in bytes

	algo := sha256.New // hash algorithm for pbkdf

	key = pbkdf2.Key(pw, salt, iter, klen, algo)
	return
}

var setSeed sync.Once

func RandBytes(len int) []byte {
	// set rand seed once
	setSeed.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})

	// generate len random bytes
	b := make([]byte, len)
	rand.Read(b)
	return b
}
