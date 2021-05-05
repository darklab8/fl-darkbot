"starting module"
import os
from types import SimpleNamespace
from dotenv import load_dotenv
from commands import created_app
import json

# nice settings loading
load_dotenv()

STORAGE = SimpleNamespace()
STORAGE.unique_tag = 'dark_info:'

STORAGE.settings = SimpleNamespace()
for item, value in os.environ.items():
    setattr(STORAGE.settings, item, value)

# loading persistent settings
try:
    with open('channels.json') as file_:
        STORAGE.channels = json.loads(file_.read())
except FileNotFoundError:
    STORAGE.channels = {}

bot = created_app(STORAGE)

if __name__ == '__main__':
    bot.run(STORAGE.settings.secret_key)
