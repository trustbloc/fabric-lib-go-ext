#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# Supported Targets:
# all : runs unit and integration tests
# depend: installs test dependencies
# unit-test: runs all the unit tests
# checks: runs all check conditions (license, spelling, linting)
# clean: stops docker conatainers used for integration testing
# populate: populates generated files (not included in git) - currently only vendor
# populate-vendor: populate the vendor directory based on the lock
# clean-populate: cleans up populated files (might become part of clean eventually)
# thirdparty-pin: pulls (and patches) pinned dependencies into the project under internal
#

# Tool commands (overridable)
GO_CMD             ?= go
DOCKER_CMD         ?= docker

# Fabric versions used in the Makefile
FABRIC_STABLE_VERSION           := 1.4.2
FABRIC_STABLE_VERSION_MINOR     := 1.4
FABRIC_STABLE_VERSION_MAJOR     := 1

FABRIC_PRERELEASE_VERSION       :=
FABRIC_PRERELEASE_VERSION_MINOR :=
FABRIC_PREV_VERSION             := 1.3.0
FABRIC_PREV_VERSION_MINOR       := 1.3
FABRIC_DEVSTABLE_VERSION_MINOR  := 1.4
FABRIC_DEVSTABLE_VERSION_MAJOR  := 1

# Build flags (overridable)
GO_LDFLAGS                 ?=
GO_TESTFLAGS               ?=
GO_TESTFLAGS_UNIT          ?= $(GO_TESTFLAGS)
GO_TESTFLAGS_INTEGRATION   ?= $(GO_TESTFLAGS) -failfast
FABRIC_LIB_GO_EXT_EXPERIMENTAL    ?= true
FABRIC_LIB_GO_EXT_EXTRA_GO_TAGS   ?=
FABRIC_LIB_GO_EXT_TEST_CHANGED  ?= false
FABRIC_LIB_GO_EXT_TESTRUN_ID    ?= $(shell date +'%Y%m%d%H%M%S')

# Dev tool versions (overridable)
GOLANGCI_LINT_VER ?= v1.19.1

# Fabric tool versions (overridable)
FABRIC_TOOLS_VERSION ?= $(FABRIC_STABLE_VERSION)

# Fabric tools docker image (overridable)
FABRIC_TOOLS_IMAGE ?= hyperledger/fabric-tools
FABRIC_TOOLS_TAG   ?= $(FABRIC_ARCH)-$(FABRIC_TOOLS_VERSION)

# Upstream fabric patching (overridable)
THIRDPARTY_FABRIC_BRANCH    ?= master
THIRDPARTY_FABRIC_COMMIT    ?= fedb583e7cb2998fef986a2a1a609f1f90beb983

# Force removal of images in cleanup (overridable)
FIXTURE_DOCKER_REMOVE_FORCE ?= false

# Options for exercising unit tests (overridable)
FABRIC_LIB_GO_EXT_DEPRECATED_UNITTEST ?= false

# Code levels
FABRIC_STABLE_CODELEVEL_TAG     := stable
FABRIC_PREV_CODELEVEL_TAG       := prev
FABRIC_PRERELEASE_CODELEVEL_TAG := prerelease
FABRIC_DEVSTABLE_CODELEVEL_TAG  := devstable
FABRIC_CODELEVEL_TAG            ?= $(FABRIC_STABLE_CODELEVEL_TAG)

# Code level version targets
FABRIC_STABLE_CODELEVEL_VER     := v$(FABRIC_STABLE_VERSION_MINOR)
FABRIC_PREV_CODELEVEL_VER       := v$(FABRIC_PREV_VERSION_MINOR)
FABRIC_PRERELEASE_CODELEVEL_VER := v$(FABRIC_PRERELEASE_VERSION_MINOR)
FABRIC_DEVSTABLE_CODELEVEL_VER  := v$(FABRIC_DEVSTABLE_VERSION_MINOR)
FABRIC_CODELEVEL_VER            ?= $(FABRIC_STABLE_CODELEVEL_VER)
FABRIC_CRYPTOCONFIG_VER         ?= v$(FABRIC_STABLE_VERSION_MAJOR)

# Code level to exercise during unit tests
FABRIC_CODELEVEL_UNITTEST_TAG ?= $(FABRIC_STABLE_CODELEVEL_TAG)
FABRIC_CODELEVEL_UNITTEST_VER ?= $(FABRIC_STABLE_CODELEVEL_VER)

# Local variables used by makefile
PROJECT_NAME           := fabric-lib-go-ext
ARCH                   := $(shell uname -m)
OS_NAME                := $(shell uname -s)
MAKEFILE_THIS          := $(lastword $(MAKEFILE_LIST))
THIS_PATH              := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
TEST_SCRIPTS_PATH      := test/scripts

ifneq ($(GO_LDFLAGS),)
GO_LDFLAGS_ARG := -ldflags=$(GO_LDFLAGS)
else
GO_LDFLAGS_ARG :=
endif

