import argparse
import os
import secrets
from enum import Enum, auto, EnumMeta
from types import SimpleNamespace
from os import path


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
    pgadmin = auto()


def shell(cmd):
    exit(os.system(cmd))


class Parser:
    def __init__(self):
        self._parser = argparse.ArgumentParser()
        self._parser.add_argument(
            "--job_id",
            type=str,
            default=secrets.token_hex(4),
            help="optional parameter random by default, ensures to run docker-compose with random -p parameter for no conflicts in parallel runs",
        )
        self._parser.add_argument("service", type=str, choices=Services.get_keys())

    @property
    def service(self):
        args, argv = self._parser.parse_known_args()
        return args.service

    def register_actions(self, *actions):
        self._parser.add_argument("action", type=str, choices=actions)
        args = self._parser.parse_args()
        self._args = args
        return self

    @property
    def args(self) -> SimpleNamespace:
        return self._args


class CommandExecutor:
    def __init__(self, parser: Parser):
        self._parser = parser

    @property
    def args(self):
        return self._parser.args

    def run_inside_container(self, command, job_id=None):
        job_id_ = job_id if job_id is not None else self.args.job_id

        main_command = (
            f"docker-compose -f docker-compose.{self.args.service}.yml -p {job_id_} build && "
            f"docker-compose -f docker-compose.{self.args.service}.yml -p {job_id_}"
            f" {command}"
        )
        print(f"running: {main_command}")
        return_code = os.system(main_command)
        eixiting_command = f"docker-compose -f docker-compose.{self.args.service}.yml  -p {job_id_} down"
        print(f"running: {eixiting_command}")
        os.system(eixiting_command)
        print(return_code)
        if return_code != 0:
            raise Exception(f"non zero returned code={return_code}")

    def run(self):
        scrappy_env_file_command = (
            "--env-file ./.env.scrappy.staging"
            if path.exists(".env.scrappy.staging")
            else ""
        )

        match (self.args.service, self.args.action):
            case (Services.scrappy, Actions.test):
                self.run_inside_container(CommonCommands.test)
            case (Services.scrappy, Actions.shell):
                self.run_inside_container(
                    f"{scrappy_env_file_command} run --rm service_shell",
                    job_id="scrappy_shell",
                )
            case (Services.scrappy, Actions.run):
                self.run_inside_container(
                    f"-f docker-compose.scrappy-network.yml {scrappy_env_file_command} up",
                    job_id="scrappy_run",
                )
            case (Services.scrappy, Actions.lint):
                self.run_inside_container(CommonCommands.lint)
            case (Services.pgadmin, Actions.run):
                self.run_inside_container(CommonCommands.run)
            case _:
                raise Exception("Not registered command for this service")


class CommonCommands:
    test = "run --rm service_base pytest"
    run = "up"
    lint = 'run --rm service_base black --exclude="alembic/.*/*.py" --check .'


def main():
    service = Parser().service
    match service:
        case Services.scrappy:
            CommandExecutor(
                parser=Parser().register_actions(
                    Actions.test, Actions.shell, Actions.run, Actions.lint
                )
            ).run()
        case Services.pgadmin:
            CommandExecutor(parser=Parser().register_actions(Actions.run)).run()


if __name__ == "__main__":
    main()
