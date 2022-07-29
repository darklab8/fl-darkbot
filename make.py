import argparse
import os
import secrets
from enum import Enum, auto, EnumMeta
from types import SimpleNamespace

class _EnumDirectValueMeta(EnumMeta):
    def __getattribute__(cls, name):
        value = super().__getattribute__(name)
        if isinstance(value, cls):
            value = value.name
        return value

class _EnumGetKey(Enum):
    @classmethod
    def get_keys(cls):
        return [e.name for e in cls]

class EnumWithValues(_EnumGetKey, metaclass=_EnumDirectValueMeta):
    pass

class Actions(EnumWithValues):
    test = auto()
    shell = auto()
    run = auto()
    lint = auto()

class Services(EnumWithValues):
    scrappy = auto()

def shell(cmd):
    exit(os.system(cmd))

class Parser:
    def __init__(self):
        self._parser = argparse.ArgumentParser()
        self._parser.add_argument('--job_id', type=str, default=secrets.token_hex(4),
            help="optional parameter random by default, ensures to run docker-compose with random -p parameter for no conflicts in parallel runs",
        )
        self._parser.add_argument('service', type=str, choices=Services.get_keys())

    def parse_all(self):
        self._args = self._parser.parse_args()
        return self

    def parse_service_only(self):
        args, argv = self._parser.parse_known_args()
        self._args = args
        return self

    @property
    def args(self) -> SimpleNamespace:
        return self._args

    def registher_actions(self, *actions):
        self._parser.add_argument('action', type=str, choices=actions)
        return self

class CommandExecutor:
    def __init__(self, parser: Parser):
        self._parser = parser

    @property
    def args(self):
        return self._parser.args

    def run_inside_container(self, command):        
        return_code = os.system(
            f"docker-compose -f docker-compose.{self.args.service}.yml -p {self.args.job_id} build && "
            f"docker-compose -f docker-compose.{self.args.service}.yml  -p {self.args.job_id}"
            f" {command}"
        )
        os.system(
            f"docker-compose -f docker-compose.{self.args.service}.yml  -p {self.args.job_id} down"
        )
        print(return_code)
        if return_code != 0:
            raise Exception(f"non zero returned code={return_code}")

class CommonCommands:
    test = "run --rm service_base pytest"
    shell = "run --rm service_base /bin/bash"
    run = "up"
    lint = "run --rm service_base black --check ."

def main():
    args = Parser().parse_service_only().args

    match args.service:
        case Services.scrappy:
            parser = Parser().registher_actions(Actions.test, Actions.shell, Actions.run, Actions.lint).parse_all()
            executor = CommandExecutor(parser)

    match (parser.args.service, parser.args.action):
        case (Services.scrappy, Actions.test):
            executor.run_inside_container(CommonCommands.test)
        case (Services.scrappy, Actions.shell):
            executor.run_inside_container(CommonCommands.shell)
        case (Services.scrappy, Actions.run):
            executor.run_inside_container(CommonCommands.run)
        case (Services.scrappy, Actions.lint):
            executor.run_inside_container(CommonCommands.lint)
        case _:
            raise Exception("Not registered command for this service")

if __name__=="__main__":
    main()