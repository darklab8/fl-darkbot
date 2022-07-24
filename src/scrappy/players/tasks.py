from celery import shared_task


@shared_task
def test(arg):
    print(arg)


@shared_task
def add(x, y):
    z = x + y
    print(z)
    return z