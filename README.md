# Container Cli
This tool is for installing projects with Dockerfiles into local commands.

# Usage
* Install the tool
  * Download a release from https://gitlab.com/locke-codes/container-cli/-/releases
  * Untar the file
    * `tar -xvf ~/Downloads/container-cli_Linux_x86_64.tar.gz`
  * Execute the file to install
    * `./container-cli install`
* Install a project
  * `ccli project install -name big-salad -url ssh://git@gitlab.com/slocke716/big-salad.git -dest /home/slocke/.local/share -command bs`
  * TODO: Replace the above command with a public repository for others
* Run project
  * `big-salad format yaml test.yaml`
* Update this tool
  * `ccli update`
* Version
  * `ccli version`

```shell
❯ ccli
NAME:
   Container CLI - Execute applications in containers

USAGE:
   Container CLI [global options] [command [command options]]

COMMANDS:
   install  ContainerCLI the ccli binary
   update   Update to the latest version of the CLI
   version  Get the version of the CLI
   project  Ccli commands for projects
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

# TODO
TODO: Replace the example with a public repository for others
TODO: Support Docker and update calls that are explicit for Podman

