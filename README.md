# Container CLI
A 'container cli' is a command line input that executes containers (docker | podman) for processing
This application can be used to add a container cli to your repo or to just install a git repository that has a 
compatible container

# Motivation
I have custom tooling repos that I have recently converted to run completely from docker with a local cli entrypoint.
This was just an experiment to see if it could be generic and shared.

# Installation
```shell
curl https://gitlab.com/locke-codes/container-cli/-/raw/main/order-container-cli.sh | bash
```

# Clone a repo and install a cli
Install tpc-c-bench to the cli command tpcb
```shell
ccli git-install -n tpcb -r 'git@gitlab.com:locke-codes/tpc_c_bench.git' -c tpcb
```
Now this repo is available to run inside the docker container from the repo like:
```shell
>> tpcb
Usage: tpcb [OPTIONS] COMMAND [ARGS]...

Options:
  --help  Show this message and exit.

Commands:
  create-schema
  drop-schema
  populate-data
  run
```
```shell
>> tpcb run --help
Usage: tpcb run [OPTIONS]

Options:
  --url TEXT
  --driver TEXT
  --warehouses INTEGER
  --duration INTEGER
  --print-interval INTEGER
  --help Show this message and exit.
```

# Add integration directly to a repository
Add templated files to a current project to provide "order" functionality and a cli entry point
```shell
ccli template -n tpcb -r 'git@gitlab.com:locke-codes/tpc_c_bench.git' -d ~/dev/tpc_c_bench -c tpcb
```
This will add the templated files for:
* .env
* docker-compose.yaml
* order-cli.sh
* server-cli.sh

to the working directory.

The provided functionality allows a simple installation script like:
```shell
curl https://gitlab.com/locke-codes/tpc_c_bench/-/raw/main/order-tpcb.sh | bash
```

# Commands in repos
The following commands are available after installation for the example command name of tpcb:
```shell
tpcb build # build the docker container. This will happen the first time you run anyways
tpcb shell # entry a shell inside the container
tpcb help # display help
tpcb {anything} # pass through any command to the container
```

# Ccli commands
```shell
##################################################################################################################
##################################################################################################################
########### Container CLI
########### A 'container cli' is a command line input that executes containers (docker | podman) for processing
########### This application can be used to add a container cli to your repo or to just install a git repository that
########### has a compatible container
##################################################################################################################
##################################################################################################################
# Commands:
  git-install: Clone a git repository into a runtime directory and install cli scripts
    :-r/--repository: git repository url. e.g. git@gitlab.com:locke-codes/container-cli.git
    :-n/--name: cli command name
    :-d/--destination: Optional destination directory (default: /home/slocke/.local/share/{name}/src
  template: install cli scripts into a directory using templates
    :-r/--repository: git repository url. e.g. git@gitlab.com:locke-codes/container-cli.git
    :-n/--name: cli command name
    :-d/--destination: Optional destination directory (default: /home/slocke/.local/share/{name}/src
  help: Display this help text
  update-ccli: Update the ccli to head or a ref
  update: Update the command repo to head or a ref
    :-n/--name: cli command to update
    :-d/--destination: Optional destination directory (default: /home/slocke/.local/share/{name}/src
    :-r/--ref: Optional git ref. (default: head)
```
