# 왁스코드 이미지 프록시 서버
## 구축하기
```sh
docker pull opentypefont/wakscord-proxy:latest
docker run -d -p <port>:<port> -e HOST=<host> -e PORT=<port> opentypefont/wakscord-proxy:latest
```

## FAQ
1. 왜 클라이언트는 net/http를 사용하나요?
    > fasthttp의 클라이언트는 리스폰스 바디를 무조건 읽도록 설계되어있어요.
    > 따라서 메모리 점유율이 지나치게 상승할 수 있기에 기본적으로 바디 스트림을 가져오는 net/http를 사용합니다.