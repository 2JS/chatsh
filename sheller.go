package main

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/pkg/term"
)

func main() {
	shell := "/bin/zsh"

	cmd := exec.Command(shell)

	// Create a new PTY.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	defer ptmx.Close()

	// Set the terminal to raw mode.
	t, err := term.Open("/dev/tty")
	if err != nil {
		panic(err)
	}
	defer t.Close()

	err = t.SetRaw()
	if err != nil {
		panic(err)
	}
	defer t.Restore()

	// Handle SIGHUP signal.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP)
	go func() {
		for range signalChan {
			// Do nothing, just allow the shell to exit gracefully.
		}
	}()

	// Copy data from the PTY to stdout.
	go func() { _, _ = io.Copy(os.Stdout, ptmx) }()

	// Copy data from stdin to the PTY.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()

	// Wait for the shell command to finish.
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

	println("Shell command completed successfully.")
}
