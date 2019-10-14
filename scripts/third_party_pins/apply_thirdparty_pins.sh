#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# This script fetches code originating from other upstream projects
# These files are checked into internal paths.

set -e

#if output=$(git status --porcelain) && [ -z "$output" ]
#then
#  echo "Working directory clean, proceeding with upstream patching"
#else
#  echo "ERROR: git status must be clean before applying upstream patches"
#  exit 1
#fi

scripts/third_party_pins/fabric/apply_upstream.sh

# Once upstream changes have been applied follow these steps:

# 1. Commit upstream changes in a commit separate from the patch which makes them work locally
# git add .
# git commit --signoff -m "Apply upstream packages"

# 2. Apply the last correct patch to upstream changes. The patch is not guaranteed to work with
#    the latest changes, so use -3 option to perform 3-way merge.
# git am -3 scripts/third_party_pins/patches/upstream.patch

# 3. If there are conflicts in step 2:

#    3.1 Fix conflicts from step 2. Make sure to run 'go mod tidy' as the last step.
# (fix conflicts to make all tests pass)
# go mod tidy

#    3.2 Complete the patch.
# git am --continue

#    3.3. Create the new correct patch to upstream changes.
# format-patch --stdout HEAD^..HEAD > scripts/third_party_pins/patches/upstream.patch

#    3.4 Commit the new correct patch to upstream changes.
# git add .
# git commit amend -m "Apply upstream packages"

# 4. Push all the changes.
#    IMPORTANT: DO NOT SQUASH! Commit from step 1 has to be present for
#    proper 3-way merge. Also notice that step 3.3 assumes that the last two commits
#    are from step 1 and step 3.2.
