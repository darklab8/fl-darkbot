import os


class ConfigParser:

    _env_storage: dict = os.environ

    def __init__(self, settings_prefix: str = ""):
        self.settings_prefix = settings_prefix.upper()

    def __getitem__(self, variable_name):
        return self._env_storage[
            self._transform_to_env_var_name(f"{self.settings_prefix}_{variable_name}")
        ]

    def get(self, variable_name, default=None):
        return self._env_storage.get(
            self._transform_to_env_var_name(f"{self.settings_prefix}_{variable_name}"),
            default,
        )

    @staticmethod
    def _transform_to_env_var_name(variable_name: str):
        return variable_name.replace(".", "_").upper()
