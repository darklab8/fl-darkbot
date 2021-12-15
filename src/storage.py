"module building all data to be stored"
import json
import os
import requests
from types import SimpleNamespace
import asyncio
from dotenv import load_dotenv

from src.forum_parser import forum_record, get_forum_threads
from src.info_controller import (AlertOnlyController, InfoController,
                                 InfoWithAlertController)
import src.settings as settings


class Storage():
    def __init__(self, unique_tag='DarkInfo:', IS_MOCKING_REQUESTS=False):
        self.unique_tag = unique_tag
        self.settings = self.load_env_settings()
        self.channels = self.load_channel_settings()

        self.IS_MOCKING_REQUESTS = IS_MOCKING_REQUESTS

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

    def get_players_data(self) -> dict:
        return requests.get(self.settings.player_request_url).json()

    def get_base_data(self) -> dict:
        return requests.get(self.settings.base_request_url).json()

    def get_only_not_repeated_records(self, forum_records,
                                      previous_forum_records):
        new_records = []
        for record in forum_records:
            if record.title not in previous_forum_records:
                new_records.append(record)
            else:
                if record.date != previous_forum_records[record.title].date:
                    new_records.append(record)
        return new_records

    def get_new_forum_records(self,
                              previous_forum_records={}) -> list[forum_record]:
        forum_records = get_forum_threads(
            forum_acc=self.settings.forum_acc,
            forum_pass=self.settings.forum_pass,
        )

        return self.get_only_not_repeated_records(forum_records,
                                                  previous_forum_records)

    def get_load_test_game_data(self, previous_forum_records):
        output = SimpleNamespace()

        with open('tests/examples/forum_records.json', 'r') as file_:
            data = file_.read()
            loaded = json.loads(data)
            morfed = [forum_record(**record) for record in loaded]

        output.new_forum_records = self.get_only_not_repeated_records(
            morfed, previous_forum_records)

        with open('tests/examples/players.json', 'r') as file_:
            output.players = json.loads(file_.read())

        with open('tests/examples/bases.json', 'r') as file_:
            output.bases = json.loads(file_.read())

        return output

    def get_game_data(self, previous_forum_records) -> SimpleNamespace:

        if self.IS_MOCKING_REQUESTS or settings.IS_MOCKING_REQUESTS:
            return self.get_load_test_game_data(previous_forum_records)

        output = SimpleNamespace()
        output.players = self.get_players_data()
        output.bases = self.get_base_data()
        output.new_forum_records = self.get_new_forum_records(
            previous_forum_records)
        return output

    async def a_get_game_data(self, previous_forum_records) -> SimpleNamespace:

        return await asyncio.to_thread(self.get_game_data,
                                       previous_forum_records)

    def base_add(self, name):
        print('adding the base')

    # def get_channel_data(self, key) -> SimpleNamespace:
    #     return deepcopy(self.storage.channels[key])
