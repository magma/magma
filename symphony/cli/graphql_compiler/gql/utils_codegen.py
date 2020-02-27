#!/usr/bin/env python3
import os
import re


SPACES = " " * 4


def camel_case_to_lower_case(string: str) -> str:
    return re.sub(r"(?<!^)(?=[A-Z])", "_", string).lower()


def get_filename_by_extension(name: str, extension: str) -> str:
    lower_case_name = camel_case_to_lower_case(name)
    if lower_case_name.endswith("".join(["_", extension])):
        return lower_case_name
    return "_".join([lower_case_name, extension])


def get_enum_filename(enum_name: str) -> str:
    return get_filename_by_extension(enum_name, "enum")


def get_input_filename(input_name: str) -> str:
    return get_filename_by_extension(input_name, "input")


class CodeChunk:
    class Block:
        def __init__(self, codegen: "CodeChunk"):
            self.gen = codegen

        def __enter__(self):
            self.gen.indent()
            return self.gen

        def __exit__(self, exc_type, exc_val, exc_tb):
            self.gen.unindent()

    def __init__(self):
        self.lines = []
        self.level = 0

    def indent(self):
        self.level += 1

    def unindent(self):
        if self.level > 0:
            self.level -= 1

    @property
    def indent_string(self):
        return self.level * SPACES

    def write(self, value: str, *args, **kwargs):
        if value != "":
            value = self.indent_string + value
        if args or kwargs:
            value = value.format(*args, **kwargs)

        self.lines.append(value)

    def write_lines(self, lines):
        for line in lines:
            self.lines.append(self.indent_string + line)

    def block(self):
        return self.Block(self)

    def write_block(self, block_header: str, *args, **kwargs):
        self.write(block_header, *args, **kwargs)
        return self.block()

    def __str__(self):
        return os.linesep.join(self.lines)
