import asyncio
import random

import firebase_admin
from firebase_admin import credentials, firestore, messaging

from lib.redis import RedisQueue


def iter_user_info_from_firestore():
    """Firestore에서 사용자 정보를 가져오는 함수"""
    try:
        db = firestore.client()
        # fcm_tokens 컬렉션에서 해당 토큰을 검색
        collection = db.collection(u'fcm_tokens')
        docs = collection.stream()
        # doc = db.collection('fcm_tokens').where('token', '==', token).get()

        for doc in docs:
            if doc.to_dict().get('device'):
                yield doc.to_dict()['user'], doc.to_dict()['device']

    except Exception as e:
        print(f"Firestore 데이터 가져오기 실패: {e}")


def send_notification_to_device(token, title, body, image=None):
    """특정 기기에 알림 보내기"""
    try:
        # 메시지 생성
        message = messaging.Message(
            notification=messaging.Notification(
                title=title,
                body=body,
                image=image
            ),
            # 앱에서 사용할 수 있는 값
            # data={
            #     'title': title,
            #     'message': body,
            #     'mode': 'test',
            #     'data': body
            # },
            token=token,
        )

        # 메시지 전송
        response = messaging.send(message)
        print(f"알림 전송 완료! 응답: {response}")
    except Exception as e:
        print(f"알림 전송 실패: {e}")


async def push_notification(cnt):
    queue = RedisQueue('notification')
    queue.push(f'지금부터 {cnt}번의 알림을 무작위로 보냅니다.')
    total_time = 0
    for i in range(cnt):
        sleep_for = random.uniform(1, 5)
        total_time += sleep_for
        await asyncio.sleep(sleep_for)
        queue.push(f'{i + 1}번째 알림입니다. 총 {round(total_time, 2)}초 경과하였습니다.')


async def check_notification():
    queue = RedisQueue('notification')
    while True:
        if queue.size():
            msg = queue.pop()
            # Firestore에서 사용자 정보 가져오기
            for user, device in iter_user_info_from_firestore():
                send_notification_to_device(
                    token=device,
                    title='test',
                    body=f'{user}님, {msg}',
                    image='https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcShSW8GGrxy-WZfbdkQjXJsG6H5_WDr2HG2XQ&s',
                )
        else:
            await asyncio.sleep(1)


async def main():
    asyncio.create_task(push_notification(10))

    # 2개의 테스크 생성
    await asyncio.wait({
        asyncio.create_task(check_notification()),
        asyncio.create_task(check_notification()),
    })


if __name__ == '__main__':
    # Firebase Admin 초기화
    cred = credentials.Certificate('serviceAccountKey.json')
    default_app = firebase_admin.initialize_app(cred)

    asyncio.run(main())
