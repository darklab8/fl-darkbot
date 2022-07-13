import argparse
import scrappy.core.databases as databases
import os

def makemigrations(migration_name):
    os.system(f'alembic -c scrappy/alembic.ini revision --autogenerate -m "{migration_name}"')

def migrate(migration_name):
    os.system(f'alembic -c scrappy/alembic.ini revision --autogenerate -m "{migration_name}"')

parser = argparse.ArgumentParser(
    description="Copying selected by regex strings to new file"
)
parser.add_argument(
    "--action",
    type=str,
)
parser.add_argument(
    "--migration_name",
    type=str,
    default=""
)

if "manage" in __name__:
    args = parser.parse_args()

    match args.action:
        case "run":
            os.system("uvicorn scrappy.core.main:app")
        case "makemigrations":
            makemigrations(args.migration_name)
        case "drop":
            databases.default.Base.metadata.drop_all(bind=databases.default.engine)
        case "create":
            databases.default.Base.metadata.create_all(bind=databases.default.engine)
