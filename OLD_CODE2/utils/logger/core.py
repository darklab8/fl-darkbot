import logging
from enum import Enum, EnumMeta
from typing import Union


class _EnumDirectValueMeta(EnumMeta):
    def __getattribute__(cls, name: str) -> int:
        value = super().__getattribute__(name)
        if isinstance(value, cls):
            value = value.value
        return value  # type: ignore

    def __getitem__(cls, name: str) -> int:  # type: ignore
        value = super().__getattribute__(name)
        if isinstance(value, cls):
            value = value.value
        return value  # type: ignore


class LoggerLevels(Enum, metaclass=_EnumDirectValueMeta):
    DEBUG = logging.DEBUG  # purely for diagnostics
    INFO = logging.INFO  # confirming normal flow of program
    WARN = logging.WARNING  # a thing nice to fix like deprecation warning
    ERROR = logging.ERROR  # minor handled error of app, program continues to run
    FATAL = logging.CRITICAL  # a error that disrupted program workflow completely


class CustomFilter(logging.Filter):

    COLOR = {
        "DEBUG": "GREEN",
        "INFO": "GREEN",
        "WARNING": "YELLOW",
        "ERROR": "RED",
        "CRITICAL": "RED",
    }

    def filter(self, record: logging.LogRecord) -> bool:
        record.color = CustomFilter.COLOR[record.levelname]
        return True


import logging.handlers


class Logger:

    levels = LoggerLevels
    _parent: Union["Logger", None] = None

    def __init__(
        self,
        console_level: str,
        _parent: Union["Logger", None] = None,
        name: str = "",
    ):
        self._name = name
        self._console_level = console_level
        self._parent = _parent
        self._logger = self._configure_logger(
            logging.getLogger("").getChild(self._name)
            if _parent is None
            else _parent._logger.getChild(self._name)  # type: ignore
        )

    def _configure_logger(self, logger: logging.Logger) -> logging.Logger:
        # global level, controlling available levels in handlers
        logger.setLevel(logging.DEBUG)

        # create console handler with a higher log level
        ch = logging.StreamHandler()
        ch.addFilter(CustomFilter())
        ch.setLevel(LoggerLevels[self._console_level])  # type: ignore
        formatter = logging.Formatter(
            " - ".join(
                [
                    "time:%(asctime)s",
                    "color:%(color)s",
                    "level:%(levelname)s",
                    "name:%(name)s",
                    "msg:%(message)s",
                ]
            )
        )
        ch.setFormatter(formatter)
        # add the handlers to the logger
        logger.addHandler(ch)
        return logger

    def getChild(self, name: str) -> "Logger":
        return self.__class__(
            console_level=self._console_level,
            name=name,
            _parent=self,
        )

    def debug(self, msg: object) -> None:
        self._logger.debug(msg)

    def info(self, msg: object) -> None:
        self._logger.info(msg)

    def warn(self, msg: object) -> None:
        self._logger.warning(msg)

    def error(self, msg: object) -> None:
        self._logger.error(msg)

    def fatal(self, msg: object) -> None:
        self._logger.critical(msg)
