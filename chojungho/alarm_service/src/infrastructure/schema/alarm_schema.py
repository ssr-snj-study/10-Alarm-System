from sqlalchemy.orm import mapped_column, relationship
from sqlalchemy import Integer, BigInteger, String, TIMESTAMP, ForeignKey

from .base import Base


class AlarmUser(Base):
    __tablename__ = "alarm_user"

    user_id = mapped_column(BigInteger, primary_key=True)
    email = mapped_column(String, nullable=False)
    country_code = mapped_column(Integer, nullable=False)
    phone_number = mapped_column(Integer, nullable=False)
    created_at = mapped_column(TIMESTAMP, nullable=False, server_default="CURRENT_TIMESTAMP")

    # Relationship
    devices = relationship("AlarmDevice", back_populates="user", cascade="all, delete")

    def __repr__(self):
        return f"<AlarmUser(user_id={self.user_id}, email='{self.email}')>"


class AlarmDevice(Base):
    __tablename__ = "alarm_device"

    id = mapped_column(BigInteger, primary_key=True)
    device_token = mapped_column(String, nullable=False)
    user_id = mapped_column(BigInteger, ForeignKey("alarm_user.user_id", ondelete="CASCADE"))
    last_logged_in_at = mapped_column(TIMESTAMP)

    # Relationship
    user = relationship("AlarmUser", back_populates="devices")

    def __repr__(self):
        return f"<AlarmDevice(id={self.id}, device_token='{self.device_token}')>"
