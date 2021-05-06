import click
from scripts.dock import dock
from scripts.test import test
from scripts.sphinx import sphinx
from scripts.manage import manage


@click.group()
@click.pass_context
def root(context):
    "root commands"
    pass


root.add_command(dock)
root.add_command(test)
root.add_command(sphinx)
root.add_command(manage)

if __name__ == '__main__':
    root(obj={})
