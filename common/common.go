package common

// Auth is the type with ws connection infomations
type Auth struct {
	Address string
	Secret  string
	Owner   string
}

//IRequest define interface for request
type IRequest interface {
	GetID() int64
}

//RequestBase contains fields that all requests will have
type RequestBase struct {
	Command string `json:"command"`
	ID      int64  `json:"id,omitempty"`
}

// GetID  return id for request
func (r *RequestBase) GetID() int64 {
	return r.ID
}
