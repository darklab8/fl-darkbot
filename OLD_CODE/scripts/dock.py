import click

from .universal import PROJECT_NAME, say


@click.group()
def dock():
    "docker commands"
    pass


@dock.command()
def logs():
    say(f"docker logs -t {PROJECT_NAME}:latest")


def builder() -> None:
    say(f"git pull && docker build -t {PROJECT_NAME}:latest .")


@dock.command()
def build():
    builder()


def runner():
    say(f"docker run --name {PROJECT_NAME} -t "
        f"-d --rm {PROJECT_NAME}:latest")


@dock.command()
def run():
    runner()


def stopper() -> None:
    say('docker stop $(docker ps -a -q --filter="' f"name={PROJECT_NAME}" '")')


@dock.command()
def stop():
    stopper()


def cleaner() -> None:
    "getting rid of already built docker layers"
    say(f"docker rmi $(docker images '{PROJECT_NAME}' -a -q)")


@dock.command()
def clean():
    cleaner()


@dock.command()
def deploy():
    "command to deploy/or redeploy from zero"
    cleaner()
    stopper()
    builder()
    runner()
