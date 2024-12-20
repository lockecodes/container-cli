package install

import (
	"os"
	"testing"
)

func TestCloneProject(t *testing.T) {
	type testCase struct {
		name                 string
		projectName          string
		projectUrl           string
		destinationDirectory string
		//expected string
		//hasError bool
	}

	tests := []testCase{
		{
			name:                 "helloworld-https",
			projectName:          "helloworld",
			projectUrl:           "https://github.com/fw876/helloworld.git",
			destinationDirectory: "/tmp/helloworld",
		},
		{
			name:                 "helloworld-ssh",
			projectName:          "helloworld",
			projectUrl:           "ssh://git@github.com/fw876/helloworld.git",
			destinationDirectory: "/tmp/helloworld",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			project := Project{
				Name:                 test.projectName,
				URL:                  test.projectUrl,
				DestinationDirectory: test.destinationDirectory,
			}
			os.RemoveAll(test.destinationDirectory)
			err := project.Clone()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
