"starting module"
from commands import created_app
from storage import storage_builder

# nice settings loading

if __name__ == '__main__':
    storage = storage_builder()
    bot = created_app(storage)
    bot.run(storage.settings.secret_key)
