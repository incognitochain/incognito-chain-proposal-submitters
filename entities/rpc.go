package entities

type RPCError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	StackTrace string `json:"stacktrace"`
}

type RPCBaseRes struct {
	Id       int       `json:"id"`
	RPCError *RPCError `json:"error"`
}
