from pydantic_settings import BaseSettings
from os import environ


class Config(BaseSettings):
    PROJECT_NAME: str = "Alarm Service"
    VERSION: str = "0.1.0"

    # Postgresql
    POSTGRES_USER: str = ""
    POSTGRES_PASSWORD: str = ""
    POSTGRES_HOST: str = ""
    POSTGRES_PORT: int
    POSTGRES_DB: str = ""

    # Redis
    REDIS_HOST: str = ""
    REDIS_PORT: int
    REDIS_PASSWORD: str = ""

    # RabbitMQ
    RABBITMQ_HOST: str = ""
    RABBITMQ_PORT: int
    RABBITMQ_USER: str = ""
    RABBITMQ_PASSWORD: str = ""

    class Config:
        env_file = ".env"


class LocalConfig(Config):
    DEBUG: bool = True
    SQL_PRINT: bool = True

    POSTGRES_SERVER: str = ""
    POSTGRES_SCHEMA: str = "public"

    REDIS_SERVER: str = ""


def conf():
    c = dict(local=LocalConfig)
    return c[environ.get("API_ENV", "local")]()
