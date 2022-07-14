import os
from utils.config_parser import ConfigParser
 
config =  ConfigParser(settings_prefix="SCRAPPY")

DATABASE_URL = config.get("database_url","postgresql://postgres:postgres@localhost/")

DATABASE_NAME = "default"

API_PLAYER_URL = config.get("api_player_url")