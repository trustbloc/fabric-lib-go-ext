/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
/*
Notice: This file has been modified for TrustBloc Fabric Lib Go EXT usage.
Please review third_party pinning scripts and patches for more details.
*/

package channelconfig

import (
	flogging "github.com/trustbloc/fabric-lib-go-ext/internal/github.com/hyperledger/fabric/libpatch/logbridge"
)

var logger = flogging.MustGetLogger("common.channelconfig")

// RootGroupKey is the key for namespacing the channel config, especially for
// policy evaluation.
const RootGroupKey = "Channel"
