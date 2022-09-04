from aiohttp.web import RouteTableDef, json_response
from aiohttp.web import BaseRequest
import json
from discorder.bot import Singleton

routes = RouteTableDef()


@routes.get("/")
async def ping(request: BaseRequest):
    return json_response(
        {"ping!": "pong"},
        status=200,
        content_type="application/json",
    )


@routes.get("/guilds")
async def get_guilds(request: BaseRequest):

    return json_response(
        {"guilds": [guild.id for guild in Singleton().get_bot().guilds]},
        status=200,
        content_type="application/json",
    )


@routes.get("/{channel}/message/{message}")
async def get_message(request: BaseRequest):
    channel_id = request.match_info["channel"]
    data = await request.post()
    message = json.dumps({"msg": request.match_info["message"]}, indent=2)

    print(channel_id)
    await Singleton().get_bot().get_channel(int(channel_id)).send(message)
    print("sent")

    return json_response(
        {"status": "ok"},
        status=200,
        content_type="application/json",
    )


@routes.post("/{channel}/message")
async def post_message(request: BaseRequest):
    channel_id = request.match_info["channel"]
    data = await request.post()
    message = json.dumps({"msg": data["msg"]}, indent=2)

    await Singleton().get_bot().get_channel(int(channel_id)).send(message)

    return json_response(
        {"status": "ok"},
        status=200,
        content_type="application/json",
    )
