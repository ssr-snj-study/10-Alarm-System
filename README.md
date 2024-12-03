# [가상 면접 사례로 배우는 대규모 시스템 설계 기초](https://www.yes24.com/Product/Goods/102819435)
저자: 알렉스쉬

---

# 8장 URL 단축기 설계

## 기본 기능

- URL 단축
- URL 리디렉션
- 높은 가용성 및 규모 확장성 그리고 장애 감내

## 개략적 추정

- 쓰기 연산: 매일 1억 개의 단축 URL 생성
- 초당 쓰기 연산 1억/24/3600 = 1160
- 읽기 연산: 초당 11,600회
- 10년간 보관 가능한 레코드 개수는 1억*365*10 = 3650억
- 10년간 필요한 저장 용량은 URL 평균 길이가 100일 때 365 Billion*100Byte = 36.5TB

## API 엔드포인트

- URL 단축 엔드포인트
    - Parameter ⇒ 단축할 URL
    - Method ⇒ POST
    - URL ⇒ /api/v1/data/shorten
    - Response ⇒ 단축 URL
- URL 리디렉션 엔드포인트
    - Parameter ⇒ 단축 URL
    - Method ⇒ GET
    - URL ⇒ /api/v1/shortUrl
    - Response ⇒ 리디렉션될 원래 URL

## URL 리디렉션
- 301 Permanently Moved
    - Location 헤더에 반환된 URL로 이전
    - 캐시된 응답
- 302 Found
    - Location 헤더가 지정하는 URL에 의해 처리
    - 원래 서버로 리디렉션

## URL 단축 플로
- 긴 URL과 해시값은 1대1 매핑
- 변환 및 복원이 가능해야 함

## 데이터 모델

- 관계형 DB에 아래와 같이 저장
    - 테이블명: `url`
    - 컬럼: `id`(PK), `shortUrl`, `longUrl`

## 해시 함수

- 해시값은 총 62개(`[0-9a-zA-Z]`)로 62^7=3.5조개의 URL을 만들 수 있으므로 7자리의 수로 만듦
- 해시 함수 구현은 2가지 방법이 있음
- 해시 후 충돌 해소 방법은 CRC32, MD5, SHA-1 등의 해시 함수로 축약하여 사용
- base-62 변환 방법은 10진수 ID를 62진수로 변환하여 62개의 해시 문자와 매핑

  | **해시 후 충돌 전략** | **base-62 변환** | 
  |----------------|----------------|
  | 단축 URL의 길이 고정  | 단축 URL 길이는 가변적 |
  | ID 생성기 필요없음    | ID 생성기 필요(유일성) |
  | 충돌 해소 필요       | 보안상 문제         |

## 상세 설계

- 긴 URL 입력 → DB 확인 → 있으면 반환, 없으면 단축 URL 생성 후 반환
- URL 단축기는 쓰기보다 읽기를 자주하는 시스템이므로 1:1 매핑(`단축-원래`) 캐싱 사용

## 마무리

- 처리율 제한 장치, 웹 서버의 규모 확장, DB 규모, 데이터 분석 솔루션, 가용성, 데이터 일관성, 안정성에 대한 내용을 심층적으로 분석 가능

## 참고문헌

[1] [**REST API Tutorial**](https://www.restapitutorial.com/)  
[2] [**Bloom filter**](https://en.wikipedia.org/wiki/Bloom_filter)

---

# Result

## URL Shortcut Architecture

![URL Shortcut Architecture](./images/architecture.png)

## Web Server

![Web Page](./images/web_page.png)

### Load Balance

![Load Balance](./images/load_balance.png)

### URL Redirection

- FastAPI 기본 모듈인 [RedirectResponse](https://fastapi.tiangolo.com/it/advanced/custom-response/#redirectresponse) 사용
- Response 받았을 때 응답 코드는 307 

![Response](./images/api_response.png)

## Database and Cache

![Database](./images/api_response.png)

![Cache](./images/postgres_table.png)