# Fabric tool docker tags at code levels
FABRIC_TOOLS_STABLE_TAG     = $(FABRIC_ARCH)-$(FABRIC_STABLE_VERSION)
FABRIC_TOOLS_PREV_TAG       = $(FABRIC_ARCH)-$(FABRIC_PREV_VERSION)
FABRIC_TOOLS_PRERELEASE_TAG = $(FABRIC_ARCH)-$(FABRIC_PRERELEASE_VERSION)
FABRIC_TOOLS_DEVSTABLE_TAG  := stable

# Detect CI
# TODO introduce nightly and adjust verify
ifdef JENKINS_URL
export FABRIC_LIB_GO_EXT_DEPEND_INSTALL=true
FABRIC_LIB_GO_EXT_TEST_CHANGED        := true
FABRIC_LIB_GO_EXT_DEPRECATED_UNITTEST   := false
endif

# Determine if internal dependency calc should be used
# If so, disable GOCACHE
ifeq ($(FABRIC_LIB_GO_EXT_TEST_CHANGED),true)
ifeq (,$(findstring $(GO_TESTFLAGS_UNIT),-count=1))
GO_TESTFLAGS_UNIT += -count=1
endif
ifeq (,$(findstring $(GO_TESTFLAGS_INTEGRATION),-count=1))
GO_TESTFLAGS_INTEGRATION += -count=1
endif
endif

# Setup Go Tags
GO_TAGS := $(FABRIC_LIB_GO_EXT_EXTRA_GO_TAGS)
ifeq ($(FABRIC_LIB_GO_EXT_EXPERIMENTAL),true)
GO_TAGS += experimental
endif

FABRIC_ARCH := $(ARCH)

ifneq ($(ARCH),x86_64)
# DEVSTABLE images are currently only x86_64
FABRIC_DEVSTABLE_INTTEST := false
else
# Recent Fabric builds follow GOARCH (e.g., amd64)
FABRIC_ARCH := amd64
endif

# Global environment exported for scripts
export GO_CMD
export ARCH
export FABRIC_ARCH
export GO_LDFLAGS
export GO_MOCKGEN_COMMIT
export GO_TAGS
export DOCKER_CMD
export FABRIC_LIB_GO_EXT_TESTRUN_ID
export GO111MODULE=on

.PHONY: all
all: version depend-noforce license unit-test

.PHONY: version
version:
	@$(TEST_SCRIPTS_PATH)/check_version.sh

.PHONY: depend
depend: version
	@$(TEST_SCRIPTS_PATH)/dependencies.sh -f

.PHONY: depend-noforce
depend-noforce: version
ifeq ($(FABRIC_LIB_GO_EXT_DEPEND_INSTALL),true)
	@$(TEST_SCRIPTS_PATH)/dependencies.sh
	@$(TEST_SCRIPTS_PATH)/dependencies.sh -c
else
	-@$(TEST_SCRIPTS_PATH)/dependencies.sh -c
endif

.PHONY: checks
checks: version depend-noforce license lint

.PHONY: license
license: version
	@$(TEST_SCRIPTS_PATH)/check_license.sh

.PHONY: lint
lint: version populate-noforce lint-submodules
	@MODULE="github.com/trustbloc/fabric-lib-go-ext" PKG_ROOT="./pkg" LINT_CHANGED_ONLY=true GOLANGCI_LINT_VER=$(GOLANGCI_LINT_VER) $(TEST_SCRIPTS_PATH)/check_lint.sh

.PHONY: lint-submodules
lint-submodules: version populate-noforce

.PHONY: lint-all
lint-all: version populate-noforce
	@MODULE="github.com/trustbloc/fabric-lib-go-ext" PKG_ROOT="./pkg" GOLANGCI_LINT_VER=$(GOLANGCI_LINT_VER) $(TEST_SCRIPTS_PATH)/check_lint.sh

.PHONY: unit-test
unit-test: clean depend-noforce populate-noforce license lint-submodules
	@TEST_CHANGED_ONLY=$(FABRIC_LIB_GO_EXT_TEST_CHANGED) TEST_WITH_LINTER=true FABRIC_LIB_GO_EXT_CODELEVEL_TAG=$(FABRIC_CODELEVEL_UNITTEST_TAG) FABRIC_LIB_GO_EXT_CODELEVEL_VER=$(FABRIC_CODELEVEL_UNITTEST_VER) \
	GO_TESTFLAGS="$(GO_TESTFLAGS_UNIT)" \
	GOLANGCI_LINT_VER="$(GOLANGCI_LINT_VER)" \
	MODULE="github.com/trustbloc/fabric-lib-go-ext" \
	PKG_ROOT="./pkg" \
	$(TEST_SCRIPTS_PATH)/unit.sh
ifeq ($(FABRIC_LIB_GO_EXT_DEPRECATED_UNITTEST),true)
	@GO_TAGS="$(GO_TAGS) deprecated" TEST_CHANGED_ONLY=$(FABRIC_LIB_GO_EXT_TEST_CHANGED) FABRIC_LIB_GO_EXT_CODELEVEL_TAG=$(FABRIC_CODELEVEL_UNITTEST_TAG) FABRIC_LIB_GO_EXT_CODELEVEL_VER=$(FABRIC_CODELEVEL_UNITTEST_VER) \
	GOLANGCI_LINT_VER="$(GOLANGCI_LINT_VER)" \
	MODULE="github.com/trustbloc/fabric-lib-go-ext" \
	PKG_ROOT="./pkg" \
	$(TEST_SCRIPTS_PATH)/unit.sh
