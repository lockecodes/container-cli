#!/usr/bin/env bash
COMMAND=${1:-install}

echo "##################################################################################################################
################ Container CLI Installer
#######################################################################################################################"
echo "Ordering $(pwd)/order-container-cli.sh ${COMMAND}"

CCLI_INSTALL_LOCATION="${HOME}/.local/share/container-cli"
CCLI_REPO="git@gitlab.com:locke-codes/container-cli.git"

imports() {
  echo "Sourcing imports"
  # "import" common dishes
  # this also includes global vars
  . "${HOME}/.local/share/container-cli/common-dishes.sh"
}

## Clone a repository. Replicates function in common-dishes but needs to be explicit here to support curl install
function _clone_repo() {
  echo "Cloning repo"
  if test -d "${CCLI_INSTALL_LOCATION}"
  then
    echo "Repo exists. Replacing"
    rm -rf "${CCLI_INSTALL_LOCATION}"
    git clone "${CCLI_REPO}" "${CCLI_INSTALL_LOCATION}"
  else
    echo "Repo does not exist. Cloning..."
    git clone "${CCLI_REPO}" "${CCLI_INSTALL_LOCATION}"
  fi
}

function _copy_repo() {
  echo "Copying repo"
  if test -d "${CCLI_INSTALL_LOCATION}"
  then
    echo "Repo exists. Deleting"
    rm -rf "${CCLI_INSTALL_LOCATION}"
  fi
  this_dir="$(pwd)"
  pushd .. || exit 1
  echo "Copying ${this_dir} to ${CCLI_INSTALL_LOCATION}"
  cp -a "${this_dir}" "${CCLI_INSTALL_LOCATION}"
  popd || exit 1
}

function has_git() {
  if test -f "$(which git)"; then
    echo true
  else
    echo false
  fi
}

case ${COMMAND} in
  requires)
    pushd "${CCLI_INSTALL_LOCATION}" || (echo failed to enter src directory; exit 1)
    imports
    check_requirements
    popd || exit 1
    ;;
  install)
    _has_git=$(has_git)
    if ! ${_has_git}
    then
      echo "Please install git"
    fi
    _clone_repo
    pushd "${CCLI_INSTALL_LOCATION}" || (echo failed to enter src directory; exit 1)
    imports
    check_requirements
    create_config "${CCLI_INSTALL_LOCATION}" "${CCLI_CONFIGURATOR_PATH}" "${CCLI_CONFIG_PATH}" "${CCLI_REPO}"
    create_cli "${CCLI_INSTALL_LOCATION}" "serve-container-cli.sh" "${BIN_LOCATION}/ccli"
    popd || exit 1
    ;;
  install-dev)
    _copy_repo
    pushd "${CCLI_INSTALL_LOCATION}" || (echo failed to enter src directory; exit 1)
    imports
    check_requirements
    create_config "${CCLI_INSTALL_LOCATION}" "${CCLI_CONFIGURATOR_PATH}" "${CCLI_CONFIG_PATH}" "${CCLI_REPO}"
    create_cli "${CCLI_INSTALL_LOCATION}" "serve-container-cli.sh" "${BIN_LOCATION}/ccli"
    popd || exit 1
    ;;
  clone)
    pushd "${CCLI_INSTALL_LOCATION}" || (echo failed to enter src directory; exit 1)
    imports
    check_requirements
    popd || exit 1
    clone_or_replace "${CCLI_REPO}" "${CCLI_INSTALL_LOCATION}"
    ;;
  config)
    pushd "${CCLI_INSTALL_LOCATION}" || (echo failed to enter src directory; exit 1)
    imports
    create_config "${CCLI_INSTALL_LOCATION}" "${CCLI_CONFIGURATOR_PATH}" "${CCLI_CONFIG_PATH}" "${CCLI_REPO}"
    popd || exit 1
    ;;
  cli)
    pushd "${CCLI_INSTALL_LOCATION}" || (echo failed to enter src directory; exit 1)
    imports
    create_cli "${CCLI_INSTALL_LOCATION}" "serve-container-cli.sh" "${BIN_LOCATION}/ccli"
    popd || exit 1
    ;;
esac
