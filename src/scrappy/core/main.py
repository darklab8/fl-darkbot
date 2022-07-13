from fastapi import FastAPI

import scrappy.players as players


def app_factory():
    app = FastAPI()

    app.include_router(players.views.router)

    @app.get("/")
    async def get_ping():
        return {"message": "pong!"}

    return app


print(__name__)
if "main" in __name__:
    app = app_factory()
