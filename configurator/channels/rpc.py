from ..core.rpc import RPCAction
from . import actions


class ActionRegisterChannel(RPCAction):
    action = actions.ActionRegisterChannel
    # override for typing only
    query_factory = actions.ActionRegisterChannel.query_factory
    response_factory = actions.ActionRegisterChannel.response_factory


class ActionDeleteChannel(RPCAction):
    action = actions.ActionDeleteChannel
    query_factory = actions.ActionDeleteChannel.query_factory
    response_factory = actions.ActionDeleteChannel.response_factory
