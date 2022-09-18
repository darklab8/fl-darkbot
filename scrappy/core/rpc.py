from utils.porto.rpc import AbstractRPCAction
from utils.config_parser import ConfigParser
from utils.logger import Logger

logger = Logger(console_level="DEBUG", name=__name__)
config = ConfigParser(settings_prefix="scrappy")

SCRAPPY_API_URL = config.get("api_url", "http://scrappy_web:8000")


class RPCAction(AbstractRPCAction):
    api_url = SCRAPPY_API_URL
