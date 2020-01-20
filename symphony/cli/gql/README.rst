GQL
===

This is a GraphQL client for Python. Plays nicely with ``graphene``,
``graphql-core``, ``graphql-js`` and any other GraphQL implementation
compatible with the spec.

GQL architecture is inspired by ``React-Relay`` and ``Apollo-Client``.

|travis| |pypi| |coveralls|

Installation
------------

::

    $ pip install gql

Usage
-----

The example below shows how you can execute queries against a local
schema.

.. code:: python

    from gql import gql, Client

    client = Client(schema=schema)
    query = gql('''
    {
      hello
    }
    ''')

    client.execute(query)

License
-------

`MIT
License <https://github.com/graphql-python/gql/blob/master/LICENSE>`__

.. |travis| image:: https://img.shields.io/travis/graphql-python/gql.svg?style=flat
   :target: https://travis-ci.org/graphql-python/gql
.. |pypi| image:: https://img.shields.io/pypi/v/gql.svg?style=flat
   :target: https://pypi.python.org/pypi/gql
.. |coveralls| image:: https://coveralls.io/repos/graphql-python/gql/badge.svg?branch=master&service=github
   :target: https://coveralls.io/github/graphql-python/gql?branch=master
