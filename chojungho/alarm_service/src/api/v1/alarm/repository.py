from contextlib import AbstractAsyncContextManager
from typing import Callable

from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from api.v1.alarm.domain import AlarmDomain
from infrastructure.schema.alarm_schema import AlarmUser, AlarmDevice
from dependency_injector.wiring import inject
import logging


class AlarmRepository:
    @inject
    def __init__(self, logger: logging, rdb_session: Callable[..., AbstractAsyncContextManager[AsyncSession]]) -> None:
        self.logger = logger
        self.rdb_session = rdb_session

    async def one(self, alarm_domain: AlarmDomain) -> list[AlarmDomain] | None:
        async with self.rdb_session() as session:
            user_alarm_info = await session.execute(
                select(
                    AlarmUser.user_id,
                    AlarmUser.phone_number,
                    AlarmUser.email,
                    AlarmUser.country_code,
                    AlarmUser.created_at,
                    AlarmDevice.device_token,
                    AlarmDevice.last_logged_in_at,
                )
                .join(AlarmDevice, AlarmDevice.user_id == AlarmUser.user_id)
                .where(AlarmUser.user_id == alarm_domain.user_id)
            )
            results = user_alarm_info.all()
            if not results:
                return None
            alarm_domains: list[AlarmDomain] = [
                alarm_domain.from_dict(user_alarm_info._asdict()) for user_alarm_info in results
            ]
            return alarm_domains
