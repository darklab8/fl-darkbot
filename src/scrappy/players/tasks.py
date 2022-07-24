from celery import shared_task
from .actions import ActionGetAndParseAndSavePlayers
from scrappy.core.databases import default


@shared_task
def update_players():
    with default.manager_to_get_session() as session:
        ActionGetAndParseAndSavePlayers(session=session)
    return True
