// Package dexec provides an interface similar to os/exec to run external
// commands inside containers.
//
// Please read documentation carefully about semantic differences between
// os/exec and dexec.
//
// Use Case
//
// This utility is intended to provide an execution environment without
// changing the existing code a lot to execute a command on a remote machine
// (or a pool of machines) running Docker engine or locally to limit resource
// usage and have extra security.
//
// Dependencies
//
// The package needs the following dependencies to work:
//  go get github.com/fsouza/go-dockerclient
//
// Known issues
//
// - You may receive empty stdout/stderr from commands if the executed command
// does not end with a trailing new line or has a different flushing behavior.
//
// - StdinPipe/StdoutPipe/StderrPipe should be used with goroutines (one that
// executes Cmd.Wait() and another one that writes to/reads from the pipe.
// Otherwise, the code may hang 10% of the time due to some timing issue.
package dexec
