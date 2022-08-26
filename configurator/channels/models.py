from enum import unique
from sqlalchemy import Column, Integer, String, DateTime, ForeignKey, UniqueConstraint
from sqlalchemy.orm import relationship
from ..core.base import Model
import datetime


class Channel(Model):
    __tablename__ = "channel"

    id = Column(Integer, primary_key=True, index=True)
    owners = relationship("ChannelOwner")


class ChannelOwnerConstraints:
    owner_channel = "_owner_channel_uc"


class ChannelOwner(Model):
    __tablename__ = "owner"
    id = Column(Integer, primary_key=True, index=True)
    owner_id = Column(Integer, index=True)
    owner_name = Column(String)
    channel_id = Column(Integer, ForeignKey("channel.id"))
    created = Column(DateTime, default=datetime.datetime.utcnow)

    __table_args__ = (
        UniqueConstraint(
            "owner_id", "channel_id", name=ChannelOwnerConstraints.owner_channel
        ),
    )
