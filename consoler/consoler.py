import argparse
from types import SimpleNamespace

from configurator.bases import actions
from .core import exceptions
import asyncio
from configurator.bases import rpc as bases_actions
from configurator.channels import rpc as channels_actions
from .commons import actions as common_actions


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
    check_parser.add_argument("--channel_id", type=int, required=True)

    # base
    base_parser = actions_choices.add_parser(
        name="base",
        help="base commands",
    )

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
    base_add_parser.add_argument("--channel_id", type=int, required=True)
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
    base_clear_parser.add_argument("--channel_id", type=int, required=True)

    args = root_parser.parse_args()
    return args


async def run_command(args: SimpleNamespace):

    match args.command:
        case "connect":
            params = {"channel_id": args.channel_id}
            if args.owner_id is not None:
                params["owner_id"] = args.owner_id
            if args.owner_name is not None:
                params["owner_name"] = args.owner_name

            action = channels_actions.ActionRegisterChannel
            await action(query=action.action.query_factory(**params)).run()
        case "disconnect":
            params = {"channel_id": args.channel_id}

            action = channels_actions.ActionDeleteChannel
            await action(query=action.action.query_factory(**params)).run()
        case "base":
            match args.action:
                case "add":
                    action = base_actions.ActionRegisterBase
                    await action(
                        query=action.action.query_factory(
                            channel_id=args.channel_id,
                            base_tags=args.base_tags,
                        )
                    ).run()
                case "clear":
                    action = base_actions.ActionDeleteBases
                    await action(
                        query=action.action.query_factory(
                            channel_id=args.channel_id,
                        )
                    ).run()
                case _:
                    raise exceptions.NotRegisteredCommand()
        case "check":
            print(await common_actions.ActionPingConfig().run())
        case _:
            raise exceptions.NotRegisteredCommand()


def main(args=None):
    args = process_cli()
    asyncio.run(run_command(args))
    print(f"executed command with arguments = {repr(args)}")
