#!/usr/bin/env bash
DEBUG=${DEBUG:-false}
if ${DEBUG}; then
  set -ex
fi
# subcommand
COMMAND=${1:-help}
# shift past the subcommand
shift

Green='\033[0;32m'        # Green
Color_Off='\033[0m'       # Text Reset
BIWhite='\033[1;97m'      # White

_imports() {
  # Source the config file
  . "${HOME}/.config/container-cli/.config.sh"
  # "import" common dishes
  . "${CCLI_INSTALL_LOCATION}/common-dishes.sh"
  # "import" app config
  . "${HOME}/.config/__NAME__/.config.sh"
}

CONTEXT_DIR="$(pwd)"


all_parms=$@
parms="${all_parms/install /}"
help()
{
    printf "%s${BIWhite}
##################################################################################################################
##################################################################################################################
########### ${Green}__NAME__ CLI${Color_Off}
########### This is a dynamically generated cli for __NAME__
##################################################################################################################
##################################################################################################################
# ${Green}Commands:${Color_Off}
  ${Green}build${Color_Off}: Build the container used to run commands
  ${Green}shell${Color_Off}: Exec into the cli container
  ${Green}help${Color_Off}: Display this help text then try to pass help command to the container
  ${Green}*${Color_Off}: Any other command is passed directory to the container
"
}

function _run_command() {
  _imports
  # shellcheck disable=SC2048
  pushd "${INSTALL_LOCATION}" \
    && CONTEXT_DIR="${CONTEXT_DIR}" \
      ${COMPOSE} \
        -f docker-compose.yaml \
        $* \
    ; ${COMPOSE} down \
    ; popd
}

if ${DEBUG} | ${CONFIG}; then
  _run_command config
fi

case "${COMMAND}" in
  build)
    _run_command build __NAME__
    ;;
  shell)
    _run_command run --rm __NAME__ bash
    ;;
  help)
    help
    _run_command run --rm __NAME__ "${COMMAND_PREFIX}" --help
    ;;
  *)
    _run_command run --rm __NAME__ "${COMMAND_PREFIX}" "${COMMAND}" "$@"
    ;;
esac