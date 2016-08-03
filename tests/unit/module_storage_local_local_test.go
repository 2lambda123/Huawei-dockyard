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
package unittest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/containerops/dockyard/module"
	_ "github.com/containerops/dockyard/module/km/local"
	sl "github.com/containerops/dockyard/module/storage/local"
)

func loadSLTestData(t *testing.T) (module.UpdateServiceStorage, string) {
	var local sl.UpdateServiceStorageLocal

	_, path, _, _ := runtime.Caller(0)
	topPath := filepath.Join(filepath.Dir(path), "testdata")
	km := "local:/" + topPath
	// In this test, storage dir and key manager dir is the same
	l, err := local.New(km, km)
	assert.Nil(t, err, "Fail to setup a local test storage")

	return l, topPath
}

// TestBasic
func TestSLBasic(t *testing.T) {
	var local sl.UpdateServiceStorageLocal

	validURL := "local://tmp/containerops_storage_cache"
	ok := local.Supported(validURL)
	assert.Equal(t, ok, true, "Fail to get supported status")
	ok = local.Supported("localInvalid://tmp/containerops_storage_cache")
	assert.Equal(t, ok, false, "Fail to get supported status")

	l, err := local.New(validURL, "")
	assert.Nil(t, err, "Fail to setup a local storage")
	assert.Equal(t, l.String(), validURL)
}

func TestSLList(t *testing.T) {
	l, _ := loadSLTestData(t)
	key := "containerops/official"
	validCount := 0

	apps, _ := l.List("app/v1", key)
	for _, app := range apps {
		if app == "os/arch/appA" || app == "os/arch/appB" {
			validCount++
		}
	}
	assert.Equal(t, validCount, 2, "Fail to get right apps")
}

func TestSLPut(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "us-test-")
	defer os.RemoveAll(tmpPath)
	assert.Nil(t, err, "Fail to create temp dir")

	protocal := "app/v1"
	testData := "this is test DATA, you can put in anything here"

	var local sl.UpdateServiceStorageLocal
	l, err := local.New(sl.LocalPrefix+":/"+tmpPath, sl.LocalPrefix+":/"+tmpPath)
	assert.Nil(t, err, "Fail to setup local repo")

	invalidKey := "containerops/official"
	_, err = l.Put(protocal, invalidKey, []byte(testData))
	assert.NotNil(t, err, "Fail to put with invalid key")

	validKey := "containerops/official/os/arch/appA"
	_, err = l.Put(protocal, validKey, []byte(testData))
	assert.Nil(t, err, "Fail to put key")

	_, err = l.GetMeta(protocal, "containerops/official")
	assert.Nil(t, err, "Fail to get meta data")

	getData, err := l.Get(protocal, validKey)
	assert.Nil(t, err, "Fail to load file")
	assert.Equal(t, string(getData), testData, "Fail to get correct file")
}

func TestSLGet(t *testing.T) {
	l, kmPath := loadSLTestData(t)

	protocal := "app/v1"
	key := "containerops/official"
	invalidKey := "containerops/official/invalid"

	defer os.RemoveAll(filepath.Join(kmPath, key, "key"))
	_, err := l.GetPublicKey(protocal, key)
	assert.Nil(t, err, "Fail to load public key")
	_, err = l.GetMetaSign(protocal, key)
	assert.Nil(t, err, "Fail to load  sign file")

	_, err = l.GetMeta(protocal, invalidKey)
	assert.NotNil(t, err, "Fail to get meta from invalid key")
	_, err = l.GetMeta(protocal, key)
	assert.Nil(t, err, "Fail to load meta data")

	_, err = l.Get(protocal, "invalidinput")
	assert.NotNil(t, err, "Fail to get by invalid key")

	data, err := l.Get(protocal, key+"/os/arch/appA")
	expectedData := "This is the content of appA."
	assert.Nil(t, err, "Fail to load file")
	assert.Equal(t, string(data), expectedData, "Fail to get correct file")
}
