from dependency_injector import containers, providers
from api.v1.alarm.repository import AlarmRepository
from api.v1.alarm.services import AlarmServices


class Container(containers.DeclarativeContainer):
    # logger
    logger = providers.Singleton()

    # PostgreSQL 리소스
    postgres_engine = providers.Resource()

    # Redis 리소스
    redis_client = providers.Resource()

    # Alarm repository
    alarm_repository = providers.Factory(
        AlarmRepository, logger=logger, rdb_session=postgres_engine.provided.get_pg_session
    )

    # Alarm services
    alarm_services = providers.Factory(AlarmServices, logger=logger, alarm_repository=alarm_repository)
