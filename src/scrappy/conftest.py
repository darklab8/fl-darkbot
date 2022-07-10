import databases as databases
import pytest
from sqlalchemy import create_engine
from sqlalchemy_utils import database_exists, create_database


@pytest.fixture
def db():
    database_url = "postgresql://postgres:postgres@localhost/test_database"
    engine = create_engine(database_url)
    if not database_exists(engine.url):
        create_database(engine.url)

    database = databases.Database(
        # url="sqlite:///./test_sql_app.db"
        url=database_url
    )

    databases.default.Base.metadata.drop_all(bind=database.engine)
    databases.default.Base.metadata.create_all(bind=database.engine)

    with database.manager_to_get_db() as db:
        yield db
