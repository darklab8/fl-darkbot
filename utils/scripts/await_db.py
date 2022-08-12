import time
import psycopg2
from dataclasses import dataclass
from contextlib import contextmanager
import argparse
import sys

parser = argparse.ArgumentParser()
parser.add_argument("--db_name", type=str, default="default")
parser.add_argument("--user", type=str, default="postgres")
parser.add_argument("--password", type=str, default="postgres")
parser.add_argument("--host", type=str, default="db")
parser.add_argument("--port", type=str, default="5432")
parser.add_argument("--timeout", type=int, default=30)
args = parser.parse_args()


@dataclass
class DatabaseParams:
    db_name: str = args.db_name
    user: str = args.user
    password: str = args.password
    host: str = args.host
    port: str = args.port


@contextmanager
def open_database(params: DatabaseParams):
    database = psycopg2.connect(
        " ".join(
            [
                f"dbname={params.db_name}",
                f"user={params.user}",
                f"password={params.password}",
                f"host={params.host}",
                f"port={params.port}",
            ]
        )
    )
    try:
        yield database
    finally:
        database.close()


@contextmanager
def open_cursor(database):
    with database.cursor() as cursor:
        yield cursor
        database.commit()


if __name__ == "__main__":
    loop_delay = 3
    for i in range(int(args.timeout / loop_delay)):
        try:
            with open_database(DatabaseParams()) as database:
                with open_cursor(database) as cursor:
                    cursor.execute("select 1;")
                    print("await_db: db is ready to accept connections")
                    sys.exit(0)
        except psycopg2.OperationalError:
            print(
                "await_db: db is not available yet. "
                "Sleeping 3 seconds more. Timeout 30 seconds."
            )
            time.sleep(loop_delay)

    raise Exception("non zero exit")
