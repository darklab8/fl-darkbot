import argparse
import subprocess
import secrets
from enum import Enum, auto
from types import SimpleNamespace
from dataclasses import dataclass
import logging
import pathlib
import re


def get_logger():
    logger = logging.getLogger("").getChild(__name__)
    logger.setLevel(logging.DEBUG)

    # create console handler with a higher log level
    ch = logging.StreamHandler()
    ch.setLevel(logging.DEBUG)
    formatter = logging.Formatter(
        " - ".join(
            [
                "time:%(asctime)s",
                "level:%(levelname)s",
                "name:%(name)s",
                "msg:%(message)s",
            ]
        )
    )
    ch.setFormatter(formatter)
    logger.addHandler(ch)
    return logger


logger = get_logger()


class AugmentedEnum(Enum):
    @classmethod
    @property
    def names(self) -> list[str]:
        names = [a.name for a in self]
        logger.info(f"actions.names={names}")
        return names

    @classmethod
    @property
    def values(self) -> list[str]:
        values = [a.value for a in self]
        logger.info(f"actions.values={values}")
        return values


class ScrappyActions(AugmentedEnum):
    test = auto()
    shell = auto()
    run = auto()
    lint = auto()
    migrate = auto()
    manage = auto()


class ListenerActions(AugmentedEnum):
    shell = auto()


class DiscorderActions(AugmentedEnum):
    shell = auto()


class ConfiguratorActions(AugmentedEnum):
    shell = auto()
    migrate = auto()


class PgadminActions(AugmentedEnum):
    run = auto()


class ShellActions(AugmentedEnum):
    test = auto()
    lint = auto()
    format = auto()
    makemigrations = auto()
    migrate = auto()
    check = auto()
    install = auto()


class Services(AugmentedEnum):
    scrappy = ScrappyActions
    pgadmin = PgadminActions
    shell = ShellActions
    listener = ListenerActions
    discorder = DiscorderActions
    configurator = ConfiguratorActions


@dataclass
class EnumType:
    name: str
    value: int


class Parser:
    def __init__(self, _parser=None):
        self._parser = argparse.ArgumentParser() if _parser is None else _parser
        self._parser.add_argument(
            "--session_id",
            type=str,
            default=secrets.token_hex(4),
            help="ensures to run docker-compose with persistent random -p parameter for no conflicts in parallel runs",
        )

    def register_group(self, *args, name, **kwargs):
        group = self._parser.add_subparsers(*args, dest=name, required=True, **kwargs)
        self.group = group
        return self

    def add_group_choice(self, *args, name, **kwargs):
        setattr(
            self, name, Parser(_parser=self.group.add_parser(*args, name, **kwargs))
        )
        return self

    def add_argument(self, *args, **kwargs):
        self._parser.add_argument(*args, **kwargs)
        return self

    def parse_known_args(self, args=None):
        self._args, self._unread_args = self._parser.parse_known_args(args=args)
        return self

    def parse_args(self, args=None):
        self._args = self._parser.parse_args(args=args)
        return self

    @property
    def args(self) -> SimpleNamespace:
        if not hasattr(self, "_args"):
            self.parse_known_args()
        return self._args

    @property
    def unread_args(self) -> SimpleNamespace:
        if not hasattr(self, "_unread_args"):
            self.parse_known_args()
        return self._unread_args

    @property
    def unread_cmd(self) -> str:
        return " " + " ".join(self.unread_args)


class MigrationFile:
    fullname = ""
    number = "-1"
    id = ""

    def __init__(self, filename: str):
        logger.debug(f"MigrationFile.filename={filename}")

        if filename is None:
            return

        found = re.search("([0-9]+)_([0-9a-z]+)\.py", filename)
        if not found:
            return

        self.fullname = found.group(0)
        self.number = found.group(1)
        self.id = found.group(2)

    def __repr__(self):
        return f"{self.__class__.__name__}(id={self.number}, name={self.id})"

    def __bool__(self):
        return self.fullname == ""

    @classmethod
    def get_max_migration(cls, app) -> "MigrationFile":
        app = app
        path = pathlib.Path(".") / app / "alembic" / "versions"

        migrations: list[MigrationFile] = [
            MigrationFile(file.name)
            for file in path.iterdir()
            if not MigrationFile(file.name)
        ]

        if not migrations:
            return MigrationFile(None)

        max_migration: MigrationFile = max(
            migrations, key=lambda migration: int(migration.number)
        )
        return max_migration


