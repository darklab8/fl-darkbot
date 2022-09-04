from fastapi import APIRouter
from fastapi import Query, Path, Body
from starlette.requests import Request
from fastapi import Depends
from utils.rest_api.message import MessageOk

# from . import actions
# from . import storage
from . import schemas
from typing import Union


router = APIRouter(
    prefix="/msg",
    tags=["items"],
)

# query_default_values = actions.ActionRegisterChannel.query_factory(channel_id=0)


# @router.post("/test", response_model=MessageOk)
# async def register_channel(
#     query: actions.ActionRegisterChannel.query_factory = Body(),
# ):

#     return MessageOk()


# @router.delete("/test", response_model=MessageOk)
# async def delete_channel(
#     query: actions.ActionDeleteChannel.query_factory = Body(),
# ):
#     return MessageOk()


@router.get("/ping")
async def ping(request: Request):
    return MessageOk()


# @routes.get("/guilds")
# async def get_guilds(request: BaseRequest):

#     return json_response(
#         {"guilds": [guild.id for guild in Singleton().get_bot().guilds]},
#         status=200,
#         content_type="application/json",
#     )


# @routes.get("/{channel}/message/{message}")
# async def get_message(request: BaseRequest):
#     channel_id = request.match_info["channel"]
#     data = await request.post()
#     message = json.dumps({"msg": request.match_info["message"]}, indent=2)

#     print(channel_id)
#     await Singleton().get_bot().get_channel(int(channel_id)).send(message)
#     print("sent")

#     return json_response(
#         {"status": "ok"},
#         status=200,
#         content_type="application/json",
#     )


# @routes.post("/{channel}/message")
# async def post_message(request: BaseRequest):
#     channel_id = request.match_info["channel"]
#     data = await request.post()
#     message = json.dumps({"msg": data["msg"]}, indent=2)

#     await Singleton().get_bot().get_channel(int(channel_id)).send(message)

#     return json_response(
#         {"status": "ok"},
#         status=200,
#         content_type="application/json",
#     )
