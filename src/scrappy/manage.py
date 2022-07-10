import database as database
import players.models as models

models.database.Base.metadata.create_all(bind=database.engine)
