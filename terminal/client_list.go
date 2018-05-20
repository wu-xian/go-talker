package terminal

import (
	"sync"

	"github.com/cjbassi/termui"
)

type ClientList struct {
	*termui.Block
	InChan       chan *Client
	OutChan      chan *Client
	clients      []*Client
	clientLocker sync.Mutex
}

type Client struct {
	Id   int32
	Name string
}

func (self *ClientList) Add(client *Client) {
	self.clientLocker.Lock()
	self.clients = append(self.clients, client)
	self.clientLocker.Unlock()
	self.InChan <- client
}

func (self *ClientList) Remove(client *Client) {
	self.clientLocker.Lock()
	index := -1
	for i, v := range self.clients {
		if v.Id == client.Id {
			index = i
			break
		}
	}
	if index != -1 {
		self.clients = append(self.clients[:index], self.clients[index+1:]...)
	}
	self.clientLocker.Unlock()
	self.OutChan <- client
}

func NewClientList() *ClientList {
	return &ClientList{
		Block:   termui.NewBlock(),
		OutChan: make(chan *Client, 0),
		InChan:  make(chan *Client, 0),
		clients: []*Client{},
	}
}

func (self *ClientList) Buffer() *termui.Buffer {
	buf := self.Block.Buffer()
	self.clientLocker.Lock()
	for i, v := range self.clients {
		buf.SetString(3, 2*(i+1), v.Name, 35, 47)
	}
	self.clientLocker.Unlock()
	return buf
}
