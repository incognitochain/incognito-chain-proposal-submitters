package entities

import "fmt"

type RPCError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	StackTrace string `json:"stacktrace"`
}

type RPCBaseRes struct {
	Id       int       `json:"id"`
	RPCError *RPCError `json:"error"`
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("%+v %+v %+v", e.Code, e.Message, e.StackTrace)
}
