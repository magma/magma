## Build python-precommit Docker image

To build the image, run the following in your host machine.
```bash
cd $MAGMA/lte/gateway/python
./precommit.py -b
```

## Run commands
To use the flak8 linter, run the following.
```bash
cd $MAGMA/lte/gateway/python
./precommit.py --lint -p PATH1 PATH2
OR 
./precommit.py --lint--diff
```

To use the formatting tools ([isort](https://pypi.org/project/isort/), 
[autopep8](https://pypi.org/project/autopep8/)), run the following.
```bash
cd $MAGMA/lte/gateway/python
./precommit.py --format -p PATH1 PATH2
OR 
./precommit.py --format--diff
```
