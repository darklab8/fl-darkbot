from celery import shared_task
from .actions import ActionGetAndParseAndSavePlayers
from scrappy.core.databases import DatabaseFactory
import scrappy.core.settings as settings


@shared_task
def update_players(database_name: str = settings.DATABASE_NAME):
    ActionGetAndParseAndSavePlayers(database=DatabaseFactory(name=database_name))
    return True
