package main

import (
	"bufio"
	"os"
	"syscall"
	"time"
)

type Pipe struct {
	pipePath string
}

func NewPipe(pipePath string) *Pipe {
	_ = syscall.Mkfifo(pipePath, 0600)

	return &Pipe{pipePath}
}

func (pipe *Pipe) ReadChannel() chan string {
	channel := make(chan string)

	go func() {
		for {
			p, err := os.Open(pipe.pipePath)
			if err != nil {
				panic(err)
			}
			defer p.Close()

			// Create a buffered reader.
			reader := bufio.NewReader(p)

			// Continuously read lines from the pipe.
			for {
				p.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

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

func (pipe *Pipe) WriteOnce(line string) {
	// Open the named pipe for writing.
	p, err := os.OpenFile(pipe.pipePath, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	// Write to the pipe.
	p.WriteString(line)
}

func (pipe *Pipe) WriteChannel() chan string {
	channel := make(chan string)

	go func() {
		p, err := os.OpenFile(pipe.pipePath, os.O_WRONLY, os.ModeNamedPipe)
		if err != nil {
			panic(err)
		}
		defer p.Close()

		for line := range channel {
			p.WriteString(line)
		}
	}()

	return channel
}

type BidirectionalPipe struct {
	reader *Pipe
	writer *Pipe
}

func NewBidirectionalPipe(readPipePath string, writePipePath string) *BidirectionalPipe {
	return &BidirectionalPipe{
		reader: NewPipe(readPipePath),
		writer: NewPipe(writePipePath),
	}
}

func (pipe *BidirectionalPipe) ReadChannel() chan string {
	return pipe.reader.ReadChannel()
}

func (pipe *BidirectionalPipe) WriteChannel() chan string {
	return pipe.writer.WriteChannel()
}

func (pipe *BidirectionalPipe) WriteOnce(line string) {
	pipe.writer.WriteOnce(line)
}
