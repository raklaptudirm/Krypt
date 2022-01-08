package build

import (
	"github.com/raklaptudirm/krypt/internal/auth"
	"github.com/raklaptudirm/krypt/pkg/pass"
)

// values of the managers can be manually changed
// to make krypt work in any environment.
var PassManager pass.Manager
var AuthManager auth.Manager
