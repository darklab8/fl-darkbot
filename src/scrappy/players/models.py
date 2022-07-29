from sqlalchemy import Column, Integer, String, DateTime
from sqlalchemy.orm import relationship
from scrappy.core.base import Model


class Player(Model):
    __tablename__ = "players"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String)
    region = Column(String)
    system = Column(String)
    time = Column(String)  # time online
    timestamp = Column(DateTime)


class Thing(Model):
    __tablename__ = "things"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String)
    region = Column(String)
    system = Column(String)
    time = Column(String)  # time online
    timestamp = Column(DateTime)
