import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.openapi.utils import get_openapi
from api.v1.alarm.endpoint import router
from container import Container

container = Container()


def create_app(_config) -> FastAPI:
    _app = FastAPI(title=_config.PROJECT_NAME)

    _app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_methods=("GET", "POST", "PUT", "DELETE"),
        allow_headers=["*"],
    )

    _app.include_router(router, prefix="/api/v1")

    def game_credit_openapi():
        if _app.openapi_schema:
            return _app.openapi_schema
        openapi_schema = get_openapi(
            title=_config.PROJECT_NAME,
            version=_config.VERSION,
            routes=_app.routes,
        )
        _app.openapi_schema = openapi_schema
        return _app.openapi_schema

    _app.openapi = game_credit_openapi

    return _app


app = create_app(container.config.from_dict())

if __name__ == "__main__":
    uvicorn.run("main:app", port=5959, reload=True)
