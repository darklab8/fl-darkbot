from celery import Celery
from ..bases import tasks as bases_tasks
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
    bases_tasks.update_all_bases.delay()
    sender.add_periodic_task(
        Time(seconds=30.0).seconds, bases_tasks.update_all_bases.s(), expires=10
    )
