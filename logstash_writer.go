package kleos

import (
	"errors"
	"net"
	"time"
)

var (
	ErrInvalidConnectionType = errors.New("not a TCP connection")
	ErrWriteFailed           = errors.New("failed to write entire log message")
)

// LogstashWriter is designed to output log messages to the Logstash TCP input.
type LogstashWriter struct {
	Host    string
	Timeout time.Duration

	conn *net.TCPConn
	ch   chan []byte
}

// NewLogstashWriter creates a new writer to connect to the Logstash host and output log messages.
// Host should include hostname or IP address and port of the Logstash TCP service.
func NewLogstashWriter(host string, timeout time.Duration, bufferSize int) *LogstashWriter {
	w := &LogstashWriter{
		Host:    host,
		Timeout: timeout,
	}

	if bufferSize > 0 {
		w.ch = make(chan []byte, bufferSize)
	} else {
		w.ch = make(chan []byte)
	}

	go w.sender()
	return w
}

func (w *LogstashWriter) Close() {
	if w.ch != nil {
		close(w.ch)
		w.ch = nil
	} else if w.conn != nil {
		_ = w.conn.Close()
		w.conn = nil
	}
}

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

func (w *LogstashWriter) Write(b []byte) error {
	w.ch <- b
	return nil
}

func (w *LogstashWriter) sender() {
	for b := range w.ch {
		_ = w.send(b)
	}
}

func (w *LogstashWriter) send(b []byte) error {
	if err := w.conn.SetDeadline(time.Now().Add(w.Timeout)); err != nil {
		return err
	}

	n, err := w.conn.Write(b)
	if err != nil {
		return err
	}

	if n != len(b) {
		return ErrWriteFailed
	}

	return nil
}
