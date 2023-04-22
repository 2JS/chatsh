package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/pkg/term"
	terminal "golang.org/x/term"
)

func main() {
	pipePath := "/tmp/chatsh.pipe"

	// Create a named pipe.
	_ = syscall.Mkfifo(pipePath, 0644)

	argument := fmt.Sprintf(`CHATSH=1 script -q -F >(sed -u 's/\x1b\[[0-9;]*[a-zA-Z]//g' > %s)`, pipePath)

	cmd := exec.Command("zsh", "-c", argument)

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

	// Wait for the shell command to finish.
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

	println("chatsh exit.")
}
