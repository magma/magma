#!/usr/bin/env python3

from datetime import datetime
from unittest.mock import MagicMock

from graphql import (
    GraphQLArgument,
    GraphQLEnumType,
    GraphQLEnumValue,
    GraphQLField,
    GraphQLInt,
    GraphQLList,
    GraphQLNonNull,
    GraphQLObjectType,
    GraphQLSchema,
    GraphQLString,
)
from graphql_compiler.gql.query_parser import QueryParser
from graphql_compiler.gql.renderer_dataclasses import DataclassesRenderer
from graphql_compiler.gql.utils_codegen import get_enum_filename

from .base_test import BaseTest


class TestRendererDataclasses(BaseTest):
    def test_simple_query(self):
        query = """
            query GetFilm {
              returnOfTheJedi: film(id: "1") {
                title
                director
              }
            }
        """

        parsed = self.swapi_parser.parse(query)
        rendered = self.swapi_dataclass_renderer.render(parsed)

        m = self.load_module(rendered)
        response = m.GetFilm.from_json(
            """
        {
            "data": {
                "returnOfTheJedi": {
                    "title": "Return of the Jedi",
                    "director": "George Lucas"
                }
            }
        }
        """
        )

        assert response

        data = response.data
        assert data.returnOfTheJedi.title == "Return of the Jedi"
        assert data.returnOfTheJedi.director == "George Lucas"

    def test_simple_query_with_variables(self):
        query = """
            query GetFilm($id: ID!) {
              returnOfTheJedi: film(id: $id) {
                title
                director
              }
            }
        """

        parsed = self.swapi_parser.parse(query)
        rendered = self.swapi_dataclass_renderer.render(parsed)

        m = self.load_module(rendered)

        mock_client = MagicMock()
        mock_client.call = MagicMock(
            return_value="""
           {
               "data": {
                   "returnOfTheJedi": {
                       "title": "Return of the Jedi",
                       "director": "George Lucas"
                   }
               }
           }
        """
        )

        result = m.GetFilm.execute(mock_client, "luke")
        assert result
        assert isinstance(result, m.GetFilm.GetFilmData)

        assert result.returnOfTheJedi.title == "Return of the Jedi"
        assert result.returnOfTheJedi.director == "George Lucas"

    def test_simple_query_with_fragment(self):
        query = """
            query GetFilm {
              returnOfTheJedi: film(id: "1") {
                ...FilmFields
                openingCrawl

              }
            }

            fragment FilmFields on Film {
                title
                director
            }
        """

        parsed = self.swapi_parser.parse(query)
        rendered = self.swapi_dataclass_renderer.render(parsed)

        m = self.load_module(rendered)
        response = m.GetFilm.from_json(
            """
        {
            "data": {
                "returnOfTheJedi": {
                    "title": "Return of the Jedi",
                    "director": "George Lucas",
                    "openingCrawl": "la la la"
                }
            }
        }
        """
        )

        assert response

        data = response.data
        assert data.returnOfTheJedi.title == "Return of the Jedi"
        assert data.returnOfTheJedi.director == "George Lucas"
        assert data.returnOfTheJedi.openingCrawl == "la la la"

    def test_simple_query_with_complex_fragment(self):
        query = """
            query GetPerson {
              luke: character(id: "luke") {
                ...CharacterFields
              }
            }

            fragment CharacterFields on Person {
                name

                home: homeworld {
                    name
                }
            }
        """

        parsed = self.swapi_parser.parse(query)
        rendered = self.swapi_dataclass_renderer.render(parsed)

        m = self.load_module(rendered)
        response = m.GetPerson.from_json(
            """
        {
            "data": {
                "luke": {
                    "name": "Luke Skywalker",
                    "home": {
                        "name": "Arakis"
                    }
                }
            }
        }
        """
        )

        assert response

        data = response.data
        assert data.luke.name == "Luke Skywalker"
        assert data.luke.home.name == "Arakis"

    def test_simple_query_with_complex_fragments(self):
        query = """
            fragment PlanetFields on Planet {
              name
              population
              terrains
            }

            fragment CharacterFields on Person {
              name
              home: homeworld {
                ...PlanetFields
              }
            }

            query GetPerson {
              luke: character(id: "luke") {
                ...CharacterFields
              }
            }
        """

        parsed = self.swapi_parser.parse(query)
        rendered = self.swapi_dataclass_renderer.render(parsed)

        m = self.load_module(rendered)
        response = m.GetPerson.from_json(
            """
        {
            "data": {
                "luke": {
                    "name": "Luke Skywalker",
                    "home": {
                        "name": "Arakis",
                        "population": "1,000,000",
                        "terrains": ["Desert"]
                    }
                }
            }
        }
        """
        )

        assert response

        data = response.data
        assert data.luke.name == "Luke Skywalker"
        assert data.luke.home.name == "Arakis"

    def test_simple_query_with_complex_inline_fragment(self):
        query = """
            query GetPerson {
              luke: character(id: "luke") {
                ... on Person {
                  name
                  home: homeworld {
                    name
                  }
                }
              }
            }
        """

        parsed = self.swapi_parser.parse(query)
        rendered = self.swapi_dataclass_renderer.render(parsed)

        m = self.load_module(rendered)
        response = m.GetPerson.from_json(
            """
            {
                "data": {
                    "luke": {
                        "name": "Luke Skywalker",
                        "home": {
                            "name": "Arakis"
                        }
                    }
                }
            }
            """
        )

        assert response

        data = response.data
        assert data.luke.name == "Luke Skywalker"
        assert data.luke.home.name == "Arakis"

    def test_simple_query_with_enums(self):
        query = """
            query MyIssues {
              viewer {
                issues(first: 5) {
                  edges {
                    node {
                      author { login }
                      authorAssociation
                    }
                  }
                }
              }
            }
        """
        parsed = self.github_parser.parse(query)

        rendered_enums = self.github_dataclass_renderer.render_enums(parsed)
        for enum_name, enum_code in rendered_enums.items():
            self.load_module(enum_code, module_name=get_enum_filename(enum_name))

        rendered = self.github_dataclass_renderer.render(parsed)
        m = self.load_module(rendered)

        response = m.MyIssues.from_json(
            """
            {
                "data": {
                    "viewer": {
                        "issues": {
                            "edges": [
                                {
                                    "node": {
                                        "author": { "login": "whatever" },
                                        "authorAssociation": "FIRST_TIMER"
                                    }
                                }
                            ]
                        }
                    }
                }
            }
            """
        )

        assert response

        node = response.data.viewer.issues.edges[0].node
        assert node
        assert node.author.login == "whatever"
        assert node.authorAssociation == m.CommentAuthorAssociation.FIRST_TIMER

    def test_simple_query_with_missing_enums(self):
        query = """
            query MyIssues {
              viewer {
                issues(first: 5) {
                  edges {
                    node {
                      author { login }
                      authorAssociation
                    }
                  }
                }
              }
            }
        """
        parsed = self.github_parser.parse(query)

        rendered_enums = self.github_dataclass_renderer.render_enums(parsed)
        for enum_name, enum_code in rendered_enums.items():
            self.load_module(enum_code, module_name=get_enum_filename(enum_name))

        rendered = self.github_dataclass_renderer.render(parsed)
        m = self.load_module(rendered)

        response = m.MyIssues.from_json(
            """
            {
                "data": {
                    "viewer": {
                        "issues": {
                            "edges": [
                                {
                                    "node": {
                                        "author": { "login": "whatever" },
                                        "authorAssociation": "VALUE_THAT_DOES_NOT_EXIST"
                                    }
                                }
                            ]
                        }
                    }
                }
            }
            """
        )

        assert response

        node = response.data.viewer.issues.edges[0].node
        assert node
        assert node.author.login == "whatever"
        assert node.authorAssociation == m.CommentAuthorAssociation.MISSING_ENUM

    def test_simple_query_with_enums_default_value(self):
        """
            enum LengthUnit {
              METER
              KM
            }

            type Starship {
              id: ID!
              name: String!
              length(unit: LengthUnit = METER): Float
            }

            type Query {
                ship(id: String!): Starship
            }
        """

        length_unit_enum = GraphQLEnumType(
            "LengthUnit",
            {"METER": GraphQLEnumValue("METER"), "KM": GraphQLEnumValue("KM")},
            description="One of the films in the Star Wars Trilogy",
        )

        starship_type = GraphQLObjectType(
            "Starship",
            lambda: {
                "id": GraphQLField(
                    GraphQLNonNull(GraphQLString), description="The id of the ship."
                ),
                "name": GraphQLField(
                    GraphQLString, description="The name of the ship."
                ),
                "length": GraphQLField(
                    GraphQLInt,
                    args={
                        "unit": GraphQLArgument(
                            GraphQLNonNull(length_unit_enum),
                            default_value="METER",
                            description="id of the droid",
                        )
                    },
                ),
            },
        )

        query_type = GraphQLObjectType(
            "Query",
            lambda: {
                "ship": GraphQLField(
                    starship_type,
                    args={
                        "id": GraphQLArgument(
                            GraphQLNonNull(GraphQLString), description="id of the ship"
                        )
                    },
                )
            },
        )

        schema = GraphQLSchema(query_type, types=[length_unit_enum, starship_type])

        query = """
            query GetStarship {
                ship(id: "Enterprise") {
                    id
                    name
                    length(unit: METER)
                }
            }
        """
        query_parser = QueryParser(schema)
        query_renderer = DataclassesRenderer(schema)
        parsed = query_parser.parse(query)

        rendered_enums = query_renderer.render_enums(parsed)
        for rendered_enum in rendered_enums:
            self.load_module(rendered_enum)

        rendered = query_renderer.render(parsed)
        m = self.load_module(rendered)

        response = m.GetStarship.from_json(
            """
            {
                "data": {
                    "ship": {
                        "id": "Enterprise",
                        "name": "Enterprise",
                        "length": 100
                    }
                }
            }
        """
        )

        assert response

        ship = response.data.ship
        assert ship
        assert ship.id == "Enterprise"
        assert ship.name == "Enterprise"
        assert ship.length == 100

    def test_simple_query_with_datetime(self):
        query = """
            query GetFilm($id: ID!) {
              returnOfTheJedi: film(id: $id) {
                title
                director
                releaseDate
              }
            }
        """

        parsed = self.swapi_parser.parse(query)
        rendered = self.swapi_dataclass_renderer.render(parsed)

        m = self.load_module(rendered)

        now = datetime.now()

        mock_client = MagicMock()
        mock_client.call = MagicMock(
            return_value="""
           {
               "data": {
                   "returnOfTheJedi": {
                       "title": "Return of the Jedi",
                       "director": "George Lucas",
                       "releaseDate": "%s"
                   }
               }
           }
        """
            % now.isoformat()
        )

        result = m.GetFilm.execute(mock_client, "luke")
        assert isinstance(result, m.GetFilm.GetFilmData)

        assert result.returnOfTheJedi.title == "Return of the Jedi"
        assert result.returnOfTheJedi.director == "George Lucas"
        assert result.returnOfTheJedi.releaseDate == now

    def test_non_nullable_list(self):

        PersonType = GraphQLObjectType(
            "Person", lambda: {"name": GraphQLField(GraphQLString)}
        )

        schema = GraphQLSchema(
            query=GraphQLObjectType(
                name="RootQueryType",
                fields={
                    "people": GraphQLField(
                        GraphQLList(GraphQLNonNull(PersonType)),
                        resolve=lambda obj, info: {"name": "eran"},
                    )
                },
            )
        )

        query = """
                query GetPeople {
                  people {
                    name
                  }
                }
            """

        parser = QueryParser(schema)
        dataclass_renderer = DataclassesRenderer(schema)

        parsed = parser.parse(query)
        rendered = dataclass_renderer.render(parsed)

        m = self.load_module(rendered)

        mock_client = MagicMock()
        mock_client.call = MagicMock(
            return_value="""
           {
               "data": {
                   "people": [
                      {
                        "name": "eran"
                      },
                      {
                        "name": "eran1"
                      }
                   ]
               }
           }
        """
        )

        result = m.GetPeople.execute(mock_client)
        assert result
        assert isinstance(result, m.GetPeople.GetPeopleData)

        assert len(result.people) == 2
        assert result.people[0].name == "eran"
        assert result.people[1].name == "eran1"
