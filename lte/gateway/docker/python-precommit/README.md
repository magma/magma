## Build python-precommit Docker image

To build python-precommit base image, run the following. 
```bash
# MAGMA_ROOT should be set to repo root
export PATH_TO_DOCKERFILE=$MAGMA_ROOT/lte/gateway/docker/python-precommit/Dockerfile
docker build -t magma/py-lint -f $PATH_TO_DOCKERFILE $MAGMA_ROOT
```

## Run commands
Refer to `requirements.in` in this directory for available packages inside the image. 
Here is an example of running flake8 with this image.
```bash
docker run -it -u 0 -v $MAGMA_ROOT:/code magma/py-lint:latest flake8  lte/gateway/python/precommit.py
```

## How to use `lte/gateway/python/precommit.py`
We have a utility script that wraps all necessary Docker commands with Python.
You should refer to the script for all available commands, but the main ones are as follows.
```bash
cd $MAGMA/lte/gateway/python
# to build the base image
./precommit.py --build

# to run the flake8 linter by specifying paths
./precommit.py --lint -p PATH1 PATH2
# to run the flake8 linter on all modified files in the current commit
./precommit.py --lint--diff

# to run all available formatters by specifying paths
./precommit.py --format -p PATH1 PATH2
# to run all available formatters on all modified files in the current commit
./precommit.py --format--diff
```

## How to add/update Python dependencies via `requirements.in`
`requirements.in` is the file that manages most `Python` package dependencies for this Dockerfile.
If you need to add or update dependencies in `requirements.in`, always run the following commands.
```bash
# Install https://github.com/jazzband/pip-tools if you don't have it already
cd lte/gateway/docker/python-precommit
pip-compile requirements.in  # This should update `requirements.txt`
```