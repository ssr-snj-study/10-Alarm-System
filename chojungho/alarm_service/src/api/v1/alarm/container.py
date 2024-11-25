from dependency_injector import containers, providers


class Container(containers.DeclarativeContainer):
    # logger
    logger = providers.Singleton()

    # PostgreSQL 리소스
    postgres_engine = providers.Resource()

    # Redis 리소스
    redis_client = providers.Resource()
