import unittest
from utils.config_parser.main import ConfigParser

class Testing(unittest.TestCase):

    def test_config(self):

        config = ConfigParser("SCRAPPY")
        config._env_storage = {
            "SCRAPPY_VAR_1": "abc"
        }

        self.assertEqual(config["var_1"], "abc")
        self.assertEqual(config.get("var_2"), None)
