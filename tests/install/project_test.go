package install

import (
	"testing"

	"gitlab.com/locke-codes/container-cli/internal/install"
)

var err error

func TestCloneProject(t *testing.T) {
	type testCase struct {
		name        string
		projectName string
		projectUrl  string
	}

	tests := []testCase{
		{
			name:        "go-world-https",
			projectName: "go-world-https",
			projectUrl:  "https://gitlab.com/locke-codes/go-world.git",
		},
		{
			name:        "go-world-ssh",
			projectName: "go-world-ssh",
			projectUrl:  "ssh://git@gitlab.com/locke-codes/go-world.git",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempDir := t.TempDir()
			project := install.Project{
				Name:                 test.projectName,
				URL:                  test.projectUrl,
				DestinationDirectory: tempDir,
			}
			err = project.Clone()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
