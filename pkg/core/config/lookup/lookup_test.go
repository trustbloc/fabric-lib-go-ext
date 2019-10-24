/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lookup

import (
	"testing"

	"github.com/trustbloc/fabric-lib-go-ext/pkg/core/config"

	"os"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trustbloc/fabric-lib-go-ext/pkg/common/providers/core"
)

var backend *config.MapConfigBackend

func TestMain(m *testing.M) {
	backend = setupCustomBackend("key")
	r := m.Run()
	os.Exit(r)
}

func TestGetBool(t *testing.T) {
	//Test single backend lookup
	testLookup := New(backend)
	assert.True(t, testLookup.GetBool("key.bool.true"), "expected lookup to return true")
	assert.False(t, testLookup.GetBool("key.bool.false"), "expected lookup to return false")
	assert.False(t, testLookup.GetBool("key.bool.invalid"), "expected lookup to return false for invalid value")
	assert.False(t, testLookup.GetBool("key.bool.notexisting"), "expected lookup to return false for not existing value")

	//Test With multiple backend
	keyPrefixes := []string{"key1", "key2", "key3", "key4"}
	backends := getMultipleCustomBackends(keyPrefixes)
	testLookup = New(backends...)

	for _, prefix := range keyPrefixes {
		assert.True(t, testLookup.GetBool(prefix+".bool.true"), "expected lookup to return true")
		assert.False(t, testLookup.GetBool(prefix+".bool.false"), "expected lookup to return false")
		assert.False(t, testLookup.GetBool(prefix+".bool.invalid"), "expected lookup to return false for invalid value")
		assert.False(t, testLookup.GetBool(prefix+".bool.notexisting"), "expected lookup to return false for not existing value")
	}
}

func TestGetInt(t *testing.T) {
	testLookup := New(backend)
	assert.True(t, testLookup.GetInt("key.int.positive") == 5, "expected lookup to return valid positive value")
	assert.True(t, testLookup.GetInt("key.int.negative") == -5, "expected lookup to return valid negative value")
	assert.True(t, testLookup.GetInt("key.int.invalid") == 0, "expected lookup to return 0")
	assert.True(t, testLookup.GetInt("key.int.not.existing") == 0, "expected lookup to return 0")

	//Test With multiple backend
	keyPrefixes := []string{"key1", "key2", "key3", "key4"}
	backends := getMultipleCustomBackends(keyPrefixes)
	testLookup = New(backends...)

	for _, prefix := range keyPrefixes {
		assert.True(t, testLookup.GetInt(prefix+".int.positive") == 5, "expected lookup to return valid positive value")
		assert.True(t, testLookup.GetInt(prefix+".int.negative") == -5, "expected lookup to return valid negative value")
		assert.True(t, testLookup.GetInt(prefix+".int.invalid") == 0, "expected lookup to return 0")
		assert.True(t, testLookup.GetInt(prefix+".int.not.existing") == 0, "expected lookup to return 0")
	}
}

func TestGetString(t *testing.T) {
	testLookup := New(backend)
	assert.True(t, testLookup.GetString("key.string.valid") == "valid-string", "expected lookup to return valid string value")
	assert.True(t, testLookup.GetString("key.string.valid.lower.case") == "valid-string", "expected lookup to return valid string value")
	assert.True(t, testLookup.GetString("key.string.valid.upper.case") == "VALID-STRING", "expected lookup to return valid string value")
	assert.True(t, testLookup.GetString("key.string.valid.mixed.case") == "VaLiD-StRiNg", "expected lookup to return valid string value")
	assert.True(t, testLookup.GetString("key.string.empty") == "", "expected lookup to return empty string value")
	assert.True(t, testLookup.GetString("key.string.nil") == "", "expected lookup to return empty string value")
	assert.True(t, testLookup.GetString("key.string.number") == "1234", "expected lookup to return valid string value")
	assert.True(t, testLookup.GetString("key.string.not existing") == "", "expected lookup to return empty string value")

	//Test With multiple backend
	keyPrefixes := []string{"key1", "key2", "key3", "key4"}
	backends := getMultipleCustomBackends(keyPrefixes)
	testLookup = New(backends...)

	for _, prefix := range keyPrefixes {
		assert.True(t, testLookup.GetString(prefix+".string.valid") == "valid-string", "expected lookup to return valid string value")
		assert.True(t, testLookup.GetString(prefix+".string.valid.lower.case") == "valid-string", "expected lookup to return valid string value")
		assert.True(t, testLookup.GetString(prefix+".string.valid.upper.case") == "VALID-STRING", "expected lookup to return valid string value")
		assert.True(t, testLookup.GetString(prefix+".string.valid.mixed.case") == "VaLiD-StRiNg", "expected lookup to return valid string value")
		assert.True(t, testLookup.GetString(prefix+".string.empty") == "", "expected lookup to return empty string value")
		assert.True(t, testLookup.GetString(prefix+".string.nil") == "", "expected lookup to return empty string value")
		assert.True(t, testLookup.GetString(prefix+".string.number") == "1234", "expected lookup to return valid string value")
		assert.True(t, testLookup.GetString(prefix+".string.not existing") == "", "expected lookup to return empty string value")
	}
}

