package terminal

import (
	"fmt"
	"learn/src/go-talker/log"
	"sync"

	"github.com/cjbassi/termui"
)

type MessageBox struct {
	*termui.Block
	Messages    []Message
	Boxes       []Box
	MessageChan chan *Message
	InChan      chan *Message
	locker      sync.Mutex
	showIndex   int32
	lastIndex   int32
}

func (self *MessageBox) AddMessage(message Message) {
	//self.locker.Lock()
	log.Logger.Info("add message !!!!!!!")
	self.Messages = append(self.Messages, message)
	//self.locker.Unlock()
	self.InChan <- &message
	log.Logger.Info("add message !!!!!!! over")
}

func NewMessageBox() *MessageBox {
	return &MessageBox{
		Block:       termui.NewBlock(),
		Messages:    []Message{},
		Boxes:       make([]Box, 0),
		MessageChan: make(chan *Message, 0),
		InChan:      make(chan *Message, 0),
	}
}

type Message struct {
	Content string //消息内容
	Name    string //消息发送者名称
	Time    int64  //发送消息时间
}

type Box struct {
	message *Message
	index   int32
}

func NewMessageList() *MessageBox {
	return &MessageBox{
		Block: termui.NewBlock(),
	}
}

func (self *MessageBox) Buffer() *termui.Buffer {
	log.Logger.Info("bbbbbbbbbbbbb", len(self.Messages))
	self.locker.Lock()
	buf := self.Block.Buffer()
	width := self.X - 4
	y := self.Block.Y
	y2 := y/2 - 1
	if len(self.Messages) > y2 {
		i := len(self.Messages) - y2
		shown := self.Messages[i:]
		for _i, _v := range shown {
			fmtMessage := messageFormatter(_v.Name, _v.Content)
			if len(fmtMessage) > width {
				firstLine := fmtMessage[:width+1]
				secondLine := fmtMessage[width+1:]
				if len(secondLine) > width-6 {
					secondLine = secondLine[:width-6] + "…"
				}
				buf.SetString(2, 1+2*_i, firstLine, 35, 47)
				buf.SetString(4, 2+2*_i, secondLine, 35, 47)
			} else {
				buf.SetString(2, 1+2*_i, fmtMessage, 35, 47)
			}
		}
	} else {
		for _i, _v := range self.Messages {
			fmtMessage := messageFormatter(_v.Name, _v.Content)
			if len(fmtMessage) > width {
				firstLine := fmtMessage[:width+1]
				secondLine := fmtMessage[width+1:]
				if len(secondLine) > width-6 {
					secondLine = secondLine[:width-6] + "…"
				}
				buf.SetString(2, 1+2*_i, firstLine, 35, 47)
				buf.SetString(4, 2+2*_i, secondLine, 35, 47)
			} else {
				buf.SetString(2, 1+2*_i, fmtMessage, 35, 47)
			}
		}
	}
	self.locker.Unlock()
	return buf
}

func messageFormatter(name, content string) string {
	return fmt.Sprintf("[%s]: %s", name, content)
}
