from dependency_injector.wiring import inject
from api.v1.alarm.domain import AlarmDomain
from api.v1.alarm.repository import AlarmRepository
from api.v1.alarm.models import RequestAlarm
import logging


class AlarmServices:
    @inject
    def __init__(self, logger: logging, alarm_repository: AlarmRepository):
        self.logger = logger
        self.alarm_repository = alarm_repository

    @inject
    async def get_one(
        self,
        request_alarm: RequestAlarm,
    ) -> AlarmDomain | None:
        _alarm_domain: AlarmDomain = AlarmDomain(
            user_id=request_alarm.to["user_id"],
            email=request_alarm.from_["email"],
        )
        alarm_domain: AlarmDomain = await self.alarm_repository.one(_alarm_domain)

        return alarm_domain
