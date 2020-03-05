#!/usr/bin/env python3

from dataclasses import field
from enum import Enum
from functools import partial
from typing import Type


# pyre-ignore
def enum_field(enum_type: Type[Enum]):
    def encode_enum(value: Enum) -> str:
        return value.value

    def decode_enum(t: Type[Enum], value: str) -> Enum:
        return t(value)

    return field(
        metadata={
            "dataclasses_json": {
                "encoder": encode_enum,
                "decoder": partial(decode_enum, enum_type),
            }
        }
    )
