package main

import (
	"github.com/ecosys/synd"
)

func newService() (service, error) {
	svc := service{}
	return svc, nil
}

type service struct {
}

func (s *service) Publish(acts []*synd.Action) (synd.Report, error) {
	syn, err := synd.NewSyndicator(acts)
	rep, err := syn.Async()
	return rep, err
}
