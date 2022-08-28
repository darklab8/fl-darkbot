import argparse
from types import SimpleNamespace
import aiohttp
from . import exceptions
from . import settings
import asyncio
from typing import Any
from configurator.channels import actions as configurator_channel_actions
from configurator.bases import actions as configurator_base_actions


def process_cli() -> SimpleNamespace:
    root_parser = argparse.ArgumentParser(description="Process some integers.")

    actions_choices = root_parser.add_subparsers(
        dest="command",
        help="Darkbot help",
        required=True,
    )

    # connect
    channel_connect_parser = actions_choices.add_parser(
        name="connect",
        help="",
    )
    channel_connect_parser.add_argument("--channel_id", type=int, required=True)
    channel_connect_parser.add_argument("--owner_id", type=int)
    channel_connect_parser.add_argument("--owner_name", type=str)

    # disconnect
    channel_delete_parser = actions_choices.add_parser(
        name="disconnect",
        help="",
    )
    channel_delete_parser.add_argument("--channel_id", type=int, required=True)

    # check
    check_parser = actions_choices.add_parser(
        name="check",
        help="check commands",
    )

    # base
    base_parser = actions_choices.add_parser(
        name="base",
        help="base commands",
    )
    base_parser.add_argument("--channel_id", type=int, required=True)

    base_choices = base_parser.add_subparsers(
        dest="action",
        help="base actions",
        required=True,
    )

    # base add
    base_add_parser = base_choices.add_parser(
        name="add",
        help="adding base for tracking",
    )
    base_add_parser.add_argument(
        "base_tags",
        metavar="tag",
        type=str,
        nargs="+",
    )
    # base clear
    base_clear_parser = base_choices.add_parser(
        name="clear",
        help="clearing settings",
    )

    args = root_parser.parse_args()
    return args


def is_succesful_request(code: int):
    if code < 299:
        return True
    return False


async def config_request(path: str, method: str, json: dict[str, Any] = {}):
    async with aiohttp.ClientSession() as session:
        match method:
            case "get":
                async with session.get(settings.CONFIGURATOR_API + path) as resp:
                    if is_succesful_request(resp.status):
                        print("successful request")
                    return resp
            case "post":
                async with session.post(
                    settings.CONFIGURATOR_API + path, json=json
                ) as resp:
                    if is_succesful_request(resp.status):
                        print("successful request")
                    return resp
            case "delete":
                async with session.delete(
                    settings.CONFIGURATOR_API + path, json=json
                ) as resp:
                    if is_succesful_request(resp.status):
                        print("successful request")
                    return resp
            case _:
                raise exceptions.NotImplementedMethod()


async def run_command(args: SimpleNamespace):

    match args.command:
        case "connect":
            params = {"channel_id": args.channel_id}
            if args.owner_id is not None:
                params["owner_id"] = args.owner_id
            if args.owner_name is not None:
                params["owner_name"] = args.owner_name

            registry_channel = configurator_channel_actions.ActionRegisterChannel
            await config_request(
                path=registry_channel.url,
                method=registry_channel.method.name,
                json=dict(registry_channel.query_factory(**params)),
            )
        case "disconnect":
            params = {"channel_id": args.channel_id}

            delete_channel = configurator_channel_actions.ActionDeleteChannel
            await config_request(
                path=delete_channel.url,
                method=delete_channel.method.name,
                json=dict(delete_channel.query_factory(**params)),
            )
        case "base":
            match args.action:
                case "add":
                    registry_base = configurator_base_actions.ActionRegisterBase
                    await config_request(
                        path=registry_base.url,
                        method=registry_base.method.name,
                        json=dict(
                            registry_base.query_factory(
                                channel_id=args.channel_id,
                                base_tags=args.base_tags,
                            )
                        ),
                    )
                case "clear":
                    delete_base = configurator_base_actions.ActionDeleteBases
                    await config_request(
                        path=delete_base.url,
                        method=delete_base.method.name,
                        json=dict(
                            delete_base.query_factory(
                                channel_id=args.channel_id,
                            )
                        ),
                    )
                case _:
                    raise exceptions.NotRegisteredCommand()
        case "check":
            await config_request(path="/", method="get")
        case _:
            raise exceptions.NotRegisteredCommand()


def main(args=None):
    args = process_cli()
    asyncio.run(run_command(args))
    print(f"executed command with arguments = {repr(args)}")
