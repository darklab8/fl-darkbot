import random

def shuffled_dict(data):

    shuffled = list(data.values())
    random.shuffle(shuffled)
    randomed = dict(zip(data, shuffled))

    return randomed
