gql-checker
===========

|Build Status|

A `flake8 <http://flake8.readthedocs.org/en/latest/>`__ and
`Pylama <https://github.com/klen/pylama>`__ plugin that checks the
all the static gql calls given a GraphQL schema.

It will not check anything else about the gql calls. Merely that the
GraphQL syntax is correct and it validates against the provided schema.

Warnings
--------

This package adds 3 new flake8 warnings

-  ``GQL100``: The gql query is doesn't match GraphQL syntax
-  ``GQL101``: The gql query have valid syntax but doesn't validate against provided schema

Configuration
-------------

You will want to set the ``gql-introspection-schema`` option to a
file with the json introspection of the schema.


.. |Build Status| image:: https://travis-ci.org/graphql-python/gql-checker.png?branch=master
   :target: https://travis-ci.org/graphql-python/gql-checker
