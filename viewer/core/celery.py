from celery import Celery
from scrappy.players.tasks import update_players
from scrappy.bases.tasks import update_bases
from . import settings as settings
from pydantic import BaseModel

app = Celery(
    "core",
    broker=settings.CELERY_BROKER,
    backend=settings.CELERY_BACKEND,
)


class Time(BaseModel):
    seconds: float


@app.on_after_configure.connect
def setup_periodic_tasks(sender: Celery, **kwargs: dict[str, str]) -> None:
    pass
    # sender.add_periodic_task(Time(seconds=30.0).seconds, update_players.s(), expires=10)
    # sender.add_periodic_task(Time(seconds=30.0).seconds, update_bases.s(), expires=10)
