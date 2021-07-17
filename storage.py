"module building all data to be stored"
import json
import os
import requests
from types import SimpleNamespace

from dotenv import load_dotenv

from info_controller import (AlertOnlyController, InfoController,
                             InfoWithAlertController)


class Storage():
    def __init__(self, unique_tag='DarkInfo:'):
        self.unique_tag = unique_tag
        self.settings = self.load_env_settings()
        self.channels = self.load_channel_settings()

        # forum thread tracker
        self.forum = InfoController(self.channels, 'forum')
        self.base = InfoController(self.channels, 'base')
        self.system = InfoController(self.channels, 'system')
        self.region = InfoController(self.channels, 'region')
        self.friend = InfoWithAlertController(self.channels, 'friend')
        self.enemy = InfoWithAlertController(self.channels, 'enemy')
        self.unrecognized = AlertOnlyController(self.channels, 'unrecognized')

    def load_env_settings(self) -> SimpleNamespace:
        "loading settings from os environment"
        load_dotenv()

        output = SimpleNamespace()
        for item, value in os.environ.items():
            setattr(output, item, value)
        return output

    def load_channel_settings(self) -> dict:
        """loadding perssistent settings
        set by users about channels"""
        output = {}
        try:
            with open('data/channels.json', 'r') as file_:
                output = json.loads(file_.read())
        except FileNotFoundError:
            print('ERR failed to load channels.json')
        return output

    def save_channel_settings(self) -> None:
        """loadding perssistent settings
        set by users about channels"""
        try:
            with open('data/channels.json', 'w') as file_:
                file_.write(json.dumps(self.channels, indent=2))
        except OSError as error:
            print('ERR failed to save channels.json ' + str(error))

    def get_game_data(self) -> SimpleNamespace:
        output = SimpleNamespace()
        output.players = requests.get(self.settings.player_request_url).json()
        output.bases = requests.get(self.settings.base_request_url).json()
        return output

    def base_add(self, name):
        print('adding the base')

    # def get_channel_data(self, key) -> SimpleNamespace:
    #     return deepcopy(self.storage.channels[key])
