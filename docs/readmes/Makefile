.PHONY: precommit precommit_fix

precommit:
	docker build -t magma_readmes .
	docker-compose run readmes markdownlint --ignore proposals/p010_vendor_neutral_dp.md . && echo PASSED

precommit_fix:
	docker build -t magma_readmes .
	docker-compose run readmes markdownlint --fix --ignore proposals/p010_vendor_neutral_dp.md . && echo PASSED
