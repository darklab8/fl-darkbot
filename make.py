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


class Service(Enum):
    scrappy = auto()
    pgadmin = auto()
    shell = auto()
    check = auto()


class Action(Enum):
    pass

    @classmethod
    @property
    def values(self) -> list[str]:
        values = [a.name for a in self]
        logger.info(f"actions.values={values}")
        return values


class ScrappyActions(Action):
    test = auto()
    shell = auto()
    run = auto()
    lint = auto()


class PgadminActions(Action):
    run = auto()


class ShellActions(Action):
    test = auto()
    lint = auto()
    format = auto()
    makemigrations = auto()
    migrate = auto()
    check = auto()


@dataclass
class EnumType:
    name: str
    value: int


class Parser:
    def __init__(self, actions: list[str] = []):
        self._parser = argparse.ArgumentParser()
        self._parser.add_argument(
            "--session_id",
            type=str,
            default=secrets.token_hex(4),
            help="ensures to run docker-compose with persistent random -p parameter for no conflicts in parallel runs",
        )
        self._parser.add_argument(
            "service", type=str, choices=[service.name for service in Service]
        )
        if actions:
            self._parser.add_argument(
                "action",
                type=str,
                choices=[action for action in actions],
            )

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
        return self._args


class MigrationFile:
    fullname = ""
    number = ""
    id = ""

    def __init__(self, filename: str):
        logger.debug(f"MigrationFile.filename={filename}")
        found = re.search("([0-9]+)_([0-9a-z]+)\.py", filename)
        if found is None:
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

    def __init__(self, actions=[]):
        self._parser = self.parser_cls(actions=actions)

    @property
    def args(self):
        return self._parser.args

    @property
    def service(self) -> str:
        return self.args.service

    @property
    def session_id(self) -> str:
        return self.args.session_id

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

        scrappy_env = (
            "--env-file ./.env.scrappy.staging"
            if pathlib.Path(".env.scrappy.staging").exists()
            else ""
        )

        logger.debug(f"running action = {(self.args.service, self.args.action)}")
        match (self.args.service, self.args.action):

            case (Service.scrappy.name, ScrappyActions.test.name):
                self.run_in_compose(
                    command=ComposeCommands.test,
                    session_id=self.session_id,
                )
            case (Service.scrappy.name, ScrappyActions.shell.name):
                self.run_in_compose(
                    command=f"{scrappy_env} {ComposeCommands.shell}",
                )
            case (Service.scrappy.name, ScrappyActions.run.name):
                self.run_in_compose(
                    command=f"{scrappy_env} {ComposeCommands.run}",
                    compose_overrides=[f"{Service.scrappy.name}-network"],
                )
            case (Service.scrappy.name, ScrappyActions.lint.name):
                self.run_in_compose(
                    command=ComposeCommands.lint, session_id=self.session_id
                )
            case (Service.pgadmin.name, PgadminActions.run.name):
                self.run_in_compose(command=ComposeCommands.run)
            case (Service.shell.name, ShellActions.test.name):
                self.shell(ShellCommands.test)
            case (Service.shell.name, ShellActions.lint.name):
                self.shell(ShellCommands.lint)
            case (Service.shell.name, ShellActions.format.name):
                self.shell(ShellCommands.format)
            case (Service.shell.name, ShellActions.makemigrations.name):
                app = self._parser.add_argument("app", type=str).parse_args().args.app
                max_migration = MigrationFile.get_max_migration(app=app)

                command = ShellCommands.makemigrations.format(
                    app=app, number=int(max_migration.number) + 1
                )
                logger.info(f"command={command}")
                self.shell(command)
            case (Service.shell.name, ShellActions.migrate.name):
                args = (
                    self._parser.add_argument("app", type=str)
                    .add_argument("migration_id", type=str, default="head", nargs="?")
                    .parse_args()
                    .args
                )
                migration_id = args.migration_id.replace("zero", "base")

                if "+" in migration_id or "head" == migration_id:
                    self.shell(
                        ShellCommands.upgrade.format(app=args.app, id=migration_id)
                    )
                elif "-" in migration_id or "base" == migration_id:
                    self.shell(
                        ShellCommands.downgrade.format(app=args.app, id=migration_id)
                    )
                else:
                    raise Exception("not registered type of migration_id")

            case (Service.shell.name, ShellActions.check.name):
                logger.info("pong!")
            case _:
                raise Exception("Not registered command for this service")


class ShellCommands:
    lint = 'black --exclude="alembic/.*/*.py" --check .'
    format = "black . --exclude=alembic/*"
    test = "pytest"

    makemigrations = (
        'alembic -c {app}/alembic.ini revision --autogenerate -m "{number:0>4}"'
    )
    migrate = ""
    upgrade = 'alembic -c {app}/alembic.ini upgrade "{id}"'
    downgrade = 'alembic -c {app}/alembic.ini downgrade "{id}"'


class ComposeCommands:
    test = f"run --rm service_base {ShellCommands.test}"
    lint = f"run --rm service_base {ShellCommands.lint}"
    shell = 'run --user 0 --rm -v "$(pwd):/code" service_base bash'
    run = "up"


def main():
    service = Makefile().service
    match service:
        case Service.scrappy.name:
            Makefile(actions=ScrappyActions.values).run_action()
        case Service.pgadmin.name:
            Makefile(actions=PgadminActions.values).run_action()
        case Service.shell.name:
            Makefile(actions=ShellActions.values).run_action()
        case Service.check.name:
            logger.info("pong!")
        case _:
            raise Exception("not registed service")


if __name__ == "__main__":
    main()
