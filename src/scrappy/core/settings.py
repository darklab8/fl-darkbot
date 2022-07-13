import os


class ConfigParser:

    def __init__(self, settings_prefix=""):
        self.settings_prefix = settings_prefix

    def __getitem__(self, variable_name):
        return os.environ[
            self._transform_to_env_var_name(f"{self.settings_prefix}_{variable_name}")
        ]

    def get(self, variable_name, default=None):
        return os.environ.get(
            self._transform_to_env_var_name(f"{self.settings_prefix}_{variable_name}"), default
        )

    @staticmethod
    def _transform_to_env_var_name(variable_name: str):
        return variable_name.replace(".", "_")
    


config =  ConfigParser(settings_prefix="SCRAPPY")


API_PLAYER_URL = os.environ["API_PLAYER_URL"]