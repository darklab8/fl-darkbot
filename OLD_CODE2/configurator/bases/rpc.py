from ..core.rpc import RPCAction
from . import actions


class ActionRegisterBase(RPCAction):
    action = actions.ActionRegisterBase
    # override for typing only
    query_factory = actions.ActionRegisterBase.query_factory
    response_factory = actions.ActionRegisterBase.response_factory


class ActionDeleteChannel(RPCAction):
    action = actions.ActionDeleteBases
    query_factory = actions.ActionDeleteBases.query_factory
    response_factory = actions.ActionDeleteBases.response_factory


class ActionGetBases(RPCAction):
    action = actions.ActionGetBases
    query_factory = actions.ActionGetBases.query_factory
    response_factory = actions.ActionGetBases.response_factory
