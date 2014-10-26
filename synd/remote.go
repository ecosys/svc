package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ecosys/synd"
	"net/http"
)

func NewRemote(url string) (remote, error) {
	svc := remote{url}
	return svc, nil
}

type remote struct {
	url string
}

func (rem *remote) Publish(acts []synd.Action) (*synd.Report, error) {

	b, err := json.Marshal(acts)

	resp, err := http.Post(fmt.Sprintf("%s/publish", rem.url), "", bytes.NewBuffer(b))

	var res synd.Report

	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(&res)

	return &res, err
}
