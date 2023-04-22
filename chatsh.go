package main

import (
	"fmt"
	"os/exec"
	"syscall"
)

func main() {
	pipePath := "/tmp/chatsh.pipe"

	// Create a named pipe.
	_ = syscall.Mkfifo(pipePath, 0644)

	argument := fmt.Sprintf(`CHATSH=1 script -q -F >(sed -u 's/\x1b\[[0-9;]*[a-zA-Z]//g' > %s)`, pipePath)

	cmd := exec.Command("zsh", "-c", argument)

	terminal := NewTerminal(cmd)
	defer terminal.Close()

	// Wait for the shell command to finish.
	err := cmd.Wait()
	if err != nil {
		panic(err)
	}

	println("chatsh exit.")
}
