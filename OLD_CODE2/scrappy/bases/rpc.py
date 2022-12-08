from ..core.rpc import RPCAction
from . import schemas
from utils.rest_api.methods import RequestMethod


class ActionGetFilteredBases(RPCAction):
    action = None
    query_factory = schemas.BaseQueryParams
    method = RequestMethod.post
    url = "/bases"
    response_factory = schemas.BasesOut
