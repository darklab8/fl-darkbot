from utils.config_parser import ConfigParser
from utils.logger import Logger

logger = Logger(console_level="DEBUG", name=__name__)
config = ConfigParser(settings_prefix="configurator")

DATABASE_USER = config["database_username"]
DATABASE_PASSWORD = config["database_password"]
DATABASE_HOST = config["database_host"]
DATABASE_URL = config.get(
    "database_url", f"{DATABASE_USER}:{DATABASE_PASSWORD}@{DATABASE_HOST}/"
)
DATABASE_NAME = "default"


LOGGER_CONSOLE_LEVEL = config.get("logger.console.level", "INFO")
