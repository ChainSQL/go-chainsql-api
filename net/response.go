package net

// Response manage a websocket response
type Response struct {
	Value        string
	Request      *Request
	Error        string
	ErrorMessage string
}
