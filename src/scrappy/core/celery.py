from celery import Celery
from celery.schedules import crontab
import os
from scrappy.players.tasks import update_players

app = Celery(
    "core",
    broker=os.environ.get("SCRAPPY_CELERY_BROKER", "redis://redis:6379/0"),
    backend=os.environ.get("SCRAPPY_CELERY_BROKER", "redis://redis:6379/0"),
)


@app.on_after_configure.connect
def setup_periodic_tasks(sender, **kwargs):
    # Calls test('world') every 30 seconds
    sender.add_periodic_task(30.0, update_players.s(), expires=10)
