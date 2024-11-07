# Builder stage: 애플리케이션을 빌드하고, 최종 이미지를 위한 실행 파일을 준비하는 단계입니다.
FROM golang:1.22-alpine as builder

# 필요한 패키지 설치: git, 인증서, upx, tzdata 등
RUN apk update && apk add --no-cache git ca-certificates upx tzdata

# 애플리케이션 파일을 저장할 작업 디렉토리를 설정합니다.
WORKDIR /usr/src/app

# Go 모듈 시스템을 활성화하고, 프록시 설정을 지정합니다.
# - GO111MODULE=on: Go 모듈 사용을 강제합니다.
# - GOPROXY=https://proxy.golang.org,direct: 모듈 다운로드 시 프록시 서버를 통해 캐싱된 모듈을 빠르게 가져오며, 프록시에 없는 경우 직접 다운로드합니다.
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct

# 의존성 파일을 컨테이너로 복사합니다.
COPY go.mod go.sum ./

# go.mod와 go.sum 파일을 사용하여 필요한 모든 의존성을 다운로드합니다.
RUN go mod download

# 의존성 검증: go.sum 파일에 기록된 체크섬을 기반으로 다운로드된 모듈의 무결성을 검증합니다.
RUN go mod verify

# 애플리케이션의 모든 소스 코드를 컨테이너로 복사합니다.
COPY . .

# 불필요한 모듈을 정리합니다. 사용하지 않는 패키지를 제거하여 빌드 효율성을 높입니다.
RUN go mod tidy

# 애플리케이션 빌드: OS와 아키텍처를 설정한 후, 최적화된 실행 파일을 생성합니다.
# - CGO_ENABLED=0: C 라이브러리와의 의존성을 제거하여 독립적인 실행 파일을 만듭니다.
# - GOOS=linux GOARCH=arm64: 리눅스 및 ARM64 환경에 맞춰 빌드합니다.
# - -ldflags="-s -w": 디버그 정보를 제거하여 바이너리 파일 크기를 줄입니다.
# - -o bin/main: 빌드된 실행 파일을 지정된 위치에 저장합니다.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -v -ldflags="-s -w" -o bin/main main.go | tee build.log

# Optional: 실행 파일을 압축하여 크기를 줄입니다.
# - upx --best --lzma: UPX를 사용해 최적의 압축을 적용하고, LZMA 알고리즘을 사용하여 추가 압축합니다.
RUN upx --best --lzma bin/main

# Executable image stage: 빌드된 실행 파일을 최종 이미지에 복사하고 설정하는 단계입니다.
FROM scratch

# 인증서와 사용자 정보 복사: scratch 이미지는 빈 이미지이므로 인증서와 사용자 정보를 추가합니다.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
# 타임존 정보 복사
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# 빌드 단계에서 생성된 실행 파일을 최종 이미지에 복사합니다.
COPY --from=builder /usr/src/app/bin/main ./main

# 환경 변수 설정
# - SERVER_MODE=test: 애플리케이션을 테스트 모드로 실행하도록 설정합니다.
# - TZ=Asia/Seoul: Asia/Seoul 타임존을 설정하여 애플리케이션이 해당 타임존을 사용하도록 합니다.
ENV SERVER_MODE=test
ENV TZ=Asia/Seoul

# 비-루트 유저로 실행되도록 설정하여 보안을 강화합니다. (ID가 1000인 사용자)
USER 1000

# 애플리케이션이 수신할 포트를 설정합니다.
EXPOSE 8080

# 컨테이너 시작 시 애플리케이션의 실행 파일을 바로 실행하도록 설정합니다.
ENTRYPOINT ["./main"]