endif

.PHONY: unit-tests
unit-tests: unit-test

.PHONY: crypto-gen
crypto-gen:
	@echo "Generating crypto directory ..."
	@$(DOCKER_CMD) run -i \
		-v /$(abspath .):/opt/workspace/$(PROJECT_NAME) -u $(shell id -u):$(shell id -g) \
		$(FABRIC_TOOLS_IMAGE):$(FABRIC_TOOLS_TAG) \
		//bin/bash -c "FABRIC_VERSION_DIR=fabric/$(FABRIC_CRYPTOCONFIG_VER) /opt/workspace/${PROJECT_NAME}/test/scripts/generate_crypto.sh"

.PHONY: thirdparty-pin
thirdparty-pin:
	@echo "Pinning third party packages ..."
	@UPSTREAM_COMMIT=$(THIRDPARTY_FABRIC_COMMIT) UPSTREAM_BRANCH=$(THIRDPARTY_FABRIC_BRANCH) scripts/third_party_pins/fabric/apply_upstream.sh

.PHONY: populate
populate: populate-vendor populate-fixtures-stable

.PHONY: populate-vendor
populate-vendor:
	@go mod vendor

.PHONY: populate-fixtures-stable
populate-fixtures-stable:
	@FABRIC_CRYPTOCONFIG_VERSION=$(FABRIC_CRYPTOCONFIG_VER) \
	FABRIC_FIXTURE_VERSION=v$(FABRIC_STABLE_VERSION_MINOR) \
	FABRIC_LIB_GO_EXT_CODELEVEL_TAG=$(FABRIC_STABLE_CODELEVEL_TAG) \
	$(TEST_SCRIPTS_PATH)/populate-fixtures.sh -f

.PHONY: populate-noforce
populate-noforce: populate-fixtures-stable-noforce

.PHONY: populate-fixtures-stable-noforce
populate-fixtures-stable-noforce:
	@FABRIC_CRYPTOCONFIG_VERSION=$(FABRIC_CRYPTOCONFIG_VER) \
	FABRIC_FIXTURE_VERSION=v$(FABRIC_STABLE_VERSION_MINOR) \
	FABRIC_LIB_GO_EXT_CODELEVEL_TAG=$(FABRIC_STABLE_CODELEVEL_TAG) \
	$(TEST_SCRIPTS_PATH)/populate-fixtures.sh

.PHONY: populate-fixtures-prev-noforce
populate-fixtures-prev-noforce:
	@FABRIC_CRYPTOCONFIG_VERSION=$(FABRIC_CRYPTOCONFIG_VER) \
	FABRIC_FIXTURE_VERSION=v$(FABRIC_PREV_VERSION_MINOR) \
	FABRIC_LIB_GO_EXT_CODELEVEL_TAG=$(FABRIC_PREV_CODELEVEL_TAG) \
	$(TEST_SCRIPTS_PATH)/populate-fixtures.sh

.PHONY: populate-fixtures-prerelease-noforce
populate-fixtures-prerelease-noforce:
	@FABRIC_CRYPTOCONFIG_VERSION=$(FABRIC_CRYPTOCONFIG_VER) \
	FABRIC_FIXTURE_VERSION=v$(FABRIC_PRERELEASE_VERSION_MINOR) \
	FABRIC_LIB_GO_EXT_CODELEVEL_TAG=$(FABRIC_PRERELEASE_CODELEVEL_TAG) \
	$(TEST_SCRIPTS_PATH)/populate-fixtures.sh

.PHONY: populate-fixtures-devstable-noforce
populate-fixtures-devstable-noforce:
	@FABRIC_CRYPTOCONFIG_VERSION=$(FABRIC_CRYPTOCONFIG_VER) \
	FABRIC_FIXTURE_VERSION=v$(FABRIC_DEVSTABLE_VERSION_MINOR) \
	FABRIC_LIB_GO_EXT_CODELEVEL_TAG=$(FABRIC_DEVSTABLE_CODELEVEL_TAG) \
	$(TEST_SCRIPTS_PATH)/populate-fixtures.sh

.PHONY: clean
clean: clean-fixtures clean-cache clean-populate

.PHONY: clean-populate
clean-populate:
	rm -Rf vendor

.PHONY: clean-cache
clean-cache:
ifeq ($(OS_NAME),Darwin)
	rm -Rf ${HOME}/Library/Caches/fabric-lib-go-ext
else
	rm -Rf ${HOME}/.cache/fabric-lib-go-ext
endif

.PHONY: clean-fixtures
clean-fixtures:
	-rm -Rf test/fixtures/fabric/*/crypto-config
	-rm -Rf test/fixtures/fabric/*/channel
