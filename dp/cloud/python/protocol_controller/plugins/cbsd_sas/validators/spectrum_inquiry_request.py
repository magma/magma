from marshmallow import Schema, fields


class SpectrumInquiryRequestObjectSchema(Schema):
    cbsdId = fields.String(required=True)


class SpectrumInquiryRequestSchema(Schema):
    spectrumInquiryRequest = fields.Nested(SpectrumInquiryRequestObjectSchema, required=True, many=True, unknown='true')
