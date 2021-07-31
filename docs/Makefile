.PHONY: precommit dev

dev:  ## Start local docs server with live reload
	make -C docusaurus dev

precommit:  ## Run docs precommit checks
	make -C readmes precommit

precommit_fix:  ## Try to fix existing precommit issues
	make -C readmes precommit_fix

# Ref: https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:  ## Show documented commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'
