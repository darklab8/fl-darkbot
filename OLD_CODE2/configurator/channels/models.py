from enum import unique
from sqlalchemy import (
    Column,
    Integer,
    String,
    DateTime,
    ForeignKey,
    UniqueConstraint,
    BigInteger,
)
from sqlalchemy.orm import relationship
from ..core.base import Model
import datetime


class Channel(Model):
    __tablename__ = "channel"

    id = Column(BigInteger, primary_key=True, index=True)
    owners = relationship("ChannelOwner")


class ChannelOwnerConstraints:
    owner_channel = "_owner_channel_uc"


class ChannelOwner(Model):
    __tablename__ = "owner"
    id = Column(BigInteger, primary_key=True, index=True)
    name = Column(String)
    channel_id = Column(BigInteger, ForeignKey("channel.id"))
    created = Column(DateTime, default=datetime.datetime.utcnow)

    __table_args__ = (
        UniqueConstraint(
            "id", "channel_id", name=ChannelOwnerConstraints.owner_channel
        ),
    )
