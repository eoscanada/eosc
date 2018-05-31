//
// +build !windows
//

package getpass

// #include <termios.h>
// #include <unistd.h>
// #include <stdlib.h>
import "C"

import (
	"errors"
	"io"
	"os"
	"unsafe"
)

const stdinFd = C.STDIN_FILENO
const stdoutFd = C.STDOUT_FILENO

// getPassword implementation for POSIX systems using termios API
func getPassword(prompt string) (pass string, err error) {
	var oldt C.struct_termios
	n, err := C.tcgetattr(stdinFd, &oldt)
	if n == -1 {
		return "", errors.New("tcgetattr: " + err.Error())
	}

	// Disable echo
	newt := oldt
	newt.c_lflag &^= C.ECHO
	newt.c_lflag |= C.ICANON | C.ISIG
	newt.c_iflag |= C.ICRNL
	n, err = C.tcsetattr(stdinFd, C.TCSAFLUSH, &newt)
	if n == -1 {
		return "", errors.New("tcsetattr: " + err.Error())
	}

	defer func() {
		C.tcsetattr(stdinFd, C.TCSANOW, &oldt)
	}()

	// Write prompt
	os.Stdout.Write([]byte(prompt))

	// Read password from standard input
	var ch byte
	buf := make([]byte, 0, 16)
	for {
		n, err := C.read(stdinFd, unsafe.Pointer(&ch), 1)
		if n == -1 {
			return "", errors.New("read: " + err.Error())
		} else if n == 0 {
			if len(buf) == 0 {
				return "", io.EOF
			}
			break
		}
		if ch == '\n' {
			break
		}
		buf = append(buf, ch)
	}

	// Write LF
	os.Stdout.Write([]byte("\n"))

	return string(buf), nil
}
