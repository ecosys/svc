package main

import (
	"github.com/ecosys/synd"
	"log"
	"testing"
)

type test struct {
	ID   string
	Name string
	Ch   child
	Cha  []child
}
type child struct {
	ID   string
	Name string
}

func TestCom(t *testing.T) {
	log.Println("conn to http")

	act := synd.Action{}
	act.Provider = synd.Provider{0, "smtp"}

	auth := make(map[string]string)
	auth["server"] = "smtp.gmail.com"
	auth["username"] = "ecosys13@gmail.com"
	auth["password"] = "$cosys13"

	config := make(map[string]string)
	config["server"] = "smtp.gmail.com:587"
	config["sender"] = "ecosys13@gmail.com"

	param := make(map[string][]string)
	param["recipients"] = []string{"ecosys13@gmail.com", "brentmn@gmail.com"}
	param["subject"] = []string{"testing from http"}
	param["body"] = []string{"body testing from http"}

	act.Configure(auth, config, param)
	acts := []synd.Action{act}

	rem, err := NewRemote("http://localhost:8899")

	resp, err := rem.Publish(acts)

	log.Printf("http resp: %v\n", resp)

	intern, err := NewInternal(8898)

	res, err := intern.Publish(acts)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("internal response: ", res)

}
