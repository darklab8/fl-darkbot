from fastapi import FastAPI

from .players import views


def app_factory():
    app = FastAPI()

    app.include_router(views.router)

    @app.get("/")
    async def get_ping():
        return {"message": "pong!"}

    return app


if __name__ == "__main__":
    app = app_factory()