class Makefile:
    parser_cls = Parser

    def __init__(self, services: Enum):
        self._parser = self.parser_cls()

        self._parser.register_group(name="service")
        for service in services:
            self._parser.add_group_choice(name=service.name)
            service_parser: Parser = getattr(self._parser, service.name)

            service_parser.register_group(name="action")
            for action in service.value:
                service_parser.add_group_choice(name=action.name)
                action_parser: Parser = getattr(service_parser, action.name)

                match (service.name, action.name):
                    case (Services.shell.name, ShellActions.migrate.name):
                        action_parser.add_argument(
                            "app",
                            type=str,
                            choices=[Services.scrappy.name, Services.configurator.name],
                            help="app: positional_argument[str] = application to migrate. Choices `scrappy` and others above",
                        )
                        action_parser.add_argument(
                            "migration_id",
                            type=str,
                            default="head",
                            nargs="?",
                            help="migrating destination. default/`head` to latest, "
                            "`zero` to zero, `-1` and `+1` are steps back and forward",
                        )
                    case (Services.shell.name, ShellActions.install.name):
                        action_parser.add_argument(
                            "--environment",
                            type=str,
                            choices=["dev", "prod"],
                            help="building with dev or prod set of deps",
                            required=True,
                        )
                        action_parser.add_argument(
                            "--app",
                            type=str,
                            required=True,
                        )

    @property
    def parser(self):
        return self._parser

    @property
    def args(self):
        return self._parser.args

    @property
    def service(self) -> str:
        return self.args.service

    @property
    def session_id(self) -> str:
        return self.args.session_id

    @property
    def unread_cmd(self) -> str:
        return self._parser.unread_cmd

    def run_in_compose(
        self, command, session_id=None, compose_file=None, compose_overrides=[]
    ):
        session_id = f"-p {self.service}" if session_id is None else f"-p {session_id}"
        compose_file = (
            f"-f docker-compose.{self.service}.yml"
            if compose_file is None
            else f"-f docker-compose.{compose_file}.yml"
        )
        compose_overrides = " ".join(
            [f"-f docker-compose.{override}.yml" for override in compose_overrides]
        )

        logger.debug(f"compose_file={compose_file}")
        logger.debug(f"compose_overrides={compose_overrides}")
        logger.debug(f"session_id={session_id}")
        main_command = (
            f"docker-compose {compose_file} {compose_overrides} {session_id} build && "
            f"docker-compose {compose_file} {compose_overrides} {session_id}"
            f" {command}"
        )
        eixiting_command = (
            f"docker-compose {compose_file} {compose_overrides} {session_id} down"
        )

        logger.info(f"main_command={main_command}")
        logger.info(f"eixiting_command={eixiting_command}")
        try:
            subprocess.run(main_command, shell=True, check=True)
        finally:
            subprocess.run(eixiting_command, shell=True, check=True)

    def shell(self, command):
        logger.info(f"command={command}")
        subprocess.run(command, shell=True, check=True)

    def run_action(self, args=None):

        self._parser.parse_known_args(args=args)

        staging_env = (
            "--env-file ./.env.staging" if pathlib.Path(".env.staging").exists() else ""
        )

        logger.debug(
            f"running command with args={self.args}, unread_args={self._parser.unread_args}"
        )

        match (self.args.service, self.args.action):
            case (Services.shell.name, ShellActions.test.name):
                self.shell(ShellCommands.test.format(optional_cmd=self.unread_cmd))
            case (Services.shell.name, ShellActions.lint.name):
                self.shell(ShellCommands.lint)
            case (Services.shell.name, ShellActions.format.name):
                self.shell(ShellCommands.format)
            case (Services.shell.name, ShellActions.makemigrations.name):
                app = self._parser.add_argument("app", type=str).parse_args().args.app
                max_migration = MigrationFile.get_max_migration(app=app)

                command = ShellCommands.makemigrations.format(
                    app=app, number=int(max_migration.number) + 1
                )
                logger.info(f"command={command}")
                self.shell(command)
            case (Services.shell.name, ShellActions.migrate.name):
                migration_id = self.args.migration_id.replace("zero", "base")

                if "+" in migration_id or "head" == migration_id:
                    self.shell(
                        ShellCommands.upgrade.format(app=self.args.app, id=migration_id)
                    )
                elif "-" in migration_id or "base" == migration_id:
                    self.shell(
                        ShellCommands.downgrade.format(
                            app=self.args.app, id=migration_id
                        )
                    )
                else:
                    raise Exception("not registered type of migration_id")

            case (Services.shell.name, ShellActions.check.name):
                logger.info("pong!")
            case (Services.shell.name, ShellActions.install.name):
                env = self.args.environment
                app = self.args.app
                base_cmd = "pip install --no-cache-dir"
                shared = "-r {app}/requirements.txt"
                dev_only = "-r {app}/requirements.dev.txt"
                end_cmd = "-c {app}/constraints.txt"
                match env:
                    case "dev":
                        self.shell(
                            " ".join([base_cmd, shared, dev_only, end_cmd]).format(
                                app=app
                            )
                        )
                    case "prod":
                        self.shell(
                            " ".join([base_cmd, shared, end_cmd]).format(app=app)
                        )
            case (service, "test"):
                self.run_in_compose(
                    command=ComposeCommands.base.format(
                        service=self.service,
                        cmd=ShellCommands.test.format(optional_cmd=self.unread_cmd),
                    ),
                    session_id=self.session_id,
                )
            case (service, "shell"):
                self.run_in_compose(
                    command=f"{staging_env} {ComposeCommands.shell.format(service=self.service)}",
                )
            case (service, "run"):
                self.run_in_compose(
                    command=f"{staging_env} {ComposeCommands.run}",
                    compose_overrides=["network-override"],
                )
            case (service, "lint"):
                self.run_in_compose(
                    command=ComposeCommands.base.format(
                        service=self.service,
                        cmd=ShellCommands.lint,
                    ),
                    session_id=self.session_id,
                )
            case (service, "manage"):
                self.run_in_compose(
                    command=ComposeCommands.base.format(
                        service=self.service,
                        cmd=self.unread_cmd,
                    ),
                    session_id=self.session_id,
                )
            case (service, "migrate"):
                self.run_in_compose(
                    command=ComposeCommands.base.format(
                        service=self.service,
                        cmd='sh -c "python3 utils/scripts/await_db.py '
                        f'--host={self.service}_db && python3 make.py shell migrate {self.service}"',
                    ),
                    session_id=self.session_id,
                )
            case _:
                raise Exception(
                    "Code that will never run. Not registered command for this service"
                )


class ShellCommands:
    lint = 'black --exclude="alembic/.*/*.py|OLD_CODE/|venv/" --check .'
    format = 'black . -exclude="alembic/.*/*.py|OLD_CODE/|venv/"'
    test = "coverage run -m pytest --cov=. --cov-report=xml {optional_cmd} "

    makemigrations = (
        'alembic -c {app}/alembic.ini revision --autogenerate -m "{number:0>4}"'
    )
    migrate = ""
    upgrade = 'alembic -c {app}/alembic.ini upgrade "{id}"'
    downgrade = 'alembic -c {app}/alembic.ini downgrade "{id}"'


class ComposeCommands:
    base = "run --rm {service}_base {cmd}"
    shell = 'run --user 0 --rm -v "$(pwd):/code" {service}_base sh'
    run = "up"


def main():
    Makefile(services=Services).run_action()


if __name__ == "__main__":
    main()
