//
// +build windows
//

package getpass

import (
	"errors"
	"io"
	"os"
	"syscall"
)

// getPassword implementation for Windows
func getPassword(prompt string) (pass string, err error) {
	msvcrt, _ := syscall.LoadLibrary("msvcrt.dll")
	defer syscall.FreeLibrary(msvcrt)
	_getch, _ := syscall.GetProcAddress(msvcrt, "_getch")

	// Write prompt
	os.Stdout.Write([]byte(prompt))

	// Read password from standard input
	buf := make([]byte, 0, 16)
	for {
		ch, _, err := syscall.Syscall(_getch, 0, 0, 0, 0)
		if err != 0 {
			return "", errors.New("_getch: " + err.Error())
		}

		if ch == '\r' || ch == '\n' {
			break
		} else if ch == '\x08' && len(buf) > 0 { // Backspace
			buf = buf[:len(buf)-1]
		} else if ch == '\x03' { // CLTR-C
			return "", io.EOF
		}
		buf = append(buf, byte(int(ch)))
	}

	// Write CRLF
	os.Stdout.Write([]byte("\r\n"))

	return string(buf), nil
}
