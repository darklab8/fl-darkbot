import os


class ConfigParser:

    _env_storage: os._Environ[str] = os.environ

    def __init__(self, settings_prefix: str = ""):
        self.settings_prefix = settings_prefix.upper()

    def __getitem__(self, variable_name: str) -> str:
        return self._env_storage[
            self._transform_to_env_var_name(f"{self.settings_prefix}.{variable_name}")
        ]

    def get(self, variable_name: str, default: str | None = None) -> str:
        result = self._env_storage.get(
            self._transform_to_env_var_name(f"{self.settings_prefix}.{variable_name}"),
            default,
        )

        if result is None:
            raise KeyError("No value is discovered, and no default value is supplied")

        return result

    @staticmethod
    def _transform_to_env_var_name(variable_name: str) -> str:
        return variable_name.replace(".", "_").upper()
