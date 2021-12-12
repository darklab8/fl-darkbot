import os

from dotenv import load_dotenv
from io import StringIO


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