import os
import shutil

import click

from .universal import say


@click.group()
def sphinx():
    "sphinx commands"
    pass


@sphinx.command()
def build():
    "build sphinx documentation"

    if os.path.exists('docs'):
        shutil.rmtree('docs')

    if not os.path.exists('docs'):
        os.mkdir('docs')

    if not os.path.exists(os.path.join('docs', '.nojekyll')):
        os.mknod(os.path.join("docs", ".nojekyll"))

    say(f"sphinx-build -b html {os.path.join('sphinx','source')} docs")
