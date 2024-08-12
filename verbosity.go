package kleos

var verbosity uint8

// SetVerbosity sets the verbosity level of the debug logging.  Zero disable debug logging.
func SetVerbosity(level uint8) {
	local.SetVerbosity(level)
}

// SetVerbosity sets the verbosity level of the debug logging.  Zero disable debug logging.
func (k *Kleos) SetVerbosity(level uint8) {
	k.Lock()
	defer k.Unlock()

	verbosity = level
}

// Verbosity represents a message's verbosity level, starting at level 0 (lowest detail,
// always logged) to level 4 (highest detail, rarely logged).  Verbosity may be adjusting
// during runtime.
//
// Recommendations:
//
//   - 0 - disabled; not logged
//   - 1 - startup info and error messages
//   - 2 - basic debug information with a minimal amount of detail, e.g. "User saved"
//   - 3 - specific low-level details, e.g. SQL queries, small JSON objects, short arrays
//   - 4 - Relentlessly specific details, such as incoming and outgoing large JSON
//     documents, full HTTP request body, etc.
func Verbosity() uint8 {
	return local.Verbosity()
}

// Verbosity represents a message's verbosity level, starting at level 0 (lowest detail,
// always logged) to level 4 (highest detail, rarely logged).  Verbosity may be adjusting
// during runtime.
//
// Recommendations:
//
//   - 0 - disabled; not logged
//   - 1 - startup info and error messages
//   - 2 - basic debug information with a minimal amount of detail, e.g. "User saved"
//   - 3 - specific low-level details, e.g. SQL queries, small JSON objects, short arrays
//   - 4 - Relentlessly specific details, such as incoming and outgoing large JSON
//     documents, full HTTP request body, etc.
func (k *Kleos) Verbosity() uint8 {
	k.RLock()
	defer k.RUnlock()

	return verbosity
}
