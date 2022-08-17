from celery import shared_task
from .actions import ActionGetAndParseAndSaveBases
from scrappy.core.databases import DatabaseFactory
import scrappy.core.settings as settings
from scrappy.core.logger import base_logger

logger = base_logger.getChild(__name__)


@shared_task
def update_bases(database_name: str = settings.DATABASE_NAME) -> bool:
    ActionGetAndParseAndSaveBases(database=DatabaseFactory(name=database_name))
    logger.info(f"task={'update_bases'.upper()} is done")
    return True
