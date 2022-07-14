import os
from utils.config_parser import ConfigParser
 
config =  ConfigParser(settings_prefix="SCRAPPY")

DATABASE_URL = config.get("database_url","postgresql://postgres:postgres@localhost/")

def get_database_url(name):
    return DATABASE_URL + name

DEFAULT_DATABASE = get_database_url("default")

API_PLAYER_URL = config.get("api_player_url")