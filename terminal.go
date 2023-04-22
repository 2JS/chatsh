package main

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/pkg/term"
	terminal "golang.org/x/term"
)

type Terminal struct {
	ptmx *os.File
	term *term.Term
}

func NewTerminal(cmd *exec.Cmd) *Terminal {
	// Create a new PTY.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}

	// Set the terminal to raw mode.
	term, err := term.Open("/dev/tty")
	if err != nil {
		panic(err)
	}

	err = term.SetRaw()
	if err != nil {
		panic(err)
	}

	// Function to resize the PTY based on the current window size.
	resizePty := func() {
		width, height, err := terminal.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			return
		}
		pty.Setsize(ptmx, &pty.Winsize{Cols: uint16(width), Rows: uint16(height)})
	}

	// Initial resize.
	resizePty()

	// Handle window size changes.
	sigwinch := make(chan os.Signal, 1)
	signal.Notify(sigwinch, syscall.SIGWINCH)
	go func() {
		for range sigwinch {
			resizePty()
		}
	}()

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

	return &Terminal{ptmx, term}
}

func (t *Terminal) Close() {
	t.term.Restore()
	t.term.Close()
	t.ptmx.Close()
}
