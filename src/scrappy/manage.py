from .database import engine
from .players.models import Base

Base.metadata.create_all(bind=engine)
