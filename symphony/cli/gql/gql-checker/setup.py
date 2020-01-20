import os
from setuptools import setup, find_packages


base_dir = os.path.dirname(__file__)

about = {}
with open(os.path.join(base_dir, "gql_checker", "__about__.py")) as f:
    exec(f.read(), about)

with open(os.path.join(base_dir, "README.rst")) as f:
    long_description = f.read()


setup(
    name=about["__title__"],
    version=about["__version__"],

    description=about["__summary__"],
    long_description=long_description,
    license=about["__license__"],
    url=about["__uri__"],
    author=about["__author__"],
    author_email=about["__email__"],

    packages=find_packages(exclude=["tests", "tests.*"]),
    zip_safe=False,

    install_requires=[
        "pycodestyle"
    ],

    tests_require=[
        "pytest",
        "flake8",
        "pycodestyle",
        "pylama"
    ],

    py_modules=['gql_checker'],
    entry_points={
        'flake8.extension': [
            'GQL = gql_checker.flake8_linter:Linter',
        ],
        'pylama.linter': [
            'gql_checker = gql_checker.pylama_linter:Linter'
        ]
    },

    classifiers=[
        "Intended Audience :: Developers",
        "Development Status :: 4 - Beta",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python",
        "Programming Language :: Python :: 2",
        "Programming Language :: Python :: 3",
        (
            "License :: OSI Approved :: "
            "GNU Lesser General Public License v3 (LGPLv3)"
        ),
        "Topic :: Software Development :: Libraries :: Python Modules",
        "Topic :: Software Development :: Quality Assurance",
        "Operating System :: OS Independent"
    ]
)
