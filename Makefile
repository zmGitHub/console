all: deps linux_build

export GOPATH	  := $(CURDIR)/_project
export GOBIN 	  := $(CURDIR)/bin
DOCKER_IMAGE_TAG  ?= $(shell git describe --abbrev=0 --tags)
COMMONENVVAR      ?= GOOS=linux GOARCH=amd64
BUILDENVVAR       ?= CGO_ENABLED=0

CURRENT_GIT_GROUP := bitbucket.org/forfd
CURRENT_GIT_REPO := custm-chat

folder_dep:
	mkdir -p $(CURDIR)/_project/src/$(CURRENT_GIT_GROUP)
	test -d $(CURDIR)/_project/src/$(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO) || ln -s $(CURDIR) $(CURDIR)/_project/src/$(CURRENT_GIT_GROUP)

vendor: # make add_dep dep="package name"
	govendor init

deps: folder_dep
	mkdir -p $(CURDIR)/vendor
	glide install

build: deps
	go build -o bin/$(CURRENT_GIT_REPO) -ldflags "-extldflags -static -X main.BuildTime=`date '+%Y-%m-%d_%I:%M:%S%p'` -X main.BuildGitHash=`git rev-parse HEAD` -X main.BuildGitTag=`git describe --tags 2>/dev/null || echo NO_TAG`" $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)

linux_build:
	$(COMMONENVVAR) $(BUILDENVVAR) make build

test: deps
	mysql -h 127.0.0.1 -uroot -e "source $(CURDIR)/db/webim/sql/all.sql"
	go test -v $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/webim/handler

clean:
	@rm -rf vendor bin _project

.PHONY: deps install test add_dep clean 

