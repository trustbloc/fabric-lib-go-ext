#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# This script pins client and common package families from Hyperledger Fabric into this project.
# These files are checked into internal paths.
# Note: This script must be adjusted as upstream makes adjustments

set -e

# Create and populate patching directory.
declare TMP=`mktemp -d 2>/dev/null || mktemp -d -t 'mytmpdir'`
declare PATCH_PROJECT_PATH=$TMP/src/$UPSTREAM_PROJECT
cp -R ${TMP_PROJECT_PATH} ${PATCH_PROJECT_PATH}
declare TMP_PROJECT_PATH=${PATCH_PROJECT_PATH}

declare -a PKGS=(

    "common/cauthdsl"
    "core/ledger/kvledger/txmgmt/rwsetutil"

    "common/crypto"
    "common/errors"
    "common/genesis"
    "common/util"
    "common/capabilities"
    "common/channelconfig"
    "common/configtx"
    "common/policies"
    "common/ledger"

    "common/tools/protolator"
    "common/tools/protolator/protoext"
    "common/tools/protolator/protoext/commonext"
    "common/tools/protolator/protoext/ledger/rwsetext"
    "common/tools/protolator/protoext/mspext"
    "common/tools/protolator/protoext/ordererext"
    "common/tools/protolator/protoext/peerext"

    "common/viperutil"

    "core/config"

    "core/ledger/util"

    "msp"

    "protoutil"

    "libinternal/configtxgen/encoder"
    "libinternal/configtxgen/localconfig"
    "libinternal/configtxlator/update"

    "libinternal/pkg/identity"

)

declare -a FILES=(

    "common/cauthdsl/cauthdsl_builder.go"
    "common/cauthdsl/policyparser.go"

    "core/ledger/kvledger/txmgmt/rwsetutil/rwset_proto_util.go"
    "core/ledger/util/txvalidationflags.go"

    "common/configtx/configtx.go"


    "common/capabilities/application.go"
    "common/capabilities/capabilities.go"
    "common/capabilities/channel.go"
    "common/capabilities/orderer.go"

    "common/genesis/genesis.go"

    "common/crypto/random.go"

    "common/channelconfig/application.go"
    "common/channelconfig/consortium.go"
    "common/channelconfig/consortiums.go"
    "common/channelconfig/applicationorg.go"
    "common/channelconfig/channel.go"
    "common/channelconfig/util.go"
    "common/channelconfig/orderer.go"
    "common/channelconfig/organization.go"
    "common/channelconfig/api.go"
    "common/channelconfig/standardvalues.go"
    "common/channelconfig/acls.go"
    "common/channelconfig/bundle.go"

    "common/policies/policy.go"
    "common/policies/util.go"
    "common/policies/implicitmetaparser.go"

    "common/tools/protolator/api.go"
    "common/tools/protolator/dynamic.go"
    "common/tools/protolator/json.go"
    "common/tools/protolator/nested.go"
    "common/tools/protolator/statically_opaque.go"
    "common/tools/protolator/variably_opaque.go"
    "common/tools/protolator/protoext/decorate.go"
    "common/tools/protolator/protoext/commonext/common.go"
    "common/tools/protolator/protoext/commonext/configtx.go"
    "common/tools/protolator/protoext/commonext/configuration.go"
    "common/tools/protolator/protoext/commonext/policies.go"
    "common/tools/protolator/protoext/ledger/rwsetext/rwset.go"
    "common/tools/protolator/protoext/mspext/msp_config.go"
    "common/tools/protolator/protoext/mspext/msp_principal.go"
    "common/tools/protolator/protoext/ordererext/configuration.go"
    "common/tools/protolator/protoext/peerext/configuration.go"
    "common/tools/protolator/protoext/peerext/proposal.go"
    "common/tools/protolator/protoext/peerext/proposal_response.go"
    "common/tools/protolator/protoext/peerext/transaction.go"

    "common/viperutil/config_util.go"

    "core/config/config.go"

    "common/util/utils.go"

    "msp/factory.go"
    "msp/configbuilder.go"
    "msp/msp.go"

    "protoutil/blockutils.go"
    "protoutil/commonutils.go"
    "protoutil/configtxutils.go"
    "protoutil/proputils.go"
    "protoutil/signeddata.go"
    "protoutil/txutils.go"
    "protoutil/configtxutils.go"
    "protoutil/unmarshalers.go"
    "protoutil/commonutils.go"

    "libinternal/configtxgen/encoder/encoder.go"
    "libinternal/configtxgen/localconfig/config.go"
    "libinternal/configtxlator/update/update.go"
    "libinternal/configtxlator/update/update.go"

    "libinternal/pkg/identity/identity.go"

)

# Create directory structure for packages
for i in "${PKGS[@]}"
do
    mkdir -p $INTERNAL_PATH/${i}
done

echo "Modifying go source files"

# Apply patching
echo "Inserting modification notice ..."
WORKING_DIR=$TMP_PROJECT_PATH FILES="${FILES[@]}" scripts/third_party_pins/common/apply_header_notice.sh

# Copy patched project into internal paths
echo "Copying patched upstream project into working directory ..."
for i in "${FILES[@]}"
do
    TARGET_PATH=`dirname $INTERNAL_PATH/${i}`
    cp $TMP_PROJECT_PATH/${i} $TARGET_PATH
done

rm -Rf ${TMP_PROJECT_PATH}