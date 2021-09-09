from marshmallow import Schema, fields


class GrantRequestObjectSchema(Schema):
    cbsdId = fields.String(required=True)


class GrantRequestSchema(Schema):
    grantRequest = fields.Nested(GrantRequestObjectSchema, required=True, many=True, unknown='true')
