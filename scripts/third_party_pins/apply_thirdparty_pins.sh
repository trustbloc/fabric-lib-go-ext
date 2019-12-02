#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# This script fetches code originating from other upstream projects
# These files are checked into internal paths.

set -e

if output=$(git status --porcelain) && [ -z "$output" ]
then
  echo "Working directory clean, proceeding with upstream patching"
else
  echo "ERROR: git status must be clean before applying upstream patches"
  exit 1
fi

scripts/third_party_pins/fabric/apply_upstream.sh

# The command above just copies a subset of Fabric files to /internal directory
# and applies proper headers. The rest of the process, described below, is about
# pathcing Fabric files so they compile and work locally.

# The first time upstream was patched in this repo, everything had to be done
# by hand. Each next time upstream is updated, we start by replaying the changes
# from the last patch, and proceed by resolving any conflicts and making any other
# changes to make all tests pass. In order for this process to work, it is
# necessary that each time we patch upstream we preserve the commit containing only
# the patch, so we can use it next time to replay the same changes.

# The first patch commit in this repo was 30f3fda5d02f0b07e4c3e9511ea9fe6c50d8bbba.
# The steps below describe how to replay changes from 30f3fda5d02f0b07e4c3e9511ea9fe6c50d8bbba
# on master at any time in the future.

# 1. Asuming we just ran 'make thirdparty-pin', we must commit upstream files first. This will
# keep the subsequent patch in its own clean commit so we can use it in the future.
# git add .
# git commit --signoff -m "Apply upstream packages"

# 2. Replay changes from the last correct patch. We do it using git, with the help of a
#    temporary branch where we first copy the changes we want to replay. The common
#    ancestor of the temporary branch has to be the commit which is the fist ancestor of
#    the last correct patch, in our case 935b01752c01cc18f15491d520337ede22eeaab5.
# git format-patch --stdout 935b01752c01cc18f15491d520337ede22eeaab5..30f3fda5d02f0b07e4c3e9511ea9fe6c50d8bbba > ~/last.patch
# git checkout -b fix 935b01752c01cc18f15491d520337ede22eeaab5
# git am ~/last.patch
# git checkout master
# git merge fix

# 3. If there are conflicts in step 2:

#    3.1 Fix conflicts from step 2. Make sure to run 'go mod tidy' as the last step.
# (fix conflicts to make all tests pass)
# go mod tidy

#    3.2 Complete the patch.
# git commit --signoff

# 4. Amend as required, and push all commits.
