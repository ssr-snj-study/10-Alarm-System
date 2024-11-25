# coding=utf-8
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

    async def one(self, alarm_domain: AlarmDomain) -> AlarmDomain | None:
        async with self.rdb_session() as session:
            user_alarm_info = await session.execute(
                select(
                    AlarmUser.user_id,
                    AlarmUser.phone_number,
                    AlarmUser.email,
                    AlarmDevice.device_token,  # 정확한 모델의 필드 이름 사용
                )
                .join(AlarmDevice, AlarmDevice.user_id == AlarmUser.user_id)
                .where(AlarmUser.user_id == alarm_domain.user_id)  # 명시적으로 join 조건 추가
            )
            results = user_alarm_info.all()
            return results
