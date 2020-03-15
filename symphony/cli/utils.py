#!/usr/bin/env python3

from zipfile import ZIP_DEFLATED, ZipFile


def extract_zip(input_zip_filepath):
    with ZipFile(input_zip_filepath) as input_zip:
        return {name: input_zip.read(name) for name in input_zip.namelist()}


def archive_zip(output_zip_filepath, zip_contents):
    with ZipFile(output_zip_filepath, "w", ZIP_DEFLATED) as output_zip:
        for name, content in zip_contents.items():
            output_zip.writestr(name, content)
