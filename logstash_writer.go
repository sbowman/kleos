package kleos

import (
	"errors"
	"net"
	"time"
)

var (
	// ErrInvalidConnectionType returned when the connection isn't a TCPConn.  Shouldn't
	// typically happen.
	ErrInvalidConnectionType = errors.New("not a TCP connection")
)

// LogstashWriter is designed to output log messages to the Logstash TCP input.  Use this with
// ELKOutput to send log messages to an ELK-compatible stack.
type LogstashWriter struct {
	Host    string
	Timeout time.Duration

	conn   *net.TCPConn
}

// NewLogstashWriter creates a new writer to connect to the Logstash host and output log messages.
// Host should include hostname or IP address and port of the Logstash TCP service.
func NewLogstashWriter(host string, timeout time.Duration) *LogstashWriter {
	return &LogstashWriter{
		Host:    host,
		Timeout: timeout,
	}
}

func (w *LogstashWriter) Close() {
	_ = w.conn.Close()
}

// Dial connects to the Logstash TCP service.
func (w *LogstashWriter) Dial() error {
	addr, err := net.ResolveTCPAddr("tcp", w.Host)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		return err
	}

	var ok bool
	if w.conn, ok = conn.(*net.TCPConn); !ok {
		return ErrInvalidConnectionType
	}

	if err := w.conn.SetKeepAlive(true); err != nil {
		return err
	}

	return nil
}

// Write sends a block of bytes to the LogstashWriter buffer.  Once an array of bytes with a
// carriage return at the end is written, the buffer is shipped to Logstash via TCP and reset.
func (w *LogstashWriter) Write(b []byte) (int, error) {
	return w.conn.Write(b)
}
