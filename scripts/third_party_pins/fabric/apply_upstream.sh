#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# This script fetches code originating from other Hyperledger Fabric projects
# These files are checked into internal paths.
# Note: This script must be adjusted as upstream makes adjustments

set -e

UPSTREAM_PROJECT="github.com/hyperledger/fabric"
UPSTREAM_BRANCH="${UPSTREAM_BRANCH:-release}"
SCRIPTS_PATH="scripts/third_party_pins/fabric"
PATCHES_PATH="${SCRIPTS_PATH}/patches"

THIRDPARTY_INTERNAL_FABRIC_PATH='internal/github.com/hyperledger/fabric'

####
# Clone and patches packages into repo

# Clone original project into temporary directory
echo "Fetching upstream project ($UPSTREAM_PROJECT:$UPSTREAM_COMMIT) ..."
CWD=`pwd`
TMP=`mktemp -d 2>/dev/null || mktemp -d -t 'mytmpdir'`

TMP_PROJECT_PATH=$TMP/src/$UPSTREAM_PROJECT
mkdir -p $TMP_PROJECT_PATH
cd ${TMP_PROJECT_PATH}/..

git clone https://${UPSTREAM_PROJECT}.git
cd $TMP_PROJECT_PATH
git checkout $UPSTREAM_BRANCH
git reset --hard $UPSTREAM_COMMIT

cd $CWD

echo 'Removing current upstream project from working directory ...'
rm -Rf "${THIRDPARTY_INTERNAL_FABRIC_PATH}"
mkdir -p "${THIRDPARTY_INTERNAL_FABRIC_PATH}"

# Create internal utility structure
mkdir -p ${TMP_PROJECT_PATH}/internal/protoutil
cp -R ${TMP_PROJECT_PATH}/protoutil ${TMP_PROJECT_PATH}/internal/

# copy required files that are under internal into non-internal structure.
mkdir -p ${TMP_PROJECT_PATH}/libinternal
cp -R ${TMP_PROJECT_PATH}/internal/* ${TMP_PROJECT_PATH}/libinternal/

# fabric client utils
echo "Pinning and patching fabric client utils..."
declare -a CLIENT_UTILS_IMPORT_SUBSTS=(
    's/\"github.com\/hyperledger\/fabric\/internal/\"github.com\/trustbloc\/fabric-lib-go-ext\/internal\/github.com\/hyperledger\/fabric\/libinternal/g'
    's/[[:space:]]logging[[:space:]]\"github.com/\"github.com/g'
    's/\"github.com\/hyperledger\/fabric\/protos/\"github.com\/hyperledger\/fabric-protos-go/g'
    's/\"github.com\/hyperledger\/fabric\//\"github.com\/trustbloc\/fabric-lib-go-ext\/internal\/github.com\/hyperledger\/fabric\//g'
)

INTERNAL_PATH=$THIRDPARTY_INTERNAL_FABRIC_PATH TMP_PROJECT_PATH=$TMP_PROJECT_PATH IMPORT_SUBSTS="${CLIENT_UTILS_IMPORT_SUBSTS[*]}" $SCRIPTS_PATH/apply_fabric_client_utils.sh
INTERNAL_PATH=$THIRDPARTY_INTERNAL_FABRIC_PATH TMP_PROJECT_PATH=$TMP_PROJECT_PATH IMPORT_SUBSTS="${CLIENT_UTILS_IMPORT_SUBSTS[*]}" $SCRIPTS_PATH/apply_fabric_common_utils.sh

# external utils
echo "Pinning and patching fabric external utils ..."
declare -a EXTERNAL_UTILS_IMPORT_SUBSTS=(
    's/\"github.com\/hyperledger\/fabric\/protoutil/\"github.com\/trustbloc\/fabric-lib-go-ext\/internal\/github.com\/hyperledger\/fabric\/internal\/protoutil/g'
    's/\"github.com\/hyperledger\/fabric\/protos/\"github.com\/hyperledger\/fabric-protos-go/g'
    's/\"github.com\/hyperledger\/fabric\//\"github.com\/trustbloc\/fabric-lib-go-ext\/internal\/github.com\/hyperledger\/fabric\//g'
)
INTERNAL_PATH=$THIRDPARTY_INTERNAL_FABRIC_PATH TMP_PROJECT_PATH=$TMP_PROJECT_PATH IMPORT_SUBSTS="${EXTERNAL_UTILS_IMPORT_SUBSTS[*]}" $SCRIPTS_PATH/apply_fabric_external_utils.sh
INTERNAL_PATH=$THIRDPARTY_INTERNAL_FABRIC_PATH TMP_PROJECT_PATH=$TMP_PROJECT_PATH IMPORT_SUBSTS="${EXTERNAL_UTILS_IMPORT_SUBSTS[*]}" $SCRIPTS_PATH/apply_fabric_common_utils.sh

# Cleanup temporary files from patches application
echo "Removing temporary files ..."
rm -Rf $TMP