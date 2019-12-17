from setuptools import setup, find_packages

setup(
    name="magma-manipulator",
    version="0.1",
    packages=find_packages(),
    install_requires=['jsonschema==3.2.0',
                      'kubernetes==10.0.1',
                      'paramiko==2.6.0',
                      'requests==2.22.0'],

    entry_points={
        'console_scripts': [
            'magma-manipulator = magma_manipulator.main:main',
        ],
    },
)
