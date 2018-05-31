getpass
=======

Go library that provides support for reading passwords on the command-line.

* Simple, easy to use
* Supports POSIX systems via ```termios``` API and Windows via ```_getch``` API

# Installation
```bash
$ go get github.com/dgiagio/getpass
```

# Example usage
```go
import (
	"fmt"
	"github.com/dgiagio/getpass"
)

pass, _ := getpass.GetPassword("Password: ")
fmt.Printf("Entered password: %s\n", pass)
```
