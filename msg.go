package svc

import ()

type Message struct {
	ID      string
	Command string
	Data    interface{}
}
type Response struct {
	Status int64
	Message
}