func TestGetLowerString(t *testing.T) {
	testLookup := New(backend)
	assert.True(t, testLookup.GetLowerString("key.string.valid") == "valid-string", "expected lookup to return valid lowercase string value")
	assert.True(t, testLookup.GetLowerString("key.string.valid.lower.case") == "valid-string", "expected lookup to return valid lowercase string value")
	assert.True(t, testLookup.GetLowerString("key.string.valid.upper.case") == "valid-string", "expected lookup to return valid lowercase string value")
	assert.True(t, testLookup.GetLowerString("key.string.valid.mixed.case") == "valid-string", "expected lookup to return valid lowercase string value")
	assert.True(t, testLookup.GetLowerString("key.string.empty") == "", "expected lookup to return empty string value")
	assert.True(t, testLookup.GetLowerString("key.string.nil") == "", "expected lookup to return empty string value")
	assert.True(t, testLookup.GetLowerString("key.string.number") == "1234", "expected lookup to return valid string value")
	assert.True(t, testLookup.GetLowerString("key.string.not existing") == "", "expected lookup to return empty string value")

	//Test With multiple backends
	keyPrefixes := []string{"key1", "key2", "key3", "key4"}
	backends := getMultipleCustomBackends(keyPrefixes)
	testLookup = New(backends...)

	for _, prefix := range keyPrefixes {
		assert.True(t, testLookup.GetLowerString(prefix+".string.valid") == "valid-string", "expected lookup to return valid lowercase string value")
		assert.True(t, testLookup.GetLowerString(prefix+".string.valid.lower.case") == "valid-string", "expected lookup to return valid lowercase string value")
		assert.True(t, testLookup.GetLowerString(prefix+".string.valid.upper.case") == "valid-string", "expected lookup to return valid lowercase string value")
		assert.True(t, testLookup.GetLowerString(prefix+".string.valid.mixed.case") == "valid-string", "expected lookup to return valid lowercase string value")
		assert.True(t, testLookup.GetLowerString(prefix+".string.empty") == "", "expected lookup to return empty string value")
		assert.True(t, testLookup.GetLowerString(prefix+".string.nil") == "", "expected lookup to return empty string value")
		assert.True(t, testLookup.GetLowerString(prefix+".string.number") == "1234", "expected lookup to return valid string value")
		assert.True(t, testLookup.GetLowerString(prefix+".string.not existing") == "", "expected lookup to return empty string value")
	}
}

