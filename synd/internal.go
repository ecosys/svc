package main

import (
	"encoding/gob"
	"fmt"
	"github.com/ecosys/svc"
	"github.com/ecosys/synd"
	"log"
	"net"
)

func NewInternal(port int) (internal, error) {
	intern := internal{port}
	return intern, nil
}

type internal struct {
	port int
}

func (intern *internal) Publish(acts []synd.Action) (*svc.Response, error) {
	//log.Println("conn to tcp")
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", intern.port))
	defer conn.Close()

	if err != nil {
		log.Println("error", err)
	}

	acts[0].Param["subject"] = []string{"testing from tcp"}
	acts[0].Param["body"] = []string{"body testing from tcp"}

	gob.Register(synd.Action{})
	gob.Register([]synd.Action{})
	gob.Register(synd.Report{})
	gob.Register(svc.Message{})
	gob.Register(svc.Response{})

	enc := gob.NewEncoder(conn)

	msg := svc.Message{}
	msg.Data = acts
	msg.Command = "publish"

	enc.Encode(&msg)

	var rmsg svc.Response

	log.Println("reading from tcp conn")

	dec := gob.NewDecoder(conn)
	dec.Decode(&rmsg)

	return &rmsg, nil
}
