from marshmallow import Schema, fields


class HeartbeatRequestObjectSchema(Schema):
    cbsdId = fields.String(required=True)
    grantId = fields.String(required=True)


class HeartbeatRequestSchema(Schema):
    heartbeatRequest = fields.Nested(HeartbeatRequestObjectSchema, required=True, many=True, unknown='true')
