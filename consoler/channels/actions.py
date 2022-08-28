from configurator.channels import actions as channel_actions
from ..commons.requests import RPCAction


class ActionRegisterChannel(RPCAction):
    action = channel_actions.ActionRegisterChannel

    def __init__(self, query: channel_actions.ActionRegisterChannel.query_factory):
        self.query = query


class ActionDeleteChannel(RPCAction):
    action = channel_actions.ActionDeleteChannel

    def __init__(self, query: channel_actions.ActionDeleteChannel.query_factory):
        self.query = query
