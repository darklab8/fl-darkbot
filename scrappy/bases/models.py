from sqlalchemy import Column, Integer, String, DateTime, Float
from scrappy.core.base import Model


class Base(Model):
    __tablename__ = "bases"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String)
    affiliation = Column(String)
    health = Column(Float)
    tid = Column(Integer)
    timestamp = Column(DateTime)
