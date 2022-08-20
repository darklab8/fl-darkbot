from .subtasks import SubTaskGetBaseData
import pytest


@pytest.mark.integration
def test_request_base_url() -> None:
    data = SubTaskGetBaseData()
    from pprint import pprint as print

    print(data)