func TestGetDuration(t *testing.T) {
	testLookup := New(backend)
	assert.True(t, testLookup.GetDuration("key.duration.valid.hour").String() == (24*time.Hour).String(), "expected valid time value")
	assert.True(t, testLookup.GetDuration("key.duration.valid.minute").String() == (24*time.Minute).String(), "expected valid time value")
	assert.True(t, testLookup.GetDuration("key.duration.valid.second").String() == (24*time.Second).String(), "expected valid time value")
	assert.True(t, testLookup.GetDuration("key.duration.valid.millisecond").String() == (24*time.Millisecond).String(), "expected valid time value")
	assert.True(t, testLookup.GetDuration("key.duration.valid.microsecond").String() == (24*time.Microsecond).String(), "expected valid time value")
	assert.True(t, testLookup.GetDuration("key.duration.valid.nanosecond").String() == (24*time.Nanosecond).String(), "expected valid time value")
	//default value tests
	assert.True(t, testLookup.GetDuration("key.duration.valid.not.existing").String() == (0*time.Second).String(), "expected valid default time value")
	assert.True(t, testLookup.GetDuration("key.duration.valid.invalid").String() == (0*time.Second).String(), "expected valid  default time value")
	assert.True(t, testLookup.GetDuration("key.duration.valid.nil").String() == (0*time.Second).String(), "expected valid  default time value")
	assert.True(t, testLookup.GetDuration("key.duration.valid.empty").String() == (0*time.Second).String(), "expected valid  default time value")
	//default when no time unit provided
	assert.True(t, testLookup.GetDuration("key.duration.valid.no.unit").String() == (12*time.Nanosecond).String(), "expected valid default time value with default unit")

	//Test With multiple backends
	keyPrefixes := []string{"key1", "key2", "key3", "key4"}
	backends := getMultipleCustomBackends(keyPrefixes)
	testLookup = New(backends...)

	for _, prefix := range keyPrefixes {
		assert.True(t, testLookup.GetDuration(prefix+".duration.valid.hour").String() == (24*time.Hour).String(), "expected valid time value")
		assert.True(t, testLookup.GetDuration(prefix+".duration.valid.minute").String() == (24*time.Minute).String(), "expected valid time value")
		assert.True(t, testLookup.GetDuration(prefix+".duration.valid.second").String() == (24*time.Second).String(), "expected valid time value")
		assert.True(t, testLookup.GetDuration(prefix+".duration.valid.millisecond").String() == (24*time.Millisecond).String(), "expected valid time value")
		assert.True(t, testLookup.GetDuration(prefix+".duration.valid.microsecond").String() == (24*time.Microsecond).String(), "expected valid time value")
		assert.True(t, testLookup.GetDuration(prefix+".duration.valid.nanosecond").String() == (24*time.Nanosecond).String(), "expected valid time value")
		//default value tests
		assert.True(t, testLookup.GetDuration(prefix+".duration.valid.not.existing").String() == (0*time.Second).String(), "expected valid default time value")
		assert.True(t, testLookup.GetDuration(prefix+".duration.valid.invalid").String() == (0*time.Second).String(), "expected valid  default time value")
		assert.True(t, testLookup.GetDuration(prefix+".duration.valid.nil").String() == (0*time.Second).String(), "expected valid  default time value")
		assert.True(t, testLookup.GetDuration(prefix+".duration.valid.empty").String() == (0*time.Second).String(), "expected valid  default time value")
		//default when no time unit provided
		assert.True(t, testLookup.GetDuration(prefix+".duration.valid.no.unit").String() == (12*time.Nanosecond).String(), "expected valid default time value with default unit")
	}
}

func setupCustomBackend(keyPrefix string) *config.MapConfigBackend {

	backendMap := make(map[string]interface{})

	backendMap[keyPrefix+".bool.true"] = true
	backendMap[keyPrefix+".bool.false"] = false
	backendMap[keyPrefix+".bool.invalid"] = "INVALID"

	backendMap[keyPrefix+".int.positive"] = 5
	backendMap[keyPrefix+".int.negative"] = -5
	backendMap[keyPrefix+".int.invalid"] = "INVALID"

	backendMap[keyPrefix+".string.valid"] = "valid-string"
	backendMap[keyPrefix+".string.valid.mixed.case"] = "VaLiD-StRiNg"
	backendMap[keyPrefix+".string.valid.lower.case"] = "valid-string"
	backendMap[keyPrefix+".string.valid.upper.case"] = "VALID-STRING"
	backendMap[keyPrefix+".string.empty"] = ""
	backendMap[keyPrefix+".string.nil"] = nil
	backendMap[keyPrefix+".string.number"] = 1234

	backendMap[keyPrefix+".duration.valid.hour"] = "24h"
	backendMap[keyPrefix+".duration.valid.minute"] = "24m"
	backendMap[keyPrefix+".duration.valid.second"] = "24s"
	backendMap[keyPrefix+".duration.valid.millisecond"] = "24ms"
	backendMap[keyPrefix+".duration.valid.microsecond"] = "24µs"
	backendMap[keyPrefix+".duration.valid.nanosecond"] = "24ns"
	backendMap[keyPrefix+".duration.valid.no.unit"] = "12"
	backendMap[keyPrefix+".duration.invalid"] = "24XYZ"
	backendMap[keyPrefix+".duration.nil"] = nil
	backendMap[keyPrefix+".duration.empty"] = ""

	return config.NewMapConfigBackend(backendMap)
}

func getMultipleCustomBackends(keyPrefixes []string) []core.ConfigBackend {
	var backends []core.ConfigBackend
	for _, prefix := range keyPrefixes {
		backends = append(backends, setupCustomBackend(prefix))
	}
	return backends
}
