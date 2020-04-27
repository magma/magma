#!/usr/bin/env python3
from dataclasses import dataclass, field
from typing import Any, List, Mapping, Union, cast

from graphql import (
    GraphQLEnumType,
    GraphQLInputObjectType,
    GraphQLList,
    GraphQLNonNull,
    GraphQLSchema,
    ListTypeNode,
    NonNullTypeNode,
    OperationDefinitionNode,
    TypeInfo,
    TypeInfoVisitor,
    TypeNode,
    Visitor,
    get_operation_ast,
    is_enum_type,
    is_scalar_type,
    parse,
    validate,
    visit,
)
from graphql.validation.rules.no_unused_fragments import NoUnusedFragmentsRule
from graphql.validation.specified_rules import specified_rules


@dataclass
class ParsedField:
    name: str
    type: str
    is_list: bool
    nullable: bool
    default_value: Any = None


@dataclass
class ParsedEnum:
    name: str
    values: Mapping[str, Any]


@dataclass
class ParsedObject:
    name: str
    fields: List[ParsedField] = field(default_factory=list)
    parents: List[str] = field(default_factory=list)
    children: List["ParsedObject"] = field(default_factory=list)
    inputs: List["ParsedObject"] = field(default_factory=list)
    input_enums: List[ParsedEnum] = field(default_factory=list)


@dataclass
class ParsedVariableDefinition:
    name: str
    type: str
    nullable: bool
    is_list: bool
    default_value: Any = None


@dataclass
class ParsedOperation:
    name: str
    type: str
    variables: List[ParsedVariableDefinition] = field(default_factory=list)
    children: List[ParsedObject] = field(default_factory=list)


NodeT = Union[ParsedOperation, ParsedObject]


@dataclass
class ParsedQuery:
    query: str
    objects: List[NodeT] = field(default_factory=list)
    input_objects: List[NodeT] = field(default_factory=list)
    enums: List[ParsedEnum] = field(default_factory=list)
    internal_enums: List[ParsedEnum] = field(default_factory=list)
    internal_inputs: List[ParsedEnum] = field(default_factory=list)
    fragment_objects: List[ParsedObject] = field(default_factory=list)
    used_fragments: List[str] = field(default_factory=list)


