VERSION=`date "+%Y%m%d"`

.DEFAULT: build

build:
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.Version=${VERSION}" -o haproxy-mysql-gr-healthcheck.amd64 main.go
	@echo "haproxy-mysql-gr-healthcheck.amd64 has been built."
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.Version=${VERSION}" -o haproxy-mysql-gr-healthcheck.arm64 main.go
	@echo "haproxy-mysql-gr-healthcheck.arm64 has been built."
	@du -sh haproxy-mysql-gr-healthcheck.*
