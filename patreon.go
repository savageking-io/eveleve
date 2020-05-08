package main

import (
	"gopkg.in/mxpv/patreon-go.v1"
)

type Patreon struct {
}

func (p *Patreon) Init() error {
	patreon.NewClient(nil)
	return nil
}
