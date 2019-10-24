/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cryptosuite

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
	"github.com/trustbloc/fabric-lib-go-ext/pkg/common/providers/core"
	"github.com/trustbloc/fabric-lib-go-ext/pkg/core/config/lookup"
	"github.com/trustbloc/fabric-lib-go-ext/pkg/util/pathvar"
)

const (
	defEnabled       = true
	defHashAlgorithm = "SHA2"
	defLevel         = 256
	defProvider      = "SW"
	defSoftVerify    = true

	keySecurityEnabled         = "client.BCCSP.security.enabled"
	keySecurityLevel           = "client.BCCSP.security.level"
	keySecurityDefaultProvider = "client.BCCSP.security.default.provider"
	keySecurityHashAlgorithm   = "client.BCCSP.security.hashAlgorithm"
	keySecuritySoftVerify      = "client.BCCSP.security.softVerify"
	keySecurityLibrary         = "client.BCCSP.security.library"
	keySecurityPin             = "client.BCCSP.security.pin"
	keySecurityLabel           = "client.BCCSP.security.label"
	keySecurityEphemeral       = "client.BCCSP.security.ephemeral"
	keyCryptoStorePath         = "client.credentialStore.cryptoStore.path"
)

// ConfigFromBackend returns CryptoSuite config implementation for given backend
func ConfigFromBackend(coreBackend ...core.ConfigBackend) core.CryptoSuiteConfig {
	return &Config{backend: lookup.New(coreBackend...)}
}

// Config represents the crypto suite configuration for the client
type Config struct {
	backend *lookup.ConfigLookup
}

// IsSecurityEnabled config used enable and disable security in cryptosuite
func (c *Config) IsSecurityEnabled() bool {
	val, ok := c.backend.Lookup(keySecurityEnabled)
	if !ok {
		return defEnabled
	}
	return cast.ToBool(val)
}

// SecurityAlgorithm returns cryptoSuite config hash algorithm
func (c *Config) SecurityAlgorithm() string {
	val, ok := c.backend.Lookup(keySecurityHashAlgorithm)
	if !ok {
		return defHashAlgorithm
	}
	return cast.ToString(val)
}

// SecurityLevel returns cryptSuite config security level
func (c *Config) SecurityLevel() int {
	val, ok := c.backend.Lookup(keySecurityLevel)
	if !ok {
		return defLevel
	}
	return cast.ToInt(val)
}

//SecurityProvider provider SW or PKCS11
func (c *Config) SecurityProvider() string {
	val, ok := c.backend.Lookup(keySecurityDefaultProvider)
	if !ok {
		return strings.ToLower(defProvider)
	}
	return strings.ToLower(cast.ToString(val))
}

//SoftVerify flag
func (c *Config) SoftVerify() bool {
	val, ok := c.backend.Lookup(keySecuritySoftVerify)
	if !ok {
		return defSoftVerify
	}
	return cast.ToBool(val)
}

//SecurityProviderLibPath will be set only if provider is PKCS11
func (c *Config) SecurityProviderLibPath() string {
	configuredLibs := c.backend.GetString(keySecurityLibrary)
	libPaths := strings.Split(configuredLibs, ",")
	logger.Debugf("Configured BCCSP Lib Paths %s", libPaths)
	var lib string
	for _, path := range libPaths {
		if _, err := os.Stat(strings.TrimSpace(path)); !os.IsNotExist(err) {
			lib = strings.TrimSpace(path)
			break
		}
	}
	if lib != "" {
		logger.Debugf("Found softhsm library: %s", lib)
	} else {
		logger.Debug("Softhsm library was not found")
	}
	return lib
}

//SecurityProviderPin will be set only if provider is PKCS11
func (c *Config) SecurityProviderPin() string {
	return c.backend.GetString(keySecurityPin)
}

//SecurityProviderLabel will be set only if provider is PKCS11
func (c *Config) SecurityProviderLabel() string {
	return c.backend.GetString(keySecurityLabel)
}

// KeyStorePath returns the keystore path used by BCCSP
func (c *Config) KeyStorePath() string {
	keystorePath := pathvar.Subst(c.backend.GetString(keyCryptoStorePath))
	return filepath.Join(keystorePath, "keystore")
}
