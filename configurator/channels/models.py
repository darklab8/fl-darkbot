from sqlalchemy import Column, Integer, String, DateTime, ForeignKey
from sqlalchemy.orm import relationship
from ..core.base import Model


class Channel(Model):
    __tablename__ = "channel"

    id = Column(Integer, primary_key=True, index=True)
    owners = relationship("ChannelOwner")


class ChannelOwner(Model):
    __tablename__ = ""
    id = Column(Integer, primary_key=True, index=True)
    channel = Column(Integer, ForeignKey("channel.id"))
    created = Column(DateTime)
