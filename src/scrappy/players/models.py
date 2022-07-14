from sqlalchemy import Column, Integer, String
from sqlalchemy.orm import relationship
import scrappy.core.databases as databases


class Player(databases.default.Base):
    __tablename__ = "players"

    id = Column(Integer, primary_key=True, index=True)
    description = Column(String)
