import os

import click

PROJECT_NAME = os.path.basename(os.getcwd())
print(f"project_name = {PROJECT_NAME}")


def say(phrase) -> None:
    click.echo(phrase)
    os.system(phrase)
