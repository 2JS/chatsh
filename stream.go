package main

import (
	"bufio"
	"os"
	"strings"
	"syscall"
	"time"
)

func streamCmd() chan string {
	return stream("/tmp/chatsh/cmd.pipe")
}

func streamIO() chan string {
	pipePath := "/tmp/chatsh/io.pipe"

	_ = syscall.Mkfifo(pipePath, 0644)

	channel := make(chan string)

	go func() {
		for {
			// Open the named pipe for reading.
			pipe, err := os.Open(pipePath)
			if err != nil {
				panic(err)
			}
			defer pipe.Close()

			// Create a buffered reader.
			reader := bufio.NewReader(pipe)

			// Continuously read lines from the pipe.
			for {
				pipe.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

				line, err := reader.ReadString('\n')
				if err != nil {
					break
				}

				channel <- line
			}
		}
	}()

	return channel
}

func readIO() string {
	pipePath := "/tmp/chatsh/io.pipe"

	// Open the named pipe for reading.
	pipe, err := os.Open(pipePath)
	if err != nil {
		panic(err)
	}
	defer pipe.Close()

	// Create a buffered reader.
	reader := bufio.NewReader(pipe)

	var builder strings.Builder

	// Continuously read lines from the pipe.
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		builder.WriteString(line)
	}

	return builder.String()
}

func stream(pipePath string) chan string {
	_ = syscall.Mkfifo(pipePath, 0644)

	channel := make(chan string)

	go func() {
		for {
			// Open the named pipe for reading.
			pipe, err := os.Open(pipePath)
			if err != nil {
				panic(err)
			}
			defer pipe.Close()

			// Create a buffered reader.
			reader := bufio.NewReader(pipe)

			var builder strings.Builder

			// Continuously read lines from the pipe.
			for {
				pipe.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

				line, err := reader.ReadString('\n')
				if err != nil {
					break
				}

				builder.WriteString(line)
			}

			line := builder.String()

			channel <- line
		}
	}()

	return channel
}
