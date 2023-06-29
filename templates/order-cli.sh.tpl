#!/usr/bin/env bash
COMMAND=${1:-install}

echo "##################################################################################################################
################ __NAME__ CLI Installer
#######################################################################################################################"
echo "Ordering $(pwd)/order-__NAME__-cli.sh ${COMMAND}"

INSTALL_LOCATION="${HOME}/.local/share/__NAME__"
CONFIG_PATH="${HOME}/.config/__NAME__/.config.sh"
CCLI_REPO="git@gitlab.com:locke-codes/container-cli.git"
CCLI_INSTALL_LOCATION="${HOME}/.local/share/container-cli"
DEV_CCLI_REPO_DIR=${2-"${HOME}/dev/container-cli"}
DEV_REPO_DIR=${3-"${HOME}/dev/__NAME__"}
REPO=__REPO__


imports() {
  echo "Sourcing imports"
  # "import" common dishes
  # this also includes global vars
  . "${HOME}/.local/share/container-cli/common-dishes.sh"
  # this can fail on install
  . "${HOME}/.config/container-cli/.config.sh" || true
}

## Clone a repository. Replicates function in common-dishes but needs to be explicit here to support curl install
function _clone_repo() {
  install_location=$1
  repo=$2
  echo "Cloning repo"
  if test -d "${install_location}"
  then
    echo "Repo exists. Replacing"
    rm -rf "${install_location}"
    git clone "${repo}" "${install_location}"
  else
    echo "Repo does not exist. Cloning..."
    git clone "${repo}" "${install_location}"
  fi
}

function _copy_repo() {
  install_location=$1
  repo_location=$2
  echo "Copying repo"
  if test -d "${install_location}"
  then
    echo "Repo exists. Deleting"
    rm -rf "${install_location}"
  fi
  pushd .. || exit 1
  echo "Copying ${repo_location} to ${install_location}"
  cp -a "${repo_location}" "${install_location}"
  popd || exit 1
}

function has_git() {
  if test -f "$(which git)"; then
    echo true
  else
    echo false
  fi
}

function _install_ccli() {
  dev=$1
  _has_git=$(has_git)
  if $dev
  then
    # TODO: ask where the repo directory
    _copy_repo "${CCLI_INSTALL_LOCATION}" "${DEV_CCLI_REPO_DIR}"
  else
    if ! ${_has_git}
    then
      echo "Please install git"
    fi
    _clone_repo "${CCLI_INSTALL_LOCATION}" "${CCLI_REPO}"
  fi

  pushd "${CCLI_INSTALL_LOCATION}" || (echo failed to enter src directory; exit 1)
  . "${HOME}/.local/share/container-cli/common-dishes.sh"
  check_requirements
  create_config "${CCLI_INSTALL_LOCATION}" "${CCLI_CONFIGURATOR_PATH}" "${CCLI_CONFIG_PATH}"
  create_cli "${CCLI_INSTALL_LOCATION}" "serve-container-cli.sh" "${BIN_LOCATION}/ccli"
  popd || exit 1
}

function _install() {
  dev=$1
  _has_git=$(has_git)
  if $dev
  then
    _copy_repo "${INSTALL_LOCATION}" "${DEV_REPO_DIR}"
  else
    if ! ${_has_git}
    then
      echo "Please install git"
    fi
    _clone_repo "${INSTALL_LOCATION}" "${REPO}"
  fi

  pushd "${INSTALL_LOCATION}" || (echo failed to enter src directory; exit 1)
  . "${HOME}/.local/share/container-cli/common-dishes.sh"
  check_requirements
  create_config "${INSTALL_LOCATION}" "${CCLI_CONFIGURATOR_PATH}" "${CONFIG_PATH}"
  create_cli "${INSTALL_LOCATION}" "serve-container-cli.sh" "${BIN_LOCATION}/ccli"
  popd || exit 1
}

case ${COMMAND} in
  requires)
    pushd "${CCLI_INSTALL_LOCATION}" || (echo failed to enter src directory; exit 1)
    imports
    check_requirements
    popd || exit 1
    ;;
  install-ccli)
    _install_ccli false
    ;;
  install-ccli-dev)
    _install_ccli true
    ;;
  install)
    _install false
    ;;
  install-dev)
    _install true
    ;;
  config)
    pushd "${CCLI_INSTALL_LOCATION}" || (echo failed to enter src directory; exit 1)
    imports
    create_config "${CCLI_INSTALL_LOCATION}" "${CCLI_CONFIGURATOR_PATH}" "${CCLI_CONFIG_PATH}"
    popd || exit 1
    ;;
  cli)
    pushd "${CCLI_INSTALL_LOCATION}" || (echo failed to enter src directory; exit 1)
    imports
    create_cli "${CCLI_INSTALL_LOCATION}" "serve-container-cli.sh" "${BIN_LOCATION}/ccli"
    popd || exit 1
    ;;
esac
