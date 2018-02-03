package service

// Config provides configuration for the tldrfeed service
type Config struct {
	// Port to bind to
	Port int
	// Should JSON be indented nicely
	IndentJSON bool
}
