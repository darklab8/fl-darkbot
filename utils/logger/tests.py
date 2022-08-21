from .core import LoggerLevels
import logging


def test_logger_level_transformation() -> None:
    assert LoggerLevels["DEBUG"] == LoggerLevels.DEBUG
    assert LoggerLevels.DEBUG == logging.DEBUG
