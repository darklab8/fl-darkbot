from utils.logger import Logger
from . import settings as settings

base_logger = Logger(
    console_level=settings.LOGGER_CONSOLE_LEVEL,
    name="",
)

logger = base_logger.getChild(__name__)
