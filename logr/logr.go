package logr

// Logr declares the minimal logging interface used by loggregator clients.
type Logr interface {
	Printf(string, ...interface{})
	Panicf(string, ...interface{})
}
