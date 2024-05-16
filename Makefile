BUILD_ORG   := talkincode
BUILD_VERSION   := latest
BUILD_TIME      := $(shell date "+%F %T")
BUILD_NAME      := logsight
RELEASE_VERSION := v1.0.1
SOURCE          := main.go
RELEASE_DIR     := ./release
COMMIT_SHA1     := $(shell git show -s --format=%H )
COMMIT_DATE     := $(shell git show -s --format=%cD )
COMMIT_USER     := $(shell git show -s --format=%ce )
COMMIT_SUBJECT     := $(shell git show -s --format=%s )

buildpre:
	echo "BuildVersion=${BUILD_VERSION} ${RELEASE_VERSION} ${BUILD_TIME}" > assets/buildinfo.txt
	echo "ReleaseVersion=${RELEASE_VERSION}" >> assets/buildinfo.txt
	echo "BuildTime=${BUILD_TIME}" >> assets/buildinfo.txt
	echo "BuildName=${BUILD_NAME}" >> assets/buildinfo.txt
	echo "CommitID=${COMMIT_SHA1}" >> assets/buildinfo.txt
	echo "CommitDate=${COMMIT_DATE}" >> assets/buildinfo.txt
	echo "CommitUser=${COMMIT_USER}" >> assets/buildinfo.txt
	echo "CommitSubject=${COMMIT_SUBJECT}" >> assets/buildinfo.txt

fastpub:
	docker buildx build --platform=linux/amd64 --build-arg BTIME="$(shell date "+%F %T")" -t logsight .
	docker tag logsight ${BUILD_ORG}/logsight:latest
	docker push ${BUILD_ORG}/logsight:latest

fastpubm1:
	make build
	docker buildx build --platform=linux/amd64 --build-arg BTIME="$(shell date "+%F %T")" -t logsight . -f Dockerfile.local
	docker tag logsight ${BUILD_ORG}/logsight:latest-amd64
	docker push ${BUILD_ORG}/logsight:latest-amd64
	make buildarm64
	docker buildx build --platform=linux/arm64 --build-arg BTIME="$(shell date "+%F %T")" -t logsight . -f Dockerfile.local
	docker tag logsight ${BUILD_ORG}/logsight:latest-arm64
	docker push ${BUILD_ORG}/logsight:latest-arm64
	docker manifest create ${BUILD_ORG}/logsight:latest ${BUILD_ORG}/logsight:latest-arm64 ${BUILD_ORG}/logsight:latest-amd64
	# 标注不同架构镜像
	docker manifest annotate ${BUILD_ORG}/logsight:latest ${BUILD_ORG}/logsight:latest-amd64 --os linux --arch amd64
	docker manifest annotate ${BUILD_ORG}/logsight:latest ${BUILD_ORG}/logsight:latest-arm64 --os linux --arch arm64
	# 推送镜像
	docker manifest push ${BUILD_ORG}/logsight:latest

build:
	test -d ./release || mkdir -p ./release
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags  '-s -w -extldflags "-static"'  -o ./release/logsight main.go
	upx ./release/logsight

buildarm64:
	test -d ./release || mkdir -p ./release
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -ldflags  '-s -w -extldflags "-static"'  -o ./release/logsight main.go
	upx ./release/logsight

syncdev:
	make buildpre
	@read -p "提示:同步操作尽量在完成一个完整功能特性后进行，请输入提交描述 (develop):  " cimsg; \
	git commit -am "$(shell date "+%F %T") : $${cimsg}"
	# 切换主分支并更新
	git checkout main
	git pull origin main
	# 切换开发分支变基合并提交
	git checkout develop
	git rebase -i main
	# 切换回主分支并合并开发者分支，推送主分支到远程，方便其他开发者合并
	git checkout main
	git merge --no-ff develop
	git push origin main
	# 切换回自己的开发分支继续工作
	git checkout develop


updev:
	make buildpre
	make build
	scp ${RELEASE_DIR}/${BUILD_NAME} trdev-server:/tmp/logsight
	ssh trdev-server "systemctl stop logsight && /tmp/logsight -install && systemctl start logsight"

swag:
	swag fmt && swag init


.PHONY: clean build

