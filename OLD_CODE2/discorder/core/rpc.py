from utils.porto.rpc import AbstractRPCAction
from utils.config_parser import ConfigParser
from utils.logger import Logger

logger = Logger(console_level="DEBUG", name=__name__)
config = ConfigParser(settings_prefix="discorder")

DISCORDER_API_URL = config.get("api_url", "http://localhost:8000")


class RPCAction(AbstractRPCAction):
    api_url = DISCORDER_API_URL
