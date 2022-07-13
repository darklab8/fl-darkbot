import argparse
from . import databases
import os

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
args = parser.parse_args()

match args.action:
    case "run":
        os.system("uvicorn scrappy.core.main:app")
    case "makemigrations":
        os.system(f'alembic -c scrappy/alembic.ini revision --autogenerate -m "{args.migration_name}"')
    case "drop":
        databases.default.Base.metadata.drop_all(bind=databases.default.engine)
    case "create":
        databases.default.Base.metadata.create_all(bind=databases.default.engine)
