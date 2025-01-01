# Documentation

## Introduction

Container CLI is a powerful tool designed to install and execute containerized projects as local commands seamlessly.
This document will provide steps on installing, using, and maintaining this CLI tool effectively.

## Installation Guide

### Step 1: Download Release

1. Navigate to the [Container CLI Releases](https://gitlab.com/locke-codes/container-cli/-/releases).
2. Download the appropriate release for your operating system.

### Step 2: Extract and Install

1. Open a terminal and navigate to the directory where the release file is downloaded.
2. Extract the downloaded tar file:
   ```bash
   tar -xvf ~/Downloads/container-cli_Linux_x86_64.tar.gz
   ```
3. Execute the binary file to install `container-cli`:
   ```bash
   ./container-cli install
   ```

## Usage Instructions

### Installing a New Project

1. Run the following command to install a project:
   ```bash
   ccli project install \
    --name <project_name> \
    --url <git_url> \
    --dest <destination_path> \
    --command <default_command> \
    --alias <command_alias>
   ```

   Example:
   ```bash
   ccli project install \
    --name big-salad \
    --url ssh://git@gitlab.com/locke-codes/big-salad.git \
    --dest ~/.local/share \
    --command bs \
    --alias bs
   ```
2. Replace `<project_name>`, `<git_url>`, `<destination_path>`, `<default_command>`, `<command_alias>` with appropriate values.
3. Run the command using the command alias like: `bs --help`. NOTE: If you do not pass the alias flag the command will default to the project name

### Interactive Mode for Installing a New Project

You can use the interactive mode to install a new project. This allows you to input the required fields step-by-step
directly in the terminal.

Run the following command:

```bash
ccli project install
```

You will then be prompted to enter the required details interactively, such as:

1. **Project Name**: Enter the name of the project.
2. **Git URL**: Provide the repository's Git URL.
3. **Destination Path**: Specify the path where the project should be installed.
4. **Command**: The default command to execute in the docker container

For example, the session may look like:

```shell
❯ ccli project install
Enter Project Name: big-salad
Enter Git URL: ssh://git@gitlab.com/locke-codes/big-salad.git
Enter Destination Path: /home/user/projects
Enter Command: bs
```

Once all inputs are provided, the project will be installed and can be executed using the specified alias.

Example:

```bash
big-salad format yaml test.yaml
```

### Running a Project

After installing a project, you can run it using the alias command you specified during installation. For example:

```bash
big-salad format yaml test.yaml
```

### Updating the Tool

To update Container CLI to the latest version:

```bash
ccli update
```

### Checking Version

To check the current version of Container CLI:

```bash
ccli version
```

## Commands Overview

Below are some of the primary commands for `ccli`:

| Command   | Description                                    |
|-----------|------------------------------------------------|
| `install` | Installs the ContainerCLI binary.              |
| `update`  | Updates the CLI tool to the latest version.    |
| `version` | Displays the current version of the CLI.       |
| `project` | Manage projects (install, remove, etc.).       |
| `help`    | Shows help for commands or a list of commands. |

## TODO

1. Add support for Docker and refactor calls currently specific to Podman.
2. Add Makefile targets for usage during ci/cd to run tests etc. Bring back Dockerfile was it is useful
3. Add more complete unit testing 
4. Create and include a usage demonstration video.
