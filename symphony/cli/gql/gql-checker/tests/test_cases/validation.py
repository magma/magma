from gql import gql


gql('''
    query NestedQueryWithFragment {
      hero {
        ...NameAndAppearances
        friends {
          ...NameAndAppearances
          friends {
            ...NameAndAppearances
          }
        }
      }
    }
    fragment NameAndAppearances on Character {
      name
      appearsIn
    }
''')

gql('''
    query HeroSpaceshipQuery {
      hero {
        favoriteSpaceship
      }
    }
''') # GQL101: Cannot query field "favoriteSpaceship" on type "Character".

gql('''
    query HeroNoFieldsQuery {
      hero
    }
''') # GQL101: Field "hero" of type "Character" must have a sub selection.


gql('''
    query HeroFieldsOnScalarQuery {
      hero {
        name {
          firstCharacterOfName
        }
      }
    }
''') # GQL101: Field "name" of type "String" must not have a sub selection.


gql('''
    query DroidFieldOnCharacter {
      hero {
        name
        primaryFunction
      }
    }
''') # GQL101: Cannot query field "primaryFunction" on type "Character". However, this field exists on "Droid". Perhaps you meant to use an inline fragment?

gql('''
    query DroidFieldInFragment {
      hero {
        name
        ...DroidFields
      }
    }
    fragment DroidFields on Droid {
      primaryFunction
    }
''')

gql('''
    query DroidFieldInFragment {
      hero {
        name
        ... on Droid {
          primaryFunction
        }
      }
    }
''')
