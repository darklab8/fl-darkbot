import click

from .universal import say


@click.group()
def manage():
    "manage commands"
    pass


@manage.command()
def run():
    "launch server"
    say("python app.py")
