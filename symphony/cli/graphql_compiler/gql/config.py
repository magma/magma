from typing import Type, TypeVar
from dataclasses import dataclass
from dataclasses_json import dataclass_json


ConfigT = TypeVar('ConfigT', bound='ConfigT')

@dataclass_json
@dataclass(frozen=True)
class Config:
    schema: str
    endpoint: str
    documents: str
    custom_header: str = ''

    @classmethod
    def load(cls: Type[ConfigT], filename: str) -> ConfigT:
        with open(filename, 'r') as fin:
            json_str = fin.read()
            return cls.from_json(json_str)  # pylint:disable=no-member

    def save(self, filename, pretty=True):
        with open(filename, 'w') as outfile:
            json_str = self.to_json(indent=2) if pretty else self.to_json()  # pylint:disable=no-member
            outfile.write(json_str)
