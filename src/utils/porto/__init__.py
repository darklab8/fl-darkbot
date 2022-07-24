import abc


class AbstractAction(abc.ABC):
    def __new__(cls, *args, **kwargs):
        instance = super(AbstractAction, cls).__new__(cls)
        instance.__init__(*args, **kwargs)
        result = instance.run()
        return result

    @abc.abstractmethod
    def run(self):
        pass
