package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	apiKey := os.Getenv("CHATSH_OPENAI_API_KEY")
	nested := os.Getenv("CHATSH")
	if nested == "1" {
		println("nested chatsh is not supported yet.")
		return
	}

	system := `You are shell assistant of a shell user. Given the shell interaction history, you respond helpfully. The user may or may not ask things relavant to the shell history. Shell history will be given after 'Shell:' keward, and the user request will be given after 'User:' keyword. Answer useful, clear, and super consize. Don't end with rhetorical words such as "Is there anything else you need help with?". Answer simply only 'Yes' to this first prompt.`

	pipePath := "/tmp/chatsh/io.pipe"

	// Create a named pipe.
	pipe := NewPipe(pipePath)

	argument := fmt.Sprintf(`CHATSH=1 script -q -F >(sed -u 's/\x1b\[[0-9;]*[a-zA-Z]//g' > %s)`, pipePath)

	cmd := exec.Command("zsh", "-c", argument)

	terminal := NewTerminal(cmd)
	defer terminal.Close()

	client := NewClient(apiKey)

	builder := new(strings.Builder)

	go func() {
		for io := range pipe.ReadChannel() {
			builder.WriteString(io)
		}
	}()

	go func() {
		channel, _ := client.stream(system)
		for token := range channel {
			fmt.Fprint(os.Stderr, strings.ReplaceAll(token, "\n", "\r\n"))
		}

		cmdPipe := NewBidirectionalPipe(
			"/tmp/chatsh/cmd/reader.pipe",
			"/tmp/chatsh/cmd/writer.pipe",
		)

		commands := cmdPipe.ReadChannel()

		for command := range commands {
			answers := cmdPipe.WriteChannel()

			_ = fmt.Sprint(command)
			shellio := builder.String()
			prompt := fmt.Sprintf("Shell:\r\n%s\r\n\r\nUser: %s", shellio, command)
			builder.Reset()

			// fmt.Fprintln(os.Stderr, prompt)

			response, err := client.stream(prompt)
			if err != nil {
				panic(err)
			}

			for token := range response {
				answers <- strings.ReplaceAll(token, "\n", "\r\n")
			}
			answers <- "\r\n"
			close(answers)
		}
	}()

	// Wait for the shell command to finish.
	err := cmd.Wait()
	if err != nil {
		panic(err)
	}

	println("chatsh exit.")
}
