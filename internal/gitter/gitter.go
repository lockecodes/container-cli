package gitter

import (
	"net/url"
	"os"

	"github.com/hashicorp/go-getter"
)

// Gitter represents a structure for managing Git operations with a repository URL and destination path.
type Gitter struct {
	Name        string
	Url         *url.URL
	Destination string
	_client     getter.GitGetter
}

// NewGitter creates and returns a new Gitter instance initialized with the provided name, Git repository URL,
// and destination path.
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

// Clone retrieves the repository from the configured URL and stores it in the destination directory.
func (g *Gitter) Clone() error {
	var err error
	err = os.RemoveAll(g.Destination)
	if err != nil {
		println(err)
	}
	err = g._client.Get(g.Destination, g.Url)
	if err != nil {
		return err
	}
	return nil
}
