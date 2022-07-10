import src.scrappy.databases as databases
import players.models as models

models.databases.Base.metadata.create_all(bind=databases.engine)
