from celery import shared_task
import asyncio

# from .actions import ActionGetAndParseAndSaveBases
from scrappy.bases import rpc as rpc_bases_data
from scrappy.bases import schemas as rpc_bases_data_schemas
from configurator.bases import rpc as rpc_bases_settings
from configurator.bases import schemas as rpc_bases_settings_schemas
from ..core import settings
from ..core.logger import base_logger

logger = base_logger.getChild(__name__)


@shared_task
def update_all_bases() -> bool:
    pass

    # for base in base_settings:
    #     render_base.delay(
    #         channel_id=base.channel_id,
    #         tags=base.tags,
    #     )


# @shared_task
# def render_base(channel_id: int, tags: list[str]) -> bool:
#     logger.debug(f"render_base.channel_id={channel_id}, tags={tags}")
