import argparse
import os
import secrets
from enum import Enum, auto, EnumMeta

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

random_string = secrets.token_hex(16)
parser = argparse.ArgumentParser()
parser.add_argument('--job_id', type=str, default=secrets.token_hex(4),
    help="optional parameter random by default, ensures to run docker-compose with random -p parameter for no conflicts in parallel runs",
)
parser.add_argument('service', type=str, choices=[Services.scrappy])
parser.add_argument('action', type=str, choices=[Actions.test, Actions.shell, Actions.run, Actions.lint])
args = parser.parse_args()

def run_inside_container(command):
    return_code = os.system(
        f"docker-compose -f docker-compose.{args.service}.yml -p {args.job_id} build && "
        f"docker-compose -f docker-compose.{args.service}.yml  -p {args.job_id}"
        f" {command}"
    )
    os.system(
        f"docker-compose -f docker-compose.{args.service}.yml  -p {args.job_id} down"
    )
    print(return_code)
    if return_code != 0:
        raise Exception(f"non zero returned code={return_code}")

match args.action:
    case Actions.test:
        run_inside_container(command="run --rm service_base pytest")
    case Actions.shell:
        run_inside_container(command="run --rm service_base /bin/bash")
    case Actions.run:
        run_inside_container(command="up")
    case Actions.lint:
        run_inside_container(command="run --rm service_base black --check .")
