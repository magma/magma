#!/usr/bin/env bash

rm -r -f ./dist ./build ./pyinventory.egg-info/
python3 -m pip install --user --upgrade setuptools wheel
python3 setup.py sdist bdist_wheel
sudo rm -r -f ./build
