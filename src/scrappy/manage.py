import argparse
from . import databases

parser = argparse.ArgumentParser(
    description="Copying selected by regex strings to new file"
)
parser.add_argument(
    "--action",
    type=str,
)
args = parser.parse_args()

match args.action:
    case "drop_tables":
        databases.default.Base.metadata.drop_all(bind=databases.default.engine)
    case "create_tables":
        databases.default.Base.metadata.create_all(bind=databases.default.engine)
