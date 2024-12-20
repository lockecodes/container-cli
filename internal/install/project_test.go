package install

import "testing"

//git@gitlab.com:slocke716/big-salad.git

func TestCloneProject(t *testing.T) {
	projectUrl := "ssh://git@gitlab.com/slocke716/big-salad.git"
	project := Project{
		Name:        "big-salad",
		URL:         projectUrl,
		Destination: "/tmp/big-salad",
	}
	err := project.Clone()
	if err != nil {
		t.Error(err)
	}
}
