from sqlalchemy import Boolean, Column, ForeignKey, Integer, String
from sqlalchemy.orm import relationship
from ..database import Base


class Player(Base):
    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True)
    description = Column(String)
