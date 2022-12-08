from fastapi import FastAPI

import scrappy.players.views as player_views
import scrappy.bases.views as base_views


def app_factory() -> FastAPI:
    from . import settings

    app = FastAPI()

    app.include_router(player_views.router)
    app.include_router(base_views.router)

    @app.get("/")
    def get_ping() -> dict[str, str]:
        return {"message": "pong!"}

    return app


if "main" in __name__:

    app = app_factory()
