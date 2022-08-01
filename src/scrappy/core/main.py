from fastapi import FastAPI

import scrappy.players.views as player_views


def app_factory():
    app = FastAPI()

    app.include_router(player_views.router)

    @app.get("/")
    def get_ping():
        return {"message": "pong!"}

    return app


print(__name__)
if "main" in __name__:
    app = app_factory()
