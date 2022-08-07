from .subtasks import SubTaskGetBaseData


def test_request_base_url():
    data = SubTaskGetBaseData()
    from pprint import pprint as print

    print(data)
