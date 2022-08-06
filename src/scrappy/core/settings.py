from utils.config_parser import ConfigParser

config = ConfigParser(settings_prefix="SCRAPPY")

DATABASE_USER = config.get("database_username")
DATABASE_PASSWORD = config.get("database_password")
DATABASE_HOST = config.get("database_host")
DATABASE_URL = config.get(
    "database_url", f"{DATABASE_USER}:{DATABASE_PASSWORD}@{DATABASE_HOST}/"
)

DATABASE_NAME = "default"
CELERY_BROKER = config.get("celery.broker", "redis://redis:6379/0")
CELERY_BACKEND = config.get("celery.backend", "redis://redis:6379/0")


API_PLAYER_URL = config.get("player.url")
API_BASE_URL = config.get("base.url")
FORUM_USERNAME = config.get("forum.username")
FORUMN_PASSWORD = config.get("forum.password")
