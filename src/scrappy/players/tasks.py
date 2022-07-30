from celery import shared_task
from .actions import ActionGetAndParseAndSavePlayers
from scrappy.core.databases import DatabaseFactory
import scrappy.core.settings as settings


@shared_task
def update_players(database_name: str = settings.DATABASE_NAME):
    with DatabaseFactory(name=database_name).manager_to_get_session() as session:
        ActionGetAndParseAndSavePlayers(session=session)
    return True
