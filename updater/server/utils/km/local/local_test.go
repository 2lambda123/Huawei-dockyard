/*
Copyright 2016 The ContainerOps Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package local

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	dus_utils "github.com/containerops/dockyard/updater/server/utils"
)

func loadTestData(t *testing.T) (dus_utils.DyKeyManager, string) {
	var local DyKeyManagerLocal
	_, path, _, _ := runtime.Caller(0)
	realPath := filepath.Join(filepath.Dir(path), "testdata")

	l, err := local.New(localPrefix + ":/" + realPath)
	assert.Nil(t, err, "Fail to setup a local test storage")

	return l, realPath
}

// TestBasic
func TestLocalBasic(t *testing.T) {
	var local DyKeyManagerLocal

	validURL := "local://tmp/containerops_km_cache"
	ok := local.Supported(validURL)
	assert.Equal(t, ok, true, "Fail to get supported status")
	ok = local.Supported("localInvalid://tmp/containerops_km_cache")
	assert.Equal(t, ok, false, "Fail to get supported status")

	_, err := local.New(validURL)
	assert.Nil(t, err, "Fail to setup a local storage")
}

func TestLocalGetPublicKey(t *testing.T) {
	l, realPath := loadTestData(t)
	key := "containerops/official"

	data, err := l.GetPublicKey(key)
	assert.Nil(t, err, "Fail to get public key")
	expectedData, _ := ioutil.ReadFile(filepath.Join(realPath, key, keyDir, defaultPublicKey))
	assert.Equal(t, data, expectedData, "Fail to get the correct public key data")
}

func TestLocalSign(t *testing.T) {
	l, realPath := loadTestData(t)
	key := "containerops/official"

	testFile := filepath.Join(realPath, "hello.txt")
	testByte, _ := ioutil.ReadFile(testFile)
	signFile := filepath.Join(realPath, "hello.sig")
	signByte, _ := ioutil.ReadFile(signFile)
	data, err := l.Sign(key, testByte)
	assert.Nil(t, err, "Fail to sign")
	assert.Equal(t, data, signByte, "Fail to sign correctly")
}
