from dependency_injector import containers, providers
from common import FilterRule


class Container(containers.DeclarativeContainer):
    # logger
    logger = providers.Singleton()

    # filter url
    filter_rule: FilterRule = providers.Singleton(FilterRule)

    # PostgreSQL 리소스
    postgres_engine = providers.Resource()

    # Redis 리소스
    redis_client = providers.Resource()
