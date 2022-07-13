from sqlalchemy import Boolean, Column, ForeignKey, Integer, String
from sqlalchemy.orm import relationship
import scrappy.core.databases as databases


class Player(databases.default.Base):
    __tablename__ = "players"

    id = Column(Integer, primary_key=True, index=True)
    description = Column(String)
