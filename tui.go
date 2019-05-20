package termui

import (
	"github.com/marcusolsson/tui-go"
	"log"
)

type TUIUserInterface struct {
	UserInterfaceBase

	Tui     tui.UI
	entry   *tui.Entry
	display *tui.Box
}

func (u *TUIUserInterface) Close() {
	defer u.Tui.Quit()
	u.UserInterfaceBase.Close()
}

func (u *TUIUserInterface) Open() (inputs <-chan string) {
	inputChannel := u.UserInterfaceBase.Open(make(chan string, 255), 4)

	u.entry.SetFocused(true)
	u.Tui.SetKeybinding("Esc", func() { inputChannel <- "quit" })
	u.entry.OnSubmit(func(e *tui.Entry) {
		inputChannel <- e.Text()
		e.SetText("")
	})

	// Start the UI in its own go thread.
	go func() {
		defer close(inputChannel)
		err := u.Tui.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// start the handler for writing to the ui
	go u.writer(func(text string) {
		u.display.Append(tui.NewHBox(tui.NewLabel(text)))
	})

	// there's no reader counter part, it's event based.
	return inputChannel
}

func (u *TUIUserInterface) Write(output []byte) (written int, err error) {
	u.display.Append(tui.NewHBox(tui.NewLabel(string(output))))
	return len(output), nil
}

func (u *TUIUserInterface) WriteString(text string) (written int, err error){
	return u.Write([]byte(text))
}

func NewTUIUserInterface() *TUIUserInterface {
	display := tui.NewVBox()
	scrollArea := tui.NewScrollArea(display)
	scrollArea.SetAutoscrollToBottom(true)
	scrollBackBox := tui.NewVBox(scrollArea)
	scrollBackBox.SetBorder(true)

	entry := tui.NewEntry()
	entry.SetSizePolicy(tui.Expanding, tui.Maximum)

	entryBox := tui.NewHBox(entry)
	entryBox.SetBorder(true)
	entryBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	window := tui.NewVBox(scrollBackBox, entryBox)
	window.SetBorder(false)

	screen := tui.NewHBox()
	screen.Append(window)

	t, err := tui.New(screen)
	if err != nil {
		log.Fatal(err)
	}

	return &TUIUserInterface{Tui: t, entry: entry, display: display}
}

