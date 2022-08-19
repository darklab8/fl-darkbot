import abc
from typing import Any


class AbstractAction(abc.ABC):
    def __new__(cls, *args: Any, **kwargs: dict[str, Any]) -> Any:
        instance = super(AbstractAction, cls).__new__(cls)
        instance.__init__(*args, **kwargs)  # type: ignore
        result = instance.run()
        return result

    @abc.abstractmethod
    def run(self) -> Any:
        pass
