package totp

import (
	"fmt"
	"os"

	"github.com/pquerna/otp/totp"
)

type Auth struct {
	OTP string
}

func (t *Auth) Validate() (bool, error) {
	secret := os.Getenv("TOTP_SECRET")
	if secret == "" {
		return false, fmt.Errorf("Secret not found")
	}

	if !totp.Validate(t.OTP, secret) {
		return false, fmt.Errorf("Totp not valid")
	}

	return true, nil
}
