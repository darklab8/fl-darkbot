from utils.config_parser import ConfigParser
from utils.logger import Logger

logger = Logger(console_level="DEBUG", name=__name__)
config = ConfigParser(settings_prefix="viewer")

CELERY_BROKER = config["celery.broker"]
CELERY_BACKEND = config["celery.backend"]

LOGGER_CONSOLE_LEVEL = config.get("logger.console.level", "INFO")

CONFIGURATOR_API = config.get(
    "CONFIGURATOR_API_URL", default="http://configurator_web:8000"
)

SCRAPPY_API = config.get("CONFIGURATOR_API_URL", default="http://scrappy_web:8000")
