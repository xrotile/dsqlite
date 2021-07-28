package server

import (
	"github.com/goraft/raft"
)

type WriteCommand struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func NewWriteCommand(key string, value string) *WriteCommand {
	return &WriteCommand{
		Key:   key,
		Value: value,
	}
}

func (w *WriteCommand) CommandName() string {
	return "write"
}

func (w *WriteCommand) Apply(server raft.Server) (interface{}, error) {
	db := server.Context().(*DB)
	db.Put(w.Key, w.Value)
	return nil, nil
}
