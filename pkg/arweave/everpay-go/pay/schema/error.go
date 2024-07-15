package schema

import "encoding/json"

type InternalErr struct {
	Index int    `json:"index"`
	Msg   string `json:"msg"`
}

// NewInternalErr if less than 0 (like -1), means not items error
func NewInternalErr(idx int, msg string) *InternalErr {
	return &InternalErr{
		Index: idx,
		Msg:   msg,
	}
}

func (e InternalErr) Error() string {
	jsErr, _ := json.Marshal(&e)
	return string(jsErr)
}
