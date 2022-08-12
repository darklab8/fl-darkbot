import argparse
import subprocess
import secrets
from enum import Enum, auto, EnumMeta
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


class _EnumDirectValueMeta(EnumMeta):
    def __getattribute__(cls, name):
        value = super().__getattribute__(name)
        if isinstance(value, cls):
            value = value.name
        return value


class _EnumGetKey(Enum):
    @classmethod
    @property
    def values(self) -> list[str]:
        values = [a.name for a in self]
        logger.info(f"actions.values={values}")
        return values


class EnumValues(_EnumGetKey, metaclass=_EnumDirectValueMeta):
    pass


class Service(EnumValues):
    scrappy = auto()
    pgadmin = auto()
    shell = auto()
    check = auto()
    listener = auto()


class Action(EnumValues):
    pass


class ScrappyActions(Action):
    test = auto()
    shell = auto()
    run = auto()
    lint = auto()
    migrate = auto()
    manage = auto()


class ListenerActions(Action):
    shell = auto()


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
        self._parser.add_argument("service", type=str, choices=Service.values)
        if actions:
            self._parser.add_argument(
                "action",
                type=str,
                choices=actions,
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
        return self._unread_args

    @property
    def unread_cmd(self) -> str:
        return " " + " ".join(self.unread_args)


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

        scrappy_env = (
            "--env-file ./.env.scrappy.staging"
            if pathlib.Path(".env.scrappy.staging").exists()
            else ""
        )

        logger.debug(f"running action = {(self.args.service, self.args.action)}")
        match (self.args.service, self.args.action):

            case (Service.scrappy, ScrappyActions.test):
                self.run_in_compose(
                    command=ComposeCommands.base.format(
                        service=Service.scrappy,
                        cmd=ShellCommands.test.format(optional_cmd=self.unread_cmd),
                    ),
                    session_id=self.session_id,
                )
            case (Service.scrappy, ScrappyActions.shell):
                self.run_in_compose(
                    command=f"{scrappy_env} {ComposeCommands.shell.format(service=Service.scrappy)}",
                )
            case (Service.scrappy, ScrappyActions.run):
                self.run_in_compose(
                    command=f"{scrappy_env} {ComposeCommands.run}",
                    compose_overrides=["network-override"],
                )
            case (Service.scrappy, ScrappyActions.lint):
                self.run_in_compose(
                    command=ComposeCommands.base.format(
                        service=Service.scrappy,
                        cmd=ShellCommands.lint,
                    ),
                    session_id=self.session_id,
                )
            case (Service.scrappy, ScrappyActions.manage):
                self.run_in_compose(
                    command=ComposeCommands.base.format(
                        service=Service.scrappy,
                        cmd=self.unread_cmd,
                    ),
                    session_id=self.session_id,
                )
            case (Service.scrappy, ScrappyActions.migrate):
                self.run_in_compose(
                    command=ComposeCommands.base.format(
                        service=Service.scrappy,
                        cmd='sh -c "python3 utils/scripts/await_db.py --host=scrappy_db && python3 make.py shell migrate scrappy"',
                    ),
                    session_id=self.session_id,
                )
            case (Service.pgadmin, PgadminActions.run):
                self.run_in_compose(command=ComposeCommands.run)
            case (Service.shell, ShellActions.test):
                print(f"unreadargs={self._parser.unread_args}")
                self.shell(ShellCommands.test.format(optional_cmd=self.unread_cmd))
            case (Service.shell, ShellActions.lint):
                self.shell(ShellCommands.lint)
            case (Service.shell, ShellActions.format):
                self.shell(ShellCommands.format)
            case (Service.shell, ShellActions.makemigrations):
                app = self._parser.add_argument("app", type=str).parse_args().args.app
                max_migration = MigrationFile.get_max_migration(app=app)

                command = ShellCommands.makemigrations.format(
                    app=app, number=int(max_migration.number) + 1
                )
                logger.info(f"command={command}")
                self.shell(command)
            case (Service.shell, ShellActions.migrate):
                parser = (
                    self._parser.add_argument(
                        "app",
                        type=str,
                        choices=Service.values + ["help"],
                        help="app: positional_argument[str] = application to migrate. Choices `scrappy` and others above",
                        nargs="?",
                        default="help",
                    )
                    .add_argument(
                        "migration_id",
                        type=str,
                        default="head",
                        nargs="?",
                        help="migrating destination. default/`head` to latest, "
                        "`zero` to zero, `-1` and `+1` are steps back and forward",
                    )
                    .parse_args()
                )
                args = parser.args

                if args.app == "help":
                    parser._parser.print_help()
                    return

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

            case (Service.shell, ShellActions.check):
                logger.info("pong!")
            case (Service.listener, ScrappyActions.shell):
                self.run_in_compose(
                    command=ComposeCommands.shell.format(service=self.service),
                )
            case _:
                raise Exception("Not registered command for this service")


class ShellCommands:
    lint = 'black --exclude="alembic/.*/*.py|OLD_CODE/|venv/" --check .'
    format = 'black . -exclude="alembic/.*/*.py|OLD_CODE/|venv/"'
    test = "pytest {optional_cmd}"

    makemigrations = (
        'alembic -c {app}/alembic.ini revision --autogenerate -m "{number:0>4}"'
    )
    migrate = ""
    upgrade = 'alembic -c {app}/alembic.ini upgrade "{id}"'
    downgrade = 'alembic -c {app}/alembic.ini downgrade "{id}"'


class ComposeCommands:
    base = "run --rm {service}_base {cmd}"
    shell = 'run --user 0 --rm -v "$(pwd):/code" {service}_base bash'
    run = "up"


def main():
    service = Makefile().service
    match service:
        case Service.scrappy:
            Makefile(actions=ScrappyActions.values).run_action()
        case Service.pgadmin:
            Makefile(actions=PgadminActions.values).run_action()
        case Service.shell:
            Makefile(actions=ShellActions.values).run_action()
        case Service.listener:
            Makefile(actions=ListenerActions.values).run_action()
        case Service.check:
            logger.info("pong!")
        case _:
            raise Exception("not registed service")


if __name__ == "__main__":
    main()
