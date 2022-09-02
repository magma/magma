# Disable W0511: todo warnings
# Disable R0903: Too few public methods
CHECK_CMD_PYLINT := -find . -name '*.py' -exec $(BIN)/pylint --disable=R0903,W0511 {} +;
CHECK_CMD_PYCODESTYLE := -find . -name '*.py' -exec $(BIN)/pycodestyle {} +;
# ’-’ ignores the exit 1 and the script is not terminated by a linter error,
# see https://www.gnu.org/software/make/manual/make.html#Errors
check: buildenv $(BIN)/pylint $(BIN)/pycodestyle
	$(CHECK_CMD_PYCODESTYLE)
	$(CHECK_CMD_PYLINT)
