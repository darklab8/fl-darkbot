from utils.logger import Logger


from utils.config_parser import ConfigParser

config = ConfigParser(settings_prefix="configurator")
LOGGER_CONSOLE_LEVEL = config.get("logger.console.level", "INFO")

base_logger = Logger(console_level=LOGGER_CONSOLE_LEVEL)

logger = base_logger.getChild(__name__)
