VERSION=`date "+%Y%m%d"`

.DEFAULT: build

build:
	@GOOS=linux go build -ldflags="-s -w -X main.Version=${VERSION}" -o haproxy-mysql-gr-healthcheck main.go
	@echo "haproxy-mysql-gr-healthcheck has been built."
	@du -sh haproxy-mysql-gr-healthcheck
