from fastapi import APIRouter, Depends
from api.v1.alarm.models import RequestAlarm

from dependency_injector.wiring import inject
from fastapi.responses import JSONResponse
from dependency_injector.wiring import Provide
from api.v1.alarm.services import AlarmServices

router = APIRouter(prefix="/alarm", tags=["Alarm"])


@router.post("")
@inject
async def post_alarm(
    request_alarm: RequestAlarm, alarm_service: AlarmServices = Depends(Provide["alarm_container.alarm_services"])
) -> JSONResponse:
    await alarm_service.get_one(request_alarm)
    return ...
