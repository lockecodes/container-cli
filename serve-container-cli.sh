#!/usr/bin/env bash

DEBUG=${DEBUG:-false}
if ${DEBUG}; then
  set -ex
fi
# subcommand
COMMAND=${1:-help}
# shift past the subcommand
shift

## Source the config file
echo "sourcing config file: ${HOME}/.config/container-cli/.config.sh"
. "${HOME}/.config/container-cli/.config.sh"

Green='\033[0;32m'        # Green
Color_Off='\033[0m'       # Text Reset
BIWhite='\033[1;97m'      # White

_imports() {
  # Source the config file
  . "${HOME}/.config/container-cli/.config.sh"
  # "import" common dishes
  . "${CCLI_INSTALL_LOCATION}/common-dishes.sh"
}

all_parms=$@
parms="${all_parms/install /}"
help()
{
    printf "%s${BIWhite}
##################################################################################################################
##################################################################################################################
########### ${Green}Container CLI${Color_Off}
########### A 'container cli' is a command line input that executes containers (docker | podman) for processing
########### This application can be used to add a container cli to your repo or to just install a git repository that
########### has a compatible container
##################################################################################################################
##################################################################################################################
# ${Green}Commands:${Color_Off}
  ${Green}git-install${Color_Off}: Clone a git repository into a runtime directory and install cli scripts
    :-r/--repository: git repository url. e.g. git@gitlab.com:locke-codes/container-cli.git
    :-n/--name: cli command name
    :-d/--destination: Optional destination directory (default: ${CCLI_INSTALL_LOCATION_ROOT}/{name}/src
    :-c/--command-prefix: Container command prefix. This is an executable like 'python' or an entrypoint command
  ${Green}template${Color_Off}: install cli scripts into a directory using templates
    :-r/--repository: git repository url. e.g. git@gitlab.com:locke-codes/container-cli.git
    :-n/--name: cli command name
    :-d/--destination: Optional destination directory (default: ${CCLI_INSTALL_LOCATION_ROOT}/{name}/src
    :-c/--command-prefix: Container command prefix. This is an executable like 'python' or an entrypoint command
  ${Green}help${Color_Off}: Display this help text
  ${Green}update-ccli${Color_Off}: Update the ccli to head or a ref
  ${Green}update${Color_Off}: Update the command repo to head or a ref
    :-n/--name: cli command to update
    :-d/--destination: Optional destination directory (default: ${CCLI_INSTALL_LOCATION_ROOT}/{name}/src
    :-r/--ref: Optional git ref. (default: head)
"
    exit 2
}

case "${COMMAND}" in
  git-install)
    _imports
    git_install "$@"
    ;;
  update-ccli)
    _imports
    pushd "${CCLI_INSTALL_LOCATION}" || exit \
    && git add . \
    && git stash \
    && git pull origin main \
    && git stash pop \
    && popd || exit
    ;;
  update)
    _imports
    update "$1"
    ;;
  template)
    _imports
    create_templates "$@"
    ;;
  help)
    help
    ;;
  *)
    echo "Unsupported command: ${COMMAND}"
    help
    ;;
esac
