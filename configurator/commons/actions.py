import abc
from typing import Any


class AbstractAction(abc.ABC):
    @abc.abstractmethod
    async def run(self, *args, **kwargs) -> Any:
        pass
