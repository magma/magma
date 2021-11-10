from marshmallow import Schema, fields


class RelinquishmentRequestObjectSchema(Schema):
    cbsdId = fields.String(required=True)


class RelinquishmentRequestSchema(Schema):
    relinquishmentRequest = fields.Nested(RelinquishmentRequestObjectSchema, required=True, many=True, unknown='true')
