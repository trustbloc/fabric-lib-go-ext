/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
/*
Notice: This file has been modified for TrustBloc Fabric Lib Go EXT usage.
Please review third_party pinning scripts and patches for more details.
*/

package msp

// From this point on, there are interfaces that are shared within the peer and client API
// of the membership service provider.

// ProviderType indicates the type of an identity provider
type ProviderType int

// The ProviderType of a member relative to the member API
const (
	FABRIC ProviderType = iota // MSP is of FABRIC type
	IDEMIX                     // MSP is of IDEMIX type
	OTHER                      // MSP is of OTHER TYPE

	// NOTE: as new types are added to this set,
	// the mspTypes map below must be extended
)

var mspTypeStrings = map[ProviderType]string{
	FABRIC: "bccsp",
	IDEMIX: "idemix",
}

// ProviderTypeToString returns a string that represents the ProviderType integer
func ProviderTypeToString(id ProviderType) string {
	if res, found := mspTypeStrings[id]; found {
		return res
	}

	return ""
}

const (
	// SHA2 is an identifier for SHA2 hash family
	SHA2 = "SHA2"

	// SHA256
	SHA256 = "SHA256"
)