CONFIG_PATH ?= "${HOME}/.config/container-cli"
INSTALL_LOCATION ?= "${HOME}/.local/share/container-cli"

.PHONY: clean
clean:
	printf "Removing config path and install location\n\tconfig path: %s\n\tinstall path: %s"${CONFIG_PATH} ${INSTALL_LOCATION}
	rm -rf ${CONFIG_PATH} ${INSTALL_LOCATION}

.PHONY: install
install:
	./order-container-cli.sh install

.PHONY: install-dev
install-dev:
	./order-container-cli.sh install-dev
