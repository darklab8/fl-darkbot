from enum import Enum, auto


class RequestMethod(Enum):
    get = auto()
    post = auto()
    delete = auto()
    put = auto()
    patch = auto()
