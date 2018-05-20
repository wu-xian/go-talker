package terminal

import (
	"learn/src/go-talker/log"

	"github.com/cjbassi/termui"
)

type InputBox struct {
	termui.Block
	Text        string
	TextLine    int
	TextFgColor int
	TextBgColor int
	WrapLength  int
}

var (
	// Word is word
	Word []string
)

func init() {
	_0, _9, a, z := int('0'), int('9'), int('a'), int('z')
	for i := a; i <= z; i++ {
		Word = append(Word, string(i))
	}

	for i := _0; i <= _9; i++ {
		Word = append(Word, string(i))
	}
}

// NewInputBox return a new input box pointer
func NewInputBox() *InputBox {
	return &InputBox{}
}

// Buffer return current input box buffer
func (self *InputBox) Buffer() *termui.Buffer {
	buf := self.Block.Buffer()
	width := self.X - 4
	y := self.Y - 2
	if width <= 0 {
		return buf
	}
	if len(self.Text) > width {
		texts := make([]string, 0)
		lost := self.Text
		for ; len(lost) > width; lost = lost[width:] {
			texts = append(texts, lost[:width+1])
		}
		texts = append(texts, lost)

		for i, v := range texts {
			if i < y {
				buf.SetString(2, 2+i, v, 35, 47)
			}
		}
	} else {
		buf.SetString(2, 2, self.Text, 35, 47)
	}
	return buf
}

func (self *InputBox) ListenInput(message chan string) {
	termui.On(Word, func(e termui.Event) {
		self.keyDown(e.Key)
	})

	termui.On("<space>", func(e termui.Event) {
		self.keyDown(" ")
	})

	termui.On("<enter>", func(e termui.Event) {
		if self.Text != "" {
			log.Logger.Info("enter")
			message <- self.Text
			self.Text = ""
			termui.Render(termui.Body)
		}
	})
}

func (self *InputBox) keyDown(key string) {
	self.Text += key
	termui.Render(termui.Body)
}
