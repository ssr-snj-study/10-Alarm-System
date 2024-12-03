import os
import redis
from dotenv import load_dotenv

load_dotenv()


class RedisQueue:
    def __init__(self, queue_name, host=os.getenv("REDIS_SERVER"), password=os.getenv("REDIS_PASSWORD"), port=os.getenv("REDIS_PORT"), db=0):
        """
        Redis 큐 초기화
        :param queue_name: 큐 이름
        :param host: Redis 서버 호스트
        :param port: Redis 서버 포트
        :param db: Redis 데이터베이스 번호
        """
        self.queue_name = queue_name
        self.redis = redis.StrictRedis(host=host, port=port, db=db, password=password, decode_responses=True)

    def push(self, value):
        """
        큐에 값을 추가 (오른쪽으로 삽입)
        :param value: 추가할 값
        """
        self.redis.rpush(self.queue_name, value)
        print(f"값 추가: {value}")

    def pop(self):
        """
        큐에서 값을 가져오기 (왼쪽에서 제거)
        :return: 큐에서 꺼낸 값
        """
        value = self.redis.lpop(self.queue_name)
        print(f"값 가져오기: {value}")
        return value

    def size(self):
        """
        큐의 현재 크기 반환
        :return: 큐의 크기
        """
        size = self.redis.llen(self.queue_name)
        print(f"큐 크기: {size}")
        return size

    def peek(self):
        """
        큐에서 값을 제거하지 않고 확인
        :return: 큐의 첫 번째 값
        """
        value = self.redis.lindex(self.queue_name, 0)
        print(f"큐 첫 번째 값: {value}")
        return value
