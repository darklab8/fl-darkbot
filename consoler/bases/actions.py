from configurator.bases import actions as base_actions
from ..commons.requests import RPCAction


class ActionRegisterBase(RPCAction):
    action = base_actions.ActionRegisterBase

    def __init__(self, query: base_actions.ActionRegisterBase.query_factory):
        self.query = query


class ActionDeleteBases(RPCAction):
    action = base_actions.ActionDeleteBases

    def __init__(self, query: base_actions.ActionDeleteBases.query_factory):
        self.query = query
