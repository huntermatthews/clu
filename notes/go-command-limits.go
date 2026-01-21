package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

// LimitedWriter is an io.Writer that only writes up to a maximum number of bytes.
type LimitedWriter struct {
	W            io.Writer
	N            int
	totalWritten int
}

func (l *LimitedWriter) Write(p []byte) (n int, err error) {
	if l.totalWritten >= l.N {
		return 0, io.EOF // Stop writing once the limit is reached
	}

	remaining := l.N - l.totalWritten
	if len(p) > remaining {
		p = p[:remaining]
	}

	n, err = l.W.Write(p)
	l.totalWritten += n
	if l.totalWritten >= l.N && err == nil {
		err = io.EOF // Indicate that the limit has been reached
	}
	return n, err
}

func main() {
	cmd := exec.Command("your_command_that_spews_output")

	var buf bytes.Buffer
	limitedWriter := &LimitedWriter{W: &buf, N: 100} // Limit to 100 bytes

	cmd.Stdout = limitedWriter
	cmd.Stderr = limitedWriter // Can use the same writer for both, if desired

	err := cmd.Run()
	if err != nil {
		// The command might fail here with an EOF error if it tries to write past the limit
		fmt.Printf("Command finished with error: %v\n", err)
	}

	fmt.Printf("Limited output (first 100 bytes):\n%s\n", buf.String())
}
