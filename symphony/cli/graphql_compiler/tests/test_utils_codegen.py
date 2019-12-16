#!/usr/bin/env python3

from .base_test import BaseTest
from fbc.symphony.cli.graphql_compiler.gql.utils_codegen import CodeChunk


class TestRendererDataclasses(BaseTest):
    def test_codegen_write_simple_strings(self):
        gen = CodeChunk()
        gen.write('def sum(a, b):')
        gen.indent()
        gen.write('return a + b')

        code = str(gen)

        m = self.load_module(code)
        assert m.sum(2, 3) == 5

    def test_codegen_write_template_strings_args(self):
        gen = CodeChunk()
        gen.write('def {0}(a, b):', 'sum')
        gen.indent()
        gen.write('return a + b')

        code = str(gen)

        m = self.load_module(code)
        assert m.sum(2, 3) == 5

    def test_codegen_write_template_strings_kwargs(self):
        gen = CodeChunk()
        gen.write('def {method}(a, b):', method='sum')
        gen.indent()
        gen.write('return a + b')

        code = str(gen)

        m = self.load_module(code)
        assert m.sum(2, 3) == 5

    def test_codegen_block(self):
        gen = CodeChunk()
        gen.write('def sum(a, b):')
        with gen.block():
            gen.write('return a + b')

        code = str(gen)

        m = self.load_module(code)
        assert m.sum(2, 3) == 5

    def test_codegen_write_block(self):
        gen = CodeChunk()
        with gen.write_block('def {name}(a, b):', name='sum'):
            gen.write('return a + b')

        code = str(gen)

        m = self.load_module(code)
        assert m.sum(2, 3) == 5

    def test_codegen_write_lines(self):
        lines = [
            '@staticmethod',
            'def sum(a, b):'
            '    return a + b'
        ]
        gen = CodeChunk()
        gen.write('class Math:')
        gen.indent()
        gen.write_lines(lines)

        code = str(gen)

        m = self.load_module(code)
        assert m.Math.sum(2, 3) == 5
