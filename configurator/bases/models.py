from sqlalchemy import Column, Integer, String, DateTime, ForeignKey, UniqueConstraint
from sqlalchemy.orm import relationship
from ..core.base import Model
import datetime


class BaseConstraints:
    channel_tag = "_tag_channel_id_uc"


class Base(Model):
    __tablename__ = "base"
    id = Column(Integer, primary_key=True, index=True)
    tag = Column(String)
    channel_id = Column(Integer, ForeignKey("channel.id"))
    created = Column(DateTime, default=datetime.datetime.utcnow)

    __table_args__ = (
        UniqueConstraint("tag", "channel_id", name=BaseConstraints.channel_tag),
    )
