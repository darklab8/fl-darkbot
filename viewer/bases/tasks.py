from celery import shared_task
import asyncio

# from .actions import ActionGetAndParseAndSaveBases
from scrappy.bases import rpc as rpc_bases_data
from configurator.bases import rpc as rpc_bases_settings
from ..core import settings
from ..core.logger import base_logger

logger = base_logger.getChild(__name__)


@shared_task
def update_all_bases() -> bool:
    result = asyncio.run(
        rpc_bases_settings.ActionGetBases(
            query=rpc_bases_settings.ActionGetBases.query_factory()
        ).run()
    )
    logger.debug(f"update_all_bases.result={result}")


# def render_to_channel(channel_id, tags: list[str]):
# invent some unique ID of a msg.
#

# ActionGetAndParseAndSaveBases(database=DatabaseFactory(name=database_name))
# logger.info(f"task={'update_bases'.upper()} is done")
# return True
