PROJECT        :=carousel
GOOS           :=$(shell go env GOOS)
GOARCH         :=$(shell go env GOARCH)
GOMODULECMD    :=$(shell go list -m)/cmd
RELEASE_ROOT   ?=release
CAROUSEL_PATH  =./$(PROJECT)
TARGETS        ?=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64


SEMVER_VERSION    ?=0.8.1
SEMVER_PRERELEASE ?=
SEMVER_BUILDMETA  ?=
BUILD_DATE        :=$(shell date -u -Iseconds)
BUILD_VCS_URL     :=$(shell git config --get remote.origin.url) 
BUILD_VCS_ID      :=$(shell git log -n 1 --date=iso-strict-local --format="%h")
BUILD_VCS_ID_DATE :=$(shell TZ=UTC0 git log -n 1 --date=iso-strict-local --format='%ad')


build: SEMVER_PRERELEASE := dev

GO_LDFLAGS = -ldflags="-X '$(GOMODULECMD).SemVerVersion=$(SEMVER_VERSION)' \
	            -X '$(GOMODULECMD).SemVerPrerelease=$(SEMVER_PRERELEASE)' \
	            -X '$(GOMODULECMD).SemVerBuildMeta=$(SEMVER_BUILDMETA)' \
	            -X '$(GOMODULECMD).BuildDate=$(BUILD_DATE)' \
	            -X '$(GOMODULECMD).BuildVcsUrl=$(BUILD_VCS_URL)' \
	            -X '$(GOMODULECMD).BuildVcsId=$(BUILD_VCS_ID)' \
		    -X '$(GOMODULECMD).BuildVcsIdDate=$(BUILD_VCS_ID_DATE)'"

.PHONY: build test require-% release-% clean

build:
	CGO_ENABLED=0 go build $(GO_LDFLAGS) -o $(CAROUSEL_PATH)
	$(CAROUSEL_PATH) version

test: $(if $(wildcard $(CAROUSEL_PATH)),build)
	$(TEST_PATH) $(CAROUSEL_PATH) ${VAULT_VERSIONS}

require-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

RELEASES := $(foreach target,$(TARGETS),release-$(target)-$(PROJECT))

release-all: $(RELEASES)

define build-target
release-$(1)/$(2)-$(PROJECT): # require-VERSION
	@echo "Building $(PROJECT) $(VERSION) ($(1)/$(2)) ..." 
	CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build -o $(RELEASE_ROOT)/$(PROJECT)-$(1)-$(2)$(if $(patsubst windows,,$(1)),,.exe) $(GO_LDFLAGS)
	@ls -la $(RELEASE_ROOT)/$(PROJECT)-$(1)-$(2)$(if $(patsubst windows,,$(1)),,.exe)
	@echo ""
endef

$(foreach target,$(TARGETS),$(eval $(call build-target,$(word 1, $(subst /, ,$(target))),$(word 2, $(subst /, ,$(target))))))

clean:
	rm -rf $(CAROUSEL_PATH) $(RELEASE_ROOT)/* 

.DEFAULT_GOAL := release-all

# test:
#	ginkgo watch ./...
	
# test-ci:
#	ginkgo  ./...

# gen:
#	go generate ./...
# docker:
#	docker build -t $(docker_registry) .

# publish: docker
#	docker push $(docker_registry)

# fmt:
#	find . -name '*.go' | while read -r f; do \
#		gofmt -w -s "$$f"; \
#	done

# .DEFAULT_GOAL := docker

# .PHONY: go-mod docker-build docker-push docker test fmt
