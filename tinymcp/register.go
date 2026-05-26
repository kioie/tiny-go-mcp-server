package tinymcp

import "fmt"

func registerRecover(label string, fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("register %s: %v", label, r)
		}
	}()
	fn()
	return nil
}
