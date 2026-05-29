package tinymcp

import "errors"

// Sentinel errors for programmatic handling at registration time.
var (
	ErrNilServer          = errors.New("tinymcp: nil server")
	ErrNilTool            = errors.New("tinymcp: nil tool")
	ErrNilHandler         = errors.New("tinymcp: nil handler")
	ErrRegistrationFailed = errors.New("tinymcp: registration failed")
)
