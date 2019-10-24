/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cryptosuite

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/trustbloc/fabric-lib-go-ext/pkg/core/config"

	"github.com/stretchr/testify/assert"
	"github.com/trustbloc/fabric-lib-go-ext/pkg/common/providers/core"
)

func getConfigTestBackend() *config.MapConfigBackend {
	m := map[string]interface{}{
		keySecurityEnabled:         true,
		keySecurityHashAlgorithm:   "SHA2",
		keySecurityDefaultProvider: "SW",
		keySecuritySoftVerify:      true,
	}
	return config.NewMapConfigBackend(m)
}

func TestEmptyConfig(t *testing.T) {
	emptyMap := map[string]interface{}{}
	backend := config.NewMapConfigBackend(emptyMap)

	cryptoConfig := ConfigFromBackend(backend).(*Config)

	// Test for defaults
	assert.Equal(t, true, cryptoConfig.IsSecurityEnabled())
	assert.Equal(t, "SHA2", cryptoConfig.SecurityAlgorithm())
	assert.Equal(t, 256, cryptoConfig.SecurityLevel())
	// Note that we transform to lower case in SecurityProvider()
	assert.Equal(t, "sw", cryptoConfig.SecurityProvider())
	assert.Equal(t, true, cryptoConfig.SoftVerify())
}

func TestConfigKeyStorePath(t *testing.T) {

	customBackend := getCustomBackend(getConfigTestBackend())
	cryptoConfig := ConfigFromBackend(customBackend).(*Config)

	// Test KeyStore Path
	val, ok := customBackend.Lookup(keyCryptoStorePath)
	if !ok || val == nil {
		t.Fatal("expected valid value")
	}

	if filepath.Join(val.(string), "keystore") != cryptoConfig.KeyStorePath() {
		t.Fatal("Incorrect keystore path")
	}
}

func TestConfigBCCSPSecurityEnabled(t *testing.T) {

	customBackend := getCustomBackend(getConfigTestBackend())
	cryptoConfig := ConfigFromBackend(customBackend).(*Config)

	// Test BCCSP security is enabled
	val, ok := customBackend.Lookup(keySecurityEnabled)
	if !ok || val == nil {
		t.Fatal("expected valid value")
	}
	if val.(bool) != cryptoConfig.IsSecurityEnabled() {
		t.Fatal("Incorrect BCCSP Security enabled flag")
	}
}

func TestConfigSecurityAlgorithm(t *testing.T) {

	customBackend := getCustomBackend(getConfigTestBackend())
	cryptoConfig := ConfigFromBackend(customBackend).(*Config)

	// Test SecurityAlgorithm
	val, ok := customBackend.Lookup(keySecurityHashAlgorithm)
	if !ok || val == nil {
		t.Fatal("expected valid value")
	}
	if val.(string) != cryptoConfig.SecurityAlgorithm() {
		t.Fatal("Incorrect BCCSP Security Hash algorithm")
	}
}

func TestConfigSecurityLevel(t *testing.T) {

	customBackend := getCustomBackend(getConfigTestBackend())
	cryptoConfig := ConfigFromBackend(customBackend).(*Config)

	// Test Security Level
	val, ok := customBackend.Lookup(keySecurityLevel)
	if !ok || val == nil {
		t.Fatal("expected valid value")
	}
	if val.(int) != cryptoConfig.SecurityLevel() {
		t.Fatal("Incorrect BCCSP Security Level")
	}
}

func TestConfigSecurityProvider(t *testing.T) {

	customBackend := getCustomBackend(getConfigTestBackend())
	cryptoConfig := ConfigFromBackend(customBackend).(*Config)

	// Test SecurityProvider provider
	val, ok := customBackend.Lookup(keySecurityDefaultProvider)
	if !ok || val == nil {
		t.Fatal("expected valid value")
	}
	if !strings.EqualFold(val.(string), cryptoConfig.SecurityProvider()) {
		t.Fatalf("Incorrect BCCSP SecurityProvider provider : %s", cryptoConfig.SecurityProvider())
	}
}

func TestConfigSoftVerifyFlag(t *testing.T) {

	customBackend := getCustomBackend(getConfigTestBackend())
	cryptoConfig := ConfigFromBackend(customBackend).(*Config)

	// Test SoftVerify flag
	val, ok := customBackend.Lookup(keySecuritySoftVerify)
	if !ok || val == nil {
		t.Fatal("expected valid value")
	}
	if val.(bool) != cryptoConfig.SoftVerify() {
		t.Fatal("Incorrect BCCSP Ephemeral flag")
	}
}

