from pydantic import BaseModel, EmailStr, Field
from typing_extensions import TypedDict


class To(TypedDict):
    user_id: str


class From(TypedDict):
    email: EmailStr


class Content(TypedDict):
    type: str
    value: str


class RequestAlarm(BaseModel):
    to: To
    from_: From = Field(..., alias="from")
    content: Content
    subject: str
