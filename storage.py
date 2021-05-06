"module building all data to be stored"
import os
from types import SimpleNamespace
import json
from dotenv import load_dotenv


def _load_env_settings() -> SimpleNamespace:
    "loading settings from os environment"
    load_dotenv()

    output = SimpleNamespace()
    for item, value in os.environ.items():
        setattr(output, item, value)
    return output


def _load_channel_settings() -> dict:
    """loadding perssistent settings
    set by users about channels"""
    output = {}
    try:
        with open('channels.json') as file_:
            output = json.loads(file_.read())
    except FileNotFoundError:
        print('ERR failed to load channels.json')
    return output


def storage_builder(unique_tag='dark_info:'):
    "building all settings for the application"
    return SimpleNamespace(unique_tag=unique_tag,
                           settings=_load_env_settings(),
                           channels=_load_channel_settings())
