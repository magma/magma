from marshmallow import Schema, fields


class RegistrationRequestObjectSchema(Schema):
    fccId = fields.String(required=True)
    cbsdSerialNumber = fields.String(required=True)


class RegistrationRequestSchema(Schema):
    registrationRequest = fields.Nested(RegistrationRequestObjectSchema, required=True, many=True, unknown='true')
