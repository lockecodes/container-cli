#!/usr/bin/env bash

BIN_LOCATION="${HOME}/.local/bin"
CCLI_SERVE_PATH="${BIN_LOCATION}/ccli"
CCLI_INSTALL_LOCATION="${HOME}/.local/share/container-cli"
CCLI_INSTALL_LOCATION_ROOT=$(dirname "${CCLI_INSTALL_LOCATION}")
CCLI_CONFIG_PATH="${HOME}/.config/container-cli/.config.sh"
CCLI_CONFIGURATOR_PATH="${HOME}/.local/share/container-cli/config.sh"
CCLI_REPO="git@gitlab.com:locke-codes/container-cli.git"

export BIN_LOCATION CCLI_SERVE_PATH CCLI_INSTALL_LOCATION CCLI_INSTALL_LOCATION_ROOT CCLI_CONFIG_PATH CCLI_CONFIGURATOR_PATH CCLI_REPO

function pushd () {
    command pushd "$@" > /dev/null
}

function popd () {
    command popd > /dev/null
}

dish_imports() {
  # Source the config file
  . "${HOME}/.config/container-cli/.config.sh"
  # "import" common dishes
  . "${CCLI_INSTALL_LOCATION}/common-dishes.sh"
}

exec_containers() {
  PODMAN=$(which podman)
  DOCKER=$(which docker)
  PODMAN_COMPOSE=$(which podman-compose || false)
  DOCKER_COMPOSE=$(which docker-compose || false)
  export PODMAN DOCKER PODMAN_COMPOSE DOCKER_COMPOSE
}

########
### Return true if podman is installed
########
function has_podman() {
  exec_containers
  if test -f "${PODMAN}"; then
      echo true
  else
      echo false
  fi
}

########
### Return true if podman-compose is installed
########
function has_podman_compose() {
  exec_containers
  if test -f "${PODMAN_COMPOSE}"; then
      echo true
  else
      echo false
  fi
}

########
### Return true if docker is installed
########
function has_docker() {
  exec_containers
  if test -f "${DOCKER}"; then
      echo true
  else
      echo false
  fi
}

########
### Return true if docker-compose is installed
########
function has_docker_compose() {
  exec_containers
  if test -f "${DOCKER_COMPOSE}"; then
      echo true
  else
      echo false
  fi
}

########
### Which container: Return docker, podman, or false
########
function which_container() {
  pod=$(has_podman)
  dock=$(has_docker)
  if ${dock}; then { echo "docker"; return 0; }; fi
  if ${pod}; then { echo "podman"; return 0; }; fi
  echo false
}

########
### Which compose: Return docker-compose, podman-compose, or false
########
function which_compose() {
  pod=$(has_podman_compose)
  dock=$(has_docker_compose)
  if ${dock}; then { echo "docker-compose"; return 0; }; fi
  if ${pod}; then { echo "podman-compose"; return 0; }; fi
  echo false
}

function has_git() {
  if test -f "$(which git)"; then
    echo true
  else
    echo false
  fi
}

########
### Clone the project to the install location
########
function clone_repo() {
  repository=$1
  install_location=$2
  mkdir -p "$(dirname "${install_location}")"
  git clone "${repository}" "${install_location}"
}

########
### Clone the repo
########
function clone_or_replace() {
  repository=$1
  install_location=$2
  echo "install_location passed as ${install_location}"
  if test -d "${install_location}"
  then
    echo "Installation directory exists"
    read -rp "Update? [y/n]" confirm

    if [[ ${confirm} == [yY] || ${confirm} == [yY][eE][sS] ]]
    then
      echo "Updating"
      rm -rf "${install_location}"
      clone_repo "${repository}" "${install_location}"
    else
      echo "Exiting..."
      exit 0
    fi
  else
    clone_repo "${repository}" "${install_location}"
  fi
}

########
### Configuration
########
function create_config() {
  install_location=$1
  configurator_path=$2
  config_path=$3
  repository=$4
  command_prefix=$5
  config_dir=$(dirname "${config_path}")
  mkdir -p "${config_dir}"
  # TODO: get docker/podman prefs from user
  bash "${configurator_path}" "${config_path}" \
    "${install_location}" docker docker-compose "${repository}" "${command_prefix}"
}

########
### Setup cli
########
function create_cli() {
  install_location=$1
  script_name=$2
  bin_path=$3
  echo "install_location passed as ${install_location}"
  # TODO: Ensure local bin location is in path
  mkdir -p "${install_location}"
  ln -sf "${install_location}/${script_name}" "${bin_path}"
}

