package net

import "sync"

//Request manage a websocket request
type Request struct {
	ID       int64
	JSON     string
	Response *Response
	Wait     *sync.WaitGroup
}

//NewRequest is constructor
func NewRequest(id int64, json string) *Request {
	request := &Request{
		ID:   id,
		JSON: json,
		Wait: new(sync.WaitGroup),
	}
	return request
}
