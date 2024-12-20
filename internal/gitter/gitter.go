package gitter

import (
	"github.com/hashicorp/go-getter"
	"net/url"
)

type Gitter struct {
	Name        string
	Url         *url.URL
	Destination string
	_client     getter.GitGetter
}

func NewGitter(name, gitUrl, destination string) *Gitter {
	thisUrl, err := url.Parse(gitUrl)
	if err != nil {
		panic(err)
	}
	return &Gitter{
		Name:        name,
		Url:         thisUrl,
		Destination: destination,
		_client:     getter.GitGetter{},
	}
}

func (g *Gitter) Clone() error {
	err := g._client.Get(g.Destination, g.Url)
	if err != nil {
		return err
	}
	return nil
}
