package tinymcp

import (
	"fmt"
)

func registerRecover(label string, fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: register %s: %v", ErrRegistrationFailed, label, r)
		}
	}()
	fn()
	return nil
}
