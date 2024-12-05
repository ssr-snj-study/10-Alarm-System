from dependency_injector.wiring import inject
from api.v1.alarm.domain import AlarmDomain
from api.v1.alarm.repository import AlarmRepository
from api.v1.alarm.models import RequestAlarm
import logging
import pika
import json


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

    async def insert_queue(self, message: str, alarm_domain: AlarmDomain) -> None:
        # 큐에 넣는 작업
        alarm_domain_to_dict: dict = alarm_domain.to_dict()
        alarm_domain_to_dict.update(
            {
                "message": message,
                "created_at": str(alarm_domain_to_dict["created_at"]),
                "last_logged_in_at": str(alarm_domain_to_dict["last_logged_in_at"]),
                "phone_number": str(alarm_domain_to_dict["phone_number"]),
            }
        )
        queue_body = json.dumps(alarm_domain_to_dict).encode()
        channel = self.rabbitmq.channel()
        device_type = {1: "Android", 2: "ios", 3: "sms", 4: "email"}
        channel.queue_declare(queue=device_type[alarm_domain.device_type])
        channel.basic_publish(exchange="", routing_key=device_type[alarm_domain.device_type], body=queue_body)
