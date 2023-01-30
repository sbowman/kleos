package kleos

// Logger represents a simple logger interface commonly used by third-party packages.
// type Logger interface {
// 	Printf(msg string, args ...interface{})
// }
//
// type logger struct {
// 	verbosity int
// }
//
// // Printf logs a message to Kleos logger.
// func (l logger) Printf(msg string, args ...interface{}) {
// 	if verbosity < l.verbosity {
// 		return
// 	}
//
// 	m := Message{
// 		when:  time.Now(),
// 		debug: true,
// 	}
//
// 	m.FileAndLineNumberBack(3).Output(msg, args...)
// }
//
// // NewLogger creates a new logger to use with simple logging interfaces.  Logs to the debug
// // verbosity.
// func NewLogger(verbosity int) Logger {
// 	return logger{verbosity: verbosity}
// }
