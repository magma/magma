.PHONY: dev precommit precommit_fix help

dev:  ## Start local docs server with live reload
	make -C docusaurus dev

precommit:  ## Run docs precommit checks
	make -C readmes precommit

precommit_fix:  ## Try to fix existing precommit issues
	make -C readmes precommit_fix

sidebar_check:  ## Check if all pages are implemented with sidebars for 1.7.0, 1.8.0, and latest
	python3 check_sidebars.py

# Ref: https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:  ## Show documented commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'
