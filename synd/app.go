package main

import (
	//"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ecosys/svc"
	"github.com/ecosys/synd"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	RECV_BUF_LEN = 1024
)

var (
	serv, err = newService()
)

func main() {
	//0 - internal port
	//1 - external port

	iport := os.Args[1]
	eport := os.Args[2]

	log.Println("listening", iport, eport)

	go func(p string) {
		log.Println("listening outside on: ", p)
		http.HandleFunc("/", handle)
		http.ListenAndServe(fmt.Sprintf(":%s", p), nil)
	}(eport)

	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", iport))

	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening inside on: ", iport)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go handleConnection(conn)
	}

}
func handle(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.Write([]byte("we like POSTs around these parts"))
		return
	}

	path := r.URL.Path[1:]

	bts, err := ioutil.ReadAll(r.Body)

	var acts []*synd.Action
	err = json.Unmarshal(bts, &acts)

	if err != nil {
		log.Println("unmarshall error: ", err.Error())
		return
	}

	ref := make([]*synd.Action, 0)

	for _, a := range acts {
		ref = append(ref, a)
	}

	var rep synd.Report

	switch path {
	case "publish":
		rep, err = serv.Publish(ref)
	}

	enc := json.NewEncoder(w)

	enc.Encode(rep)
	return

}
func handleConnection(conn net.Conn) {
	defer conn.Close()

	dec := gob.NewDecoder(conn)

	var msg svc.Message

	gob.Register(synd.Action{})
	gob.Register([]synd.Action{})
	gob.Register(synd.Report{})
	gob.Register(svc.Message{})
	gob.Register(svc.Response{})

	err = dec.Decode(&msg)

	if err != nil && err != io.EOF {
		println("Error reading:", err.Error())
		return
	}

	resp := handleCommand(msg.Command, msg.Data)

	enc := gob.NewEncoder(conn)
	enc.Encode(&resp)

	if err != nil && err != io.EOF {
		println("Error encoding:", err.Error())
		return
	}
	return
}
func handleCommand(path string, data interface{}) svc.Response {
	var err error
	var d interface{}
	resp := svc.Response{}
	resp.Command = path

	switch path {
	case "publish":
		//even if data are pointers, it is flattened and we need to redo here.
		acts, ok := data.([]synd.Action)
		if ok {
			ref := make([]*synd.Action, 0)
			for _, a := range acts {
				ref = append(ref, &a)
			}
			rep, _ := serv.Publish(ref)
			resp.Data = rep
		} else {
			log.Println("could not decode into action array: ", data)
			resp.Data = "need array of Action"
		}
	default:
		log.Println("unknown command: ", path)
		err = errors.New("unknown command: " + path)
	}

	if err != nil {
		resp.Status = -1
		resp.Data = err.Error()
	} else {
		resp.Status = 1
		resp.Data = d
	}

	return resp
}
