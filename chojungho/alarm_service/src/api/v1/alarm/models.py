from pydantic import BaseModel, EmailStr
from typing_extensions import TypedDict


class To(TypedDict):
    user_id: str


class From(TypedDict):
    email: EmailStr


class Content(TypedDict):
    type: str
    value: str


class RequestAlarm(BaseModel):
    to_target: To
    from_target: From
    content: Content
    subject: str
