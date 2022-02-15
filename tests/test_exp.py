import random


def test_shifflying_dict():
    data = {"a": 1, "b": 2, "c": 3, "d": 4}

    shuffled = list(data.values())
    random.shuffle(shuffled)
    randomed = dict(zip(data, shuffled))

    print(randomed)