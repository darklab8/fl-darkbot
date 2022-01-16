import os

from dotenv import load_dotenv
from io import StringIO

import logging
logging.basicConfig(format='%(levelname)s: %(asctime)s %(message)s', datefmt='%m/%d/%Y %I:%M:%S %p', level=logging.INFO)
logging.info("init logging")


def testing_enable_debug_in_env_and_get_status():
    config = StringIO("debug=true")
    load_dotenv(config)
    return get_debug_status()


def get_debug_status():
    """
    >>> get_debug_status();
    False

    >>> testing_enable_debug_in_env_and_get_status()
    True
    
    """
    return bool(os.environ.get("debug", False))


load_dotenv()
DEBUG = bool(os.environ.get("debug", False))

IS_MOCKING_REQUESTS = bool(os.environ.get("IS_MOCKING_REQUESTS", False))