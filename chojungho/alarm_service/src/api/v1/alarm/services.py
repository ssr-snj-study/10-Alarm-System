from dependency_injector.wiring import inject
from api.v1.alarm.domain import AlarmDomain
from api.v1.alarm.repository import AlarmRepository
from api.v1.alarm.models import RequestAlarm
import logging
import pika


class AlarmServices:
    @inject
    def __init__(self, logger: logging, alarm_repository: AlarmRepository, rabbitmq: pika.BlockingConnection):
        self.logger = logger
        self.alarm_repository = alarm_repository
        self.rabbitmq = rabbitmq

    async def get_one(
        self,
        request_alarm: RequestAlarm,
    ) -> list[AlarmDomain] | None:
        _alarm_domain: AlarmDomain = AlarmDomain(
            user_id=request_alarm.to["user_id"],
            email=request_alarm.from_["email"],
        )
        alarm_domains: list[AlarmDomain] | None = await self.alarm_repository.one(_alarm_domain)
        if alarm_domains is None:
            return None

        return alarm_domains

    async def insert_queue(self, message: str, alarm_domain: AlarmDomain) -> str | None:
        # 큐에 넣는 작업
        channel = self.rabbitmq.channel()
        channel.queue_declare(queue="android")
        channel.basic_publish(exchange="", routing_key="android", body=message)
