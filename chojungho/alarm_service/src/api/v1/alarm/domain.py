from dataclasses import dataclass, asdict, field


@dataclass
class AlarmDomain:
    user_id: int = field(default=None)
    email: str = field(default=None)
    country_code: str = field(default=None)
    phone_number: str = field(default=None)
    created_at: str = field(default=None)
    device_token: str = field(default=None)
    device_type: int = field(default=None)
    last_logged_in_at: str = field(default=None)

    @classmethod
    def from_dict(cls, d):
        return cls(**d)

    def to_dict(self):
        return asdict(self)

    def delete_to_dict_none_data(self):
        return {k: v for k, v in self.to_dict().items() if v is not None}
