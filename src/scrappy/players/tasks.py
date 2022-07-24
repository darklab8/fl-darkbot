from celery import shared_task
from .actions import ActionGetAndParseAndSavePlayers
from scrappy.core.databases import default


@shared_task
def test(arg):
    print(arg)


@shared_task
def add(x, y):
    z = x + y
    print(z)
    return z


@shared_task
def update_players():
    with default.manager_to_get_session() as session:
        ActionGetAndParseAndSavePlayers(session=session)
    return True
