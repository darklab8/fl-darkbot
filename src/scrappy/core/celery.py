from celery import Celery
from celery.schedules import crontab
import os
from scrappy.players.tasks import update_players
from . import settings as settings

app = Celery(
    "core",
    broker=settings.CELERY_BROKER,
    backend=settings.CELERY_BACKEND,
)


@app.on_after_configure.connect
def setup_periodic_tasks(sender, **kwargs):
    # Calls test('world') every 30 seconds
    sender.add_periodic_task(30.0, update_players.s(), expires=10)
