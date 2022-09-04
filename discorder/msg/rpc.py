import discord
from ..core.rpc import RPCAction
from utils.rest_api.message import MessageOk
from . import queries
from . import actions


class CreateOrReplaceMessage(RPCAction):
    action = actions.CreateOrReplaceMessage
    # override for typing only
    query_factory = actions.CreateOrReplaceMessage.query_factory
    response_factory = actions.CreateOrReplaceMessage.response_factory


class DeleteMessage(RPCAction):
    action = actions.DeleteMessage
    query_factory = actions.DeleteMessage.query_factory
    response_factory = actions.DeleteMessage.response_factory
