from sqlalchemy import Boolean, Column, ForeignKey, Integer, String
from sqlalchemy.orm import relationship
import src.scrappy.database as database


class Player(database.Base):
    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True)
    description = Column(String)
