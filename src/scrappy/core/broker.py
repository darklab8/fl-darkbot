from celery import Celery
from celery.schedules import crontab
import os

app = Celery(
    "core",
    broker=os.environ.get("SCRAPPY_CELERY_BROKER", "redis://redis:6379/0"),
)


@app.on_after_configure.connect
def setup_periodic_tasks(sender, **kwargs):
    # Calls test('hello') every 10 seconds.
    sender.add_periodic_task(10.0, test.s("hello"), name="add every 10")

    # Calls test('world') every 30 seconds
    sender.add_periodic_task(30.0, test.s("world"), expires=10)

    # Executes every Monday morning at 7:30 a.m.
    sender.add_periodic_task(
        crontab(hour=7, minute=30, day_of_week=1),
        test.s("Happy Mondays!"),
    )


@app.task
def test(arg):
    print(arg)


@app.task
def add(x, y):
    z = x + y
    print(z)
