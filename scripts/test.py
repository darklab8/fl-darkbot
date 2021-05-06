import click
from .universal import say


@click.group()
def test():
    "testing commands"
    pass


@test.command()
def pylint():
    "link with pylint"
    say("pylint `ls -d *.py | grep -v \"venv\" | grep"
        " -v \"scripts\" | grep -v \"__pycache__\" | grep -v \"sphinx\"`")


@test.command()
def flake():
    "lint with flake8"
    say("".join(("flake8 --exclude .git,venv,*/migrations/*,.tox .")))


@test.command()
@click.option('--refresh',
              '-r',
              is_flag=True,
              help="enables refresh of data examples",
              default=False)
@click.option('--cover',
              '-c',
              'cover',
              is_flag=True,
              help="shows coverage",
              default=False)
@click.option('--mypy',
              '-m',
              'mypy',
              is_flag=True,
              help="shows hint coverage",
              default=False)
def unit(refresh, cover, mypy):
    "get unit tests"
    launcher = []
    launcher.append("pytest -n 6")
    if cover:
        launcher.append("-cov-config=.coveragerc --cov=.")

    if mypy:
        launcher.append(" --mypy tests.py")

    say(" ".join(launcher))


@test.command()
def tox():
    "full test run to be done between commits"
    say("tox -r")


@test.command()
def mypy():
    "type hinting checker"
    say("mypy .")
