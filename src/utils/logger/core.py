import logging


class LoggerLevels:
    # purely for diagnostics
    DEBUG = logging.DEBUG

    # confirming normal flow of program
    INFO = logging.INFO

    # a thing nice to fix like deprecation warning
    WARN = logging.WARNING

    # handled error
    ERROR = logging.ERROR

    # fatal unhandled error
    FATAL = logging.CRITICAL


class CustomFilter(logging.Filter):

    COLOR = {
        "DEBUG": "GREEN",
        "INFO": "GREEN",
        "WARNING": "YELLOW",
        "ERROR": "RED",
        "CRITICAL": "RED",
    }

    def filter(self, record):
        record.color = CustomFilter.COLOR[record.levelname]
        return True


import logging.handlers


class Logger:

    levels = LoggerLevels

    def __init__(self, console_level: str, name: str | None = None):
        self._name = self._generate_name(name)
        self._console_level = console_level
        self._logger = self._get_logger()

    def _transform_level_to_logging_lib(self, level: str) -> int:
        return getattr(LoggerLevels, level)

    def _get_logger(self) -> logging.Logger:
        # create logger with 'spam_application'
        logger = logging.getLogger("")
        # making sure we don't see third party dependencies debug logging
        logger = logger.getChild(self._name)
        # global level, controlling available levels in handlers
        logger.setLevel(logging.DEBUG)

        # create console handler with a higher log level
        ch = logging.StreamHandler()
        ch.addFilter(CustomFilter())
        ch.setLevel(self._transform_level_to_logging_lib(self._console_level))
        formatter = logging.Formatter(
            " - ".join(
                [
                    "time:%(asctime)s",
                    "color:%(color)s",
                    "level:%(levelname)s",
                    "path:%(name)s.%(filename)s",
                    "msg:%(message)s",
                ]
            )
        )
        ch.setFormatter(formatter)
        # add the handlers to the logger
        logger.addHandler(ch)
        return logger

    def _generate_name(self, name: str is None) -> str:
        return __name__ if name is None else name

    def getChild(self, name: str):
        return self.__class__(
            console_level=self._console_level,
            name=name,
        )

    def debug(self, msg):
        self._logger.debug(msg)

    def info(self, msg):
        self._logger.info(msg)

    def warn(self, msg):
        self._logger.warning(msg)

    def error(self, msg):
        self._logger.error(msg)

    def fatal(self, msg):
        self._logger.critical(msg)
