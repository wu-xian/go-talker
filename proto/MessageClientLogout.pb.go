// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/mauricelam/genny

package proto

import (
	p "github.com/golang/protobuf/proto"
)

func (self *MessageClientLogout) MessageMarshal() ([]byte, error) {
	return p.Marshal(self)
}

func (self *MessageClientLogout) MessageUnmarshal(bytes []byte) error {
	return p.Unmarshal(bytes, self)
}