# type: ignore
import os
from utils.config_parser.main import ConfigParser
import pytest


def test_config():
    config = ConfigParser("SCRAPPY")
    os.environ["SCRAPPY_VAR_1"] = "abc"
    assert config["var_1"] == "abc"


def test_get_exception():
    config = ConfigParser("SCRAPPY")

    with pytest.raises(KeyError):
        config.get("var_2")


def test_get_default():
    config = ConfigParser("SCRAPPY")
    assert config.get("var_2", "my_default") == "my_default"