########
### Create templated files
########
function create_templated_files() {
  destination=$1
  repository=$2
  name=$3

  create_template_file "${destination}" \
    "${name}" "${repository}" "${destination}/serve-${name}.sh" \
    "${CCLI_INSTALL_LOCATION}/templates/serve-cli.sh.tpl"
  create_template_file "${destination}" \
    "${name}" "${repository}" "${destination}/order-${name}.sh" \
    "${CCLI_INSTALL_LOCATION}/templates/order-cli.sh.tpl"
  create_template_file "${destination}" \
    "${name}" "${repository}" "${destination}/.env" \
    "${CCLI_INSTALL_LOCATION}/templates/.env.tpl"
  create_template_file "${destination}" \
    "${name}" "${repository}" "${destination}/docker-compose.yaml" \
    "${CCLI_INSTALL_LOCATION}/templates/docker-compose.yaml.tpl"
}
########
### Create templated files
########
function create_templates() {
  dish_imports
  destination=""
  repository=""
  name=""

  VALID_ARGS=$(getopt -o r:n:d: -- "$@")
  if [[ $? -ne 0 ]]; then
    help
    exit 1;
  fi

  if ! [[ ${VALID_ARGS} == *"-r"* && ${VALID_ARGS} == *"-n"* && ${VALID_ARGS} == *"-d"* ]]; then
    echo "create-templated-files requires repository, name, and destination"
    help
    exit 1;
  fi

  eval set -- "$VALID_ARGS"
  while [ : ]; do
    case "$1" in
      -r | --repository)
          repository=$2
          shift 2
          ;;
      -n | --name)
          name=$2
          shift 2
          ;;
      -d | --destination)
          destination="$2"
          shift 2
          ;;
      --) shift;
          break
          ;;
    esac
  done

  if [[ ${destination} == "" ]]; then
    destination="${INSTALL_LOCATION_ROOT}/${name}"
  fi

  create_templated_files "${destination}" "${repository}" "${name}"
}

########
### template file
########
function create_template_file() {
  destination=$1
  name=$2
  repository=$3
  script_path=$4
  template_file=$5
  rm -f "${script_path}"
  mkdir -p "${destination}"
  tpl=$(cat "${template_file}")
  tpl=${tpl//__REPO__/${repository}}
  tpl=${tpl//__NAME__/${name}}
  cat << EOF > "${script_path}"
${tpl}
EOF
  chmod +x "${script_path}"
}

########
### Make sure requirements are installed. Fail out if missing
########
function check_requirements() {
  exec_containers
  echo "################ Checking prerequisites"
  echo "################ Checking  container requirement"
  _has_container=$(which_container)
  _has_compose=$(which_compose)
  _has_git=$(has_git)
  case ${_has_container} in
    docker)
      echo "Docker support found";
      _has_container=true;
      ;;
    podman)
      echo "Podman support found";
      _has_container=true;
      ;;
    false)
      echo "
      No supported container system installed
      Install either docker or podman"
  esac
  case ${_has_compose} in
    docker-compose)
      echo "Docker Compose support found";
      _has_compose=true;
      ;;
    podman-compose)
      echo "Podman Compose support found";
      _has_compose=true;
      ;;
    false)
      echo "
      No supported container compose system installed
      Install either docker-compose or podman-compose"
  esac

  if ! ${_has_git}
  then
    echo "Please install git"
  fi

  if ! ${_has_container} || ! ${_has_compose} || ! ${_has_git}
  then
    echo "Unmet requirements...exiting"
    exit 1
  fi
}

########
### Install git repo
########
function git_install() {
  dish_imports
  repository=""
  name=""
  destination=""
  config_path=""
  command_prefix=""

  VALID_ARGS=$(getopt -o r:n:d:c: --long repository:,name:,destination:,command-prefix: -- "$@")
  if [[ $? -ne 0 ]]; then
    help
    exit 1;
  fi

  if ! [[ ${VALID_ARGS} == *"-r"* && ${VALID_ARGS} == *"-n"* ]]; then
    echo "git-install requires repository and name"
    help
    exit 1;
  fi

  eval set -- "$VALID_ARGS"
  while [ : ]; do
    case "$1" in
      -r | --repository)
          repository=$2
          shift 2
          ;;
      -n | --name)
          name=$2
          shift 2
          ;;
      -d | --destination)
          destination="$2"
          shift 2
          ;;
      -c | --command-prefix)
          command_prefix="$2"
          shift 2
          ;;
      --) shift;
          break
          ;;
    esac
  done

  if [[ ${destination} == "" ]]; then
    destination="${INSTALL_LOCATION_ROOT}/${name}"
  fi
  config_path="${HOME}/.config/${name}/.config.sh"
  check_requirements
  echo "Cloning ${repository} to ${destination} for ${name}"
  clone_or_replace "${repository}" "${destination}/src"
  create_config "${destination}" "${CCLI_CONFIGURATOR_PATH}" "${config_path}" "${repository}" "${command_prefix}"
  create_templated_files "${destination}" "${repository}" "${name}"
  create_cli "${destination}" "serve-${name}.sh" "${BIN_LOCATION}/${name}"
}

########
### update
########
function update() {
  name=$1
  dish_imports
  # shellcheck disable=SC1090
  . "${HOME}/.config/${name}/.config.sh"
  pushd "${INSTALL_LOCATION}/src" \
  && git add . \
  && git stash \
  && (git pull origin main || git pull origin master ) \
  && git stash pop \
  && (popd || true)
  create_templated_files "${INSTALL_LOCATION}" "${REPOSITORY}" "${name}"
  create_config "${INSTALL_LOCATION}" "${CCLI_CONFIGURATOR_PATH}" "${CONFIG_PATH}" "${REPOSITORY}" "${COMMAND_PREFIX}"
  create_cli "${INSTALL_LOCATION}" "serve-${name}.sh" "${BIN_LOCATION}/${name}"
}
