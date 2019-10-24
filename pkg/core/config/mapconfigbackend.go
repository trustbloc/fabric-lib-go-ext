/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package config

//MapConfigBackend is a simple map implementation of ConfigBackend
type MapConfigBackend struct {
	m map[string]interface{}
}

//NewMapConfigBackend creates a new instance of MapConfigBackend
func NewMapConfigBackend(m map[string]interface{}) *MapConfigBackend {
	mcb := &MapConfigBackend{
		m: m,
	}
	return mcb
}

//Lookup returns or unmarshals value for given key
func (b *MapConfigBackend) Lookup(key string) (interface{}, bool) {
	v, ok := b.m[key]
	return v, ok
}
