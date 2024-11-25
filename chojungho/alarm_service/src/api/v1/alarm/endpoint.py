from fastapi import APIRouter
from api.v1.alarm.models import RequestAlarm

# from dependency_injector.wiring import inject
from fastapi.responses import JSONResponse

router = APIRouter(prefix="/alarm", tags=["Alarm"])


@router.post("")
# @inject
async def post_alarm(request_alarm: RequestAlarm) -> JSONResponse:
    return ...
