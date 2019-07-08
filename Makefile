app = "whateats"

all: build run

build:
	@go build -ldflags="-X 'github.com/SergeyParamoshkin/version.BuildTime=$(date "+%F %T")' -X 'github.com/SergeyParamoshkin/version.AppName=${app}' -X 'github.com/SergeyParamoshkin/version.Version=v0.0.0' -X 'github.com/SergeyParamoshkin/version.Commit=$(git rev-parse --short HEAD)'"

run:
	@./${app} -proxy="" -proxy_user="" -proxy_password="" -token=""