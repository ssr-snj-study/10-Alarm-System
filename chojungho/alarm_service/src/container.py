from dependency_injector import containers, providers
from common import conf, setup_logging
from infrastructure.rdb.rdb_postgresql import AsyncEngine
from infrastructure.nosql.redis_client import init_redis_pool
from api.v1.alarm.container import Container as AlarmContainer
import pika


class Container(containers.DeclarativeContainer):
    wiring_config = containers.WiringConfiguration(
        packages=[
            "api.v1.alarm",
        ],
    )

    _config = conf()
    config = providers.Configuration()
    config.from_dict(_config.dict())

    # log 의존성 주입
    logger = providers.Singleton(setup_logging)

    # PostgreSQL 리소스
    postgres_engine = providers.Resource(AsyncEngine, config=config)

    # Redis 리소스
    redis_client = providers.Resource(
        init_redis_pool,
        host=config.REDIS_HOST,
        port=config.REDIS_PORT,
        password=config.REDIS_PASSWORD,
    )

    # RabbitMQ 리소스
    rabbimq_connection = providers.Factory(
        pika.BlockingConnection,
        pika.ConnectionParameters(host="localhost", port=5672, credentials=pika.PlainCredentials("admin", "admin")),
    )

    # api
    alarm_container = providers.Container(
        AlarmContainer,
        logger=logger,
        postgres_engine=postgres_engine,
        redis_client=redis_client,
        rabbimq_connection=rabbimq_connection,
    )
