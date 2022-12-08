from sqlalchemy import Column, Integer, String, DateTime
from scrappy.core.base import Model


class Player(Model):
    __tablename__ = "players"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, unique=True, index=True)
    region = Column(String)
    system = Column(String)
    time = Column(String)  # time online
    timestamp = Column(DateTime)
