from marshmallow import Schema, fields


class DeregistrationRequestObjectSchema(Schema):
    cbsdId = fields.String(required=True)


class DeregistrationRequestSchema(Schema):
    deregistrationRequest = fields.Nested(DeregistrationRequestObjectSchema, required=True, many=True, unknown='true')
