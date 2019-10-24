/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package core

//CryptoSuiteConfig contains configuration items for cryptosuite.
type CryptoSuiteConfig interface {
	IsSecurityEnabled() bool
	SecurityAlgorithm() string
	SecurityLevel() int
	SecurityProvider() string
	SoftVerify() bool
	SecurityProviderLibPath() string
	SecurityProviderPin() string
	SecurityProviderLabel() string
	KeyStorePath() string
}

// Providers represents the configured core providers context.
type Providers interface {
	CryptoSuite() CryptoSuite
	SigningManager() SigningManager
}

//ConfigProvider provides config backend
type ConfigProvider func() ([]ConfigBackend, error)

//ConfigBackend backend for all config types
type ConfigBackend interface {
	Lookup(key string) (interface{}, bool)
}
