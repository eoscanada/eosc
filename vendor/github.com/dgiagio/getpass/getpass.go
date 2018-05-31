// Package getpass provides support for reading passwords on the command-line.
package getpass

// GetPassword displays a prompt to the user and reads the password from
// standard input without echoing on the screen.
func GetPassword(prompt string) (pass string, err error) {
	return getPassword(prompt)
}
