package install

import "gitlab.com/locke-codes/container-cli/internal/gitter"

type Project struct {
	Name        string
	URL         string
	Destination string
}

func (p *Project) Clone() error {
	client := gitter.NewGitter(p.Name, p.URL, p.Destination)
	err := client.Clone()
	if err != nil {
		return err
	}
	return nil
}