class FieldToTypeMatcherVisitor(Visitor):
    def __init__(self, schema: GraphQLSchema, type_info: TypeInfo, query: str):
        self.schema = schema
        self.type_info = type_info
        self.query = query
        self.parsed = ParsedQuery(query=self.query)
        self.dfs_path: List[ParsedObject] = []

    def push(self, obj: NodeT):
        self.dfs_path.append(obj)

    def pull(self) -> NodeT:
        return self.dfs_path.pop()

    @property
    def current(self) -> ParsedObject:
        return self.dfs_path[-1]

    # Document
    def enter_operation_definition(self, node: OperationDefinitionNode, *_args):
        name, operation = node.name, node.operation

        variables = []

        input_objects = []
        for var in node.variable_definitions:
            ptype, nullable, is_list, var_type = self.__variable_type_to_python(
                var.type
            )
            variables.append(
                ParsedVariableDefinition(
                    name=var.variable.name.value,
                    type=ptype,
                    nullable=nullable,
                    is_list=is_list,
                    default_value=var.default_value.value
                    if var.default_value
                    else None,
                )
            )

            var_name = var_type.name.value

            if isinstance(self.schema.type_map[var_name], GraphQLEnumType):
                enum_type = cast(GraphQLEnumType, self.schema.type_map[var_name])
                enum_name = enum_type.name
                if not any(e.name == enum_name for e in self.parsed.enums):
                    parsed_enum = ParsedEnum(
                        name=enum_type.name,
                        values={
                            val_name: val_value.value or val_name
                            for val_name, val_value in enum_type.values.items()
                        },
                    )
                    self.parsed.enums.append(parsed_enum)
            elif isinstance(self.schema.type_map[var_name], GraphQLInputObjectType):
                input_type = cast(
                    GraphQLInputObjectType, self.schema.type_map[var_name]
                )
                obj = self.parse_input_object(input_type)
                input_objects.append(obj)

        for obj in input_objects:
            self.parsed.input_objects.append(obj)

        parsed_op = ParsedOperation(
            name=name.value,
            type=str(operation.value),
            variables=variables,
            children=[ParsedObject(name=f"{name.value}Data")],
        )

        self.parsed.objects.append(parsed_op)  # pylint:disable=no-member
        self.push(parsed_op)
        self.push(parsed_op.children[0])  # pylint:disable=unsubscriptable-object

        return node

    def parse_obj(self, obj_type):
        obj = ParsedObject(name=str(obj_type))
        for field_name, field_value in obj_type.fields.items():
            field, field_obj_type, parsed_enum = self.__parse_field(
                field_name, field_value.type
            )
            if parsed_enum is not None:
                self.parsed.enums.append(parsed_enum)  # pylint:disable=no-member
            obj.fields.append(field)

            if field_obj_type is not None:
                field_obj = self.parse_obj(field_obj_type)
                obj.children.append(field_obj)
        return obj

    def parse_input_object(self, obj_type):
        obj = ParsedObject(name=str(obj_type))
        for field_name, field_value in obj_type.fields.items():
            field, field_obj_type, parsed_enum = self.__parse_field(
                field_name, field_value.type
            )
            if parsed_enum is not None:
                obj.input_enums.append(parsed_enum)
                self.parsed.internal_enums.append(
                    parsed_enum
                )  # pylint:disable=no-member
            obj.fields.append(field)
            if field_obj_type is not None:
                field_obj = self.parse_input_object(field_obj_type)
                obj.inputs.append(field_obj)
                self.parsed.internal_inputs.append(field_obj)
        return obj

    # def enter_selection_set(self, node, *_):
    #     return node

    def leave_selection_set(self, node, *_):
        self.pull()
        return node

    # Fragments

    def enter_fragment_definition(self, node, *_):
        # Same as operation definition
        obj = ParsedObject(name=node.name.value)
        self.parsed.fragment_objects.append(obj)  # pylint:disable=no-member
        self.push(obj)
        return node

    def enter_fragment_spread(self, node, *_):
        self.current.parents.append(node.name.value)
        self.parsed.used_fragments.append(node.name.value)
        return node

    # def enter_inline_fragment(self, node, *_):
    #     return node
    #
    # def leave_inline_fragment(self, node, *_):
    #     return node

    # Field

    def enter_field(self, node, *_):
        name = node.alias.value if node.alias else node.name.value
        graphql_type = self.type_info.get_type()

        field, obj_type, parsed_enum = self.__parse_field(name, graphql_type)
        if parsed_enum is not None:
            self.parsed.enums.append(parsed_enum)  # pylint:disable=no-member
        self.current.fields.append(field)  # TODO: nullables should go to the end

        if obj_type is not None:
            obj = ParsedObject(name=str(obj_type))
            self.current.children.append(obj)
            self.push(obj)

        return node

    def __parse_field(self, name, graphql_type):
        python_type, is_list, nullable, underlying_graphql_type = self.__scalar_type_to_python(
            graphql_type
        )

        parsed_field = ParsedField(
            name=name, type=python_type, is_list=is_list, nullable=nullable
        )
        parsed_enum = None

        if not is_scalar_type(underlying_graphql_type):
            if is_enum_type(underlying_graphql_type):
                enum_type = cast(
                    GraphQLEnumType, self.schema.type_map[underlying_graphql_type.name]
                )
                name = enum_type.name
                if not any(e.name == name for e in self.parsed.enums):
                    parsed_enum = ParsedEnum(
                        name=enum_type.name,
                        values={
                            name: value.value or name
                            for name, value in enum_type.values.items()
                        },
                    )
            else:
                return parsed_field, underlying_graphql_type, parsed_enum

        return parsed_field, None, parsed_enum

    @staticmethod
    def __scalar_type_to_python(scalar):
        nullable = True
        is_list = False
        if isinstance(scalar, GraphQLNonNull):
            nullable = False
            scalar = scalar.of_type

        mapping = {
            "ID": "str",
            "String": "str",
            "Int": "int",
            "Float": "Number",
            "Boolean": "bool",
            "DateTime": "DateTime",
            "Time": "DateTime",
        }

        if isinstance(scalar, GraphQLList):
            scalar = scalar.of_type
            if isinstance(scalar, GraphQLNonNull):
                scalar = scalar.of_type
                nullable = False
            is_list = True
        mapping = mapping.get(str(scalar), str(scalar))

        return mapping, is_list, nullable, scalar

    @staticmethod
    def __variable_type_to_python(var_type: TypeNode):
        nullable = True
        is_list = False
        if isinstance(var_type, NonNullTypeNode):
            nullable = False
            var_type = var_type.type
        if isinstance(var_type, ListTypeNode):
            is_list = True
            var_type = var_type.type
        if isinstance(var_type, NonNullTypeNode):
            nullable = False
            var_type = var_type.type

        mapping = {
            "ID": "str",
            "String": "str",
            "Int": "int",
            "Float": "Number",
            "Boolean": "bool",
            "DateTime": "DateTime",
            "Time": "DateTime",
        }

        mapping = mapping.get(var_type.name.value, var_type.name.value)
        return mapping, nullable, is_list, var_type


class AnonymousQueryError(Exception):
    def __init__(self):
        super().__init__("All queries must be named")


class InvalidQueryError(Exception):
    def __init__(self, errors):
        self.errors = errors
        message = "\n".join(str(err) for err in errors)
        super().__init__(message)


class QueryParser:
    def __init__(self, schema: GraphQLSchema):
        self.schema = schema
        self.__jinja2_env = None

    def parse(
        self,
        query: str,
        full_fragments: str = "",
        should_validate: bool = True,
        is_fragment: bool = False,
    ) -> ParsedQuery:
        query_document_ast = parse("".join([full_fragments, query]))
        document_ast = parse(query)
        if not is_fragment:
            operation = get_operation_ast(document_ast)
            if not operation.name:
                raise AnonymousQueryError()

        if should_validate:
            errors = validate(
                self.schema,
                query_document_ast,
                [rule for rule in specified_rules if rule is not NoUnusedFragmentsRule],
            )
            if errors:
                raise InvalidQueryError(errors)

        type_info = TypeInfo(self.schema)
        visitor = FieldToTypeMatcherVisitor(self.schema, type_info, query)
        visit(document_ast, TypeInfoVisitor(type_info, visitor))
        result = visitor.parsed
        return result
