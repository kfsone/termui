package termui

import (
	"errors"
	"sync"
)

type UserInterfacer interface {
	Open() (inputs <-chan string)
	Close()
	Write(text []byte) (count int, err error)
	WriteString(text string) (count int, err error)
}

type UserInterfaceBase struct {
	outputs	*chan string
	wg		sync.WaitGroup
}

func (u *UserInterfaceBase) Close() {
	if u.outputs != nil {
		close(*u.outputs)
		u.outputs = nil
	}
	u.wg.Wait()
}

func (u *UserInterfaceBase) writer(callback func (string)) {
	u.wg.Add(1)
	defer u.wg.Done()
	for text := range *u.outputs {
		callback(text)
	}
}

func (u *UserInterfaceBase) reader(inputs chan<- string, callback func(inputs chan<- string) (text string, ok bool)) {
	defer close(inputs)
	for {
		text, ok := callback(inputs)
		if ok == false {
			break
		}
		inputs <- text
	}
}

func (u *UserInterfaceBase) Open(outputs chan string, queueDepth int) (inputs chan string) {
	u.outputs = &outputs
	return make(chan string, queueDepth)
}

func (u UserInterfaceBase) Write(text []byte) (count int, err error) {
	return u.WriteString(string(text))
}

func (u UserInterfaceBase) WriteString(text string) (count int, err error) {
	if u.outputs == nil {
		return 0, errors.New("UI is not open for business")
	}
	*u.outputs <- text
	return len(text), nil
}