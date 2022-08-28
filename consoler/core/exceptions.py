class ConsoleException(Exception):
    pass


class NotRegisteredCommand(ConsoleException):
    pass


class NotImplementedMethod(ConsoleException):
    pass
