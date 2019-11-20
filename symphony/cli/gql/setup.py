import sys
from setuptools import setup, find_packages

setup(
    name='gql',
    version='0.1.0',
    description='GraphQL client for Python',
    long_description=open('README.rst').read(),
    url='https://github.com/graphql-python/gql',
    author='Syrus Akbary',
    author_email='me@syrusakbary.com',
    license='MIT',
    classifiers=[
        'Development Status :: 3 - Alpha',
        'Intended Audience :: Developers',
        'Topic :: Software Development :: Libraries',
        'Programming Language :: Python :: 2',
        'Programming Language :: Python :: 2.7',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.3',
        'Programming Language :: Python :: 3.4',
        'Programming Language :: Python :: 3.5',
        'Programming Language :: Python :: Implementation :: PyPy',
    ],
    keywords='api graphql protocol rest relay gql client',
    packages=find_packages(include=["gql*"]),
    install_requires=[
        'six>=1.10.0',
        'graphql-core>=0.5.0',
        'promise>=0.4.0'
    ],
    tests_require=['pytest>=2.7.2', 'mock'],
)
