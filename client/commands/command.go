package commands

import "encoding/json"

type Command interface {
	GetName() string
	GetData() (json.RawMessage, error)
}
