.PHONY: precommit precommit_fix

FILES_TO_IGNORE = --ignore proposals/p010_vendor_neutral_dp.md \
                  --ignore proposals/p010_subscriber_scaling.md \
                  --ignore proposals/p015_tech_debt_week.md \
                  --ignore proposals/p018_control_network_metrics.md

precommit:
	docker build -t magma_readmes .
	docker compose --compatibility run readmes markdownlint $(FILES_TO_IGNORE) . && echo PASSED

precommit_fix:
	docker build -t magma_readmes .
	docker compose --compatibility run readmes markdownlint --fix $(FILES_TO_IGNORE) . && echo PASSED
