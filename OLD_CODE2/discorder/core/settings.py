from utils.config_parser import ConfigParser
from utils.logger import Logger

logger = Logger(console_level="DEBUG", name=__name__)
config = ConfigParser(settings_prefix="discorder")

DISCORD_TOKEN = config.get("discord.token")
LOGGER_CONSOLE_LEVEL = config.get("logger.console.level", "INFO")
