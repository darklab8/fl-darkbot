from listener.core import settings
from listener.core import app

if "main" in __name__:
    client = app.MyClient()
    client.run(settings.DISCORD_TOKEN)