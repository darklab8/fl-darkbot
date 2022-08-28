from fastapi import FastAPI

import configurator.channels.views as channels_views
import configurator.bases.views as bases_views
from utils.rest_api.message import MessageOk


def app_factory():
    from . import settings

    app = FastAPI()

    app.include_router(channels_views.router)
    app.include_router(bases_views.router)

    @app.get("/")
    def get_ping():
        return dict(MessageOk())

    return app


if "main" in __name__:
    app = app_factory()
