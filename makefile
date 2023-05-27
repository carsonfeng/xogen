GIT_BRANCH=$(shell git symbolic-ref --short -q HEAD)
GIT_COMMIT=$(shell git rev-list HEAD | head -n 1)
BUILD_TIME=$(shell date +"%Y-%m-%d_%H:%M:%S%Z")

BUILD_LDFLAGS = "-X go_pk/build.Branch=${GIT_BRANCH} \
-X go_pk/build.Commit=${GIT_COMMIT} \
-X go_pk/build.BuildTime=${BUILD_TIME} \
"

ifeq ($(P), LINUX)
	GOBUILD := GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go install -v -ldflags ${BUILD_LDFLAGS}
else
	GOBUILD := go install -v -ldflags ${BUILD_LDFLAGS}
endif

# 如果 makefile 文件执行有问题, 请检查是否定义了 PYCOMM_DIR 和 PYTHONPATH 两个变量
# export BRANCH := $(shell git branch | grep '*' | tr -d '* ')
export DEVELOP := develop
export CURRENT_DIR = $(shell pwd)
# 在当前目录下 bin 文件中生成 可执行文件
# export GOBIN = $(CURRENT_DIR)/bin
# 更新 GOPATH 变量, 这样可以共享 GOPATH/src 下的文件包, 避免重复下载依赖包
export GOPATH := $(CURRENT_DIR):$(GOPATH)

export TAGNAME := $(shell date "+%Y%m%d-%H%M")
export GO111MODULE=auto

tag:
	git tag $(t) -m $(m)
	git push origin $(t)

tag_del:
	git tag -d $(t)
	git push origin :$(t)