func TestConfigSecurityProviderPin(t *testing.T) {

	customBackend := getCustomBackend(getConfigTestBackend())
	cryptoConfig := ConfigFromBackend(customBackend).(*Config)

	// Test SecurityProviderPin
	val, ok := customBackend.Lookup(keySecurityPin)
	if !ok || val == nil {
		t.Fatal("expected valid value")
	}
	if val.(string) != cryptoConfig.SecurityProviderPin() {
		t.Fatal("Incorrect BCCSP SecurityProviderPin flag")
	}
}

func TestConfigSecurityProviderLabel(t *testing.T) {

	customBackend := getCustomBackend(getConfigTestBackend())
	cryptoConfig := ConfigFromBackend(customBackend).(*Config)

	// Test SecurityProviderLabel
	val, ok := customBackend.Lookup(keySecurityLabel)
	if !ok || val == nil {
		t.Fatal("expected valid value")
	}
	if val.(string) != cryptoConfig.SecurityProviderLabel() {
		t.Fatal("Incorrect BCCSP SecurityProviderPin flag")
	}
}

func TestCryptoConfigWithMultipleBackends(t *testing.T) {
	var backends []core.ConfigBackend
	backendMap := make(map[string]interface{})
	backendMap[keySecurityEnabled] = true
	backends = append(backends, config.NewMapConfigBackend(backendMap))

	backendMap = make(map[string]interface{})
	backendMap[keySecurityHashAlgorithm] = "SHA2"
	backends = append(backends, config.NewMapConfigBackend(backendMap))

	backendMap = make(map[string]interface{})
	backendMap[keySecurityDefaultProvider] = "PKCS11"
	backends = append(backends, config.NewMapConfigBackend(backendMap))

	backendMap = make(map[string]interface{})
	backendMap[keySecurityLevel] = 2
	backends = append(backends, config.NewMapConfigBackend(backendMap))

	backendMap = make(map[string]interface{})
	backendMap[keySecurityPin] = "1234"
	backends = append(backends, config.NewMapConfigBackend(backendMap))

	backendMap = make(map[string]interface{})
	backendMap[keyCryptoStorePath] = "/tmp"
	backends = append(backends, config.NewMapConfigBackend(backendMap))

	backendMap = make(map[string]interface{})
	backendMap[keySecurityLabel] = "TESTLABEL"
	backends = append(backends, config.NewMapConfigBackend(backendMap))

	cryptoConfig := ConfigFromBackend(backends...)

	assert.Equal(t, cryptoConfig.IsSecurityEnabled(), true)
	assert.Equal(t, cryptoConfig.SecurityAlgorithm(), "SHA2")
	assert.Equal(t, cryptoConfig.SecurityProvider(), "pkcs11")
	assert.Equal(t, cryptoConfig.SecurityLevel(), 2)
	assert.Equal(t, cryptoConfig.SecurityProviderPin(), "1234")
	assert.Equal(t, cryptoConfig.KeyStorePath(), "/tmp/keystore")
	assert.Equal(t, cryptoConfig.SecurityProviderLabel(), "TESTLABEL")
}

//getCustomBackend returns custom backend to override config values and to avoid using new config file for test scenarios
func getCustomBackend(configBackend ...core.ConfigBackend) *config.MapConfigBackend {
	backendMap := make(map[string]interface{})
	backendMap[keySecurityEnabled], _ = configBackend[0].Lookup(keySecurityEnabled)
	backendMap[keySecurityHashAlgorithm], _ = configBackend[0].Lookup(keySecurityHashAlgorithm)
	backendMap[keySecurityDefaultProvider], _ = configBackend[0].Lookup(keySecurityDefaultProvider)
	backendMap[keySecurityEphemeral], _ = configBackend[0].Lookup(keySecurityEphemeral)
	backendMap[keySecuritySoftVerify], _ = configBackend[0].Lookup(keySecuritySoftVerify)
	backendMap[keySecurityLevel] = 2
	backendMap[keySecurityPin] = "1234"
	backendMap[keyCryptoStorePath] = "/tmp"
	backendMap[keySecurityLabel] = "TESTLABEL"
	return config.NewMapConfigBackend(backendMap)
}
