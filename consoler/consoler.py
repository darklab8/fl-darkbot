import argparse
from types import SimpleNamespace
import aiohttp
from .exceptions import NotRegisteredCommand
from . import settings
import asyncio


def process_cli() -> SimpleNamespace:
    root_parser = argparse.ArgumentParser(description="Process some integers.")

    actions_choices = root_parser.add_subparsers(
        dest="command",
        help="Darkbot help",
        required=True,
    )

    check_parser = actions_choices.add_parser(
        name="check",
        help="check commands",
    )

    base_parser = actions_choices.add_parser(
        name="base",
        help="base commands",
    )

    base_choices = base_parser.add_subparsers(
        dest="action",
        help="base actions",
        required=True,
    )

    base_add_parser = base_choices.add_parser(
        name="add",
        help="adding base for tracking",
    )
    base_add_parser.add_argument(
        "tags",
        metavar="tag",
        type=str,
        nargs="+",
    )
    base_clear_parser = base_choices.add_parser(
        name="c;ear",
        help="clearing settings",
    )

    args = root_parser.parse_args()
    return args


async def config_request(path: str):
    async with aiohttp.ClientSession() as session:
        async with session.get(settings.CONFIGURATOR_API + path) as resp:

            if resp.status < 299:
                print("successful request")

            return resp


async def run_command(args: SimpleNamespace):

    match args.command:
        case "base":
            match args.action:
                case "add":
                    pass
                case _:
                    raise NotRegisteredCommand()
        case "check":
            await config_request(path="/")
        case _:
            raise NotRegisteredCommand()


def main(args=None):
    args = process_cli()
    asyncio.run(run_command(args))
    print(f"executed command with arguments = {repr(args)}")
