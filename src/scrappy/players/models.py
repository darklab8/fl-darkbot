from sqlalchemy import Column, Integer, String, DateTime
from sqlalchemy.orm import relationship
import scrappy.core.databases as databases


class Player(databases.default.Base):
    __tablename__ = "players"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String)
    region = Column(String)
    system = Column(String)
    time = Column(String)  # time online
    timestamp = Column(DateTime)
