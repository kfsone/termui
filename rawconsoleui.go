package termui

import (
	"bufio"
	"fmt"
	"io"
)

type RawConsoleUI struct {
	UserInterfaceBase
	source		io.Reader
}

// Close 'outputs' by calling Close()
// inputs will be closed by the ui when input ends
func (r *RawConsoleUI) Open() (inputs <-chan string) {
	inputChannel := r.UserInterfaceBase.Open(make(chan string, 255), 255)

	// dispatch messages to the output stream
	go r.writer(func(text string) {
		fmt.Println(text)
	})

	// consume messages from the input stream
	scanner := bufio.NewScanner(r.source)
	go r.reader(inputChannel, func(inputs chan<- string) (text string, ok bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	})

	return inputChannel
}

func NewRawConsoleUI(source io.Reader) *RawConsoleUI {
	return &RawConsoleUI{source: source}
}
