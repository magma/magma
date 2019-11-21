MK_DIR := $(ROOT_DIR)/.mk-lib
include $(MK_DIR)/variables.mk
-include $(MK_DIR)/version.mk
-include $(ROOT_DIR)/.make.env
-include .make.env

f ?= $(DOCKER_COMPOSE_FILE)
DOCKER_COMPOSE_FILE := $(f)

.DEFAULT_GOAL := help

help: ##@other Show this help.
	@perl -e '$(HELP_FUN)' $(MAKEFILE_LIST)

confirm:
	@( read -p "$(RED)Are you sure? [y/N]$(RESET): " sure && case "$$sure" in [yY]) true;; *) false;; esac )

mk-upgrade: ##@other Check for updates of mk-lib
	@MK_VERSION=$(MK_VERSION) MK_REPO=$(MK_REPO) $(MK_DIR)/self-upgrade.sh

mk-version: ##@other Show current version of mk-lib
	@echo $(MK_VERSION)

check-dependencies:
	@echo Checking dependencies