import asyncio

from fastapi import FastAPI
from starlette.requests import Request
from .bot import create_bot, run_bot
from .main import views as main_views
from .msg import views as msg_views


def create_app() -> FastAPI:
    app = FastAPI()

    @app.on_event("startup")
    async def on_startup():
        app.discord_bot = await create_bot()
        asyncio.create_task(run_bot(app.discord_bot))

    app.include_router(main_views.router)
    app.include_router(msg_views.router)

    @app.get("/")
    def root(request: Request):
        return repr(request.app.discord_bot)

    return app
