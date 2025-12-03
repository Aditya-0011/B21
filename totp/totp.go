package totp

import (
	"fmt"

	"github.com/pquerna/otp/totp"
)

type Auth struct {
	OTP    string
	Secret string
}

func (t *Auth) Validate() (bool, error) {
	if !totp.Validate(t.OTP, t.Secret) {
		return false, fmt.Errorf("Totp not valid")
	}

	return true, nil
}
