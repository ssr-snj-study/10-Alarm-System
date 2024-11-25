from fastapi import APIRouter, Request
from dependency_injector.wiring import inject
from fastapi.responses import JSONResponse

router = APIRouter(prefix="/alarm", tags=["Alarm"])


@router.post("")
@inject
async def post_alarm(request: Request) -> JSONResponse:

    return ...
