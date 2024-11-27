from fastapi import APIRouter, Depends
from api.v1.alarm.models import RequestAlarm

from api.v1.alarm.domain import AlarmDomain
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
    # 넘어온 파라미터로 사용자 조회
    alarm_domains: list[AlarmDomain] = await alarm_service.get_one(request_alarm)
    if not alarm_domains:
        return JSONResponse(content={"status": 404, "msg": "Not Found", "code": 404, "list": []})

    # 사용자에게 알림전송하기위해 큐에 넣는작업
    for alarm_domain in alarm_domains:
        await alarm_service.insert_queue(request_alarm.content["value"], alarm_domain)

    return JSONResponse(content={"status": 200, "msg": "알람전송 완료", "code": 200, "list": []})
