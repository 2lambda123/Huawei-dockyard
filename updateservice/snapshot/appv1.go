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
package snapshot

import (
	"crypto/md5"
	"errors"
	"io/ioutil"
)

const (
	appv1Proto = "appv1"
)

type UpdateServiceSnapshotAppv1 struct {
	ID  string
	URL string

	callback Callback
}

func (m *UpdateServiceSnapshotAppv1) New(id, url string, callback Callback) (UpdateServiceSnapshot, error) {
	if id == "" || url == "" {
		return nil, errors.New("id|url should not be empty")
	}

	m.ID, m.URL, m.callback = id, url, callback
	return m, nil
}

func (m *UpdateServiceSnapshotAppv1) Supported(proto string) bool {
	return proto == appv1Proto
}

func (m *UpdateServiceSnapshotAppv1) Process() error {
	var data UpdateServiceSnapshotOutput

	content, err := ioutil.ReadFile(m.URL)
	if m.callback == nil {
		return err
	}

	if err == nil {
		s := md5.Sum(content)
		data.Data = s[:]
	}
	data.Error = err

	return m.callback(m.ID, data)
}

func (m *UpdateServiceSnapshotAppv1) Description() string {
	return "Scan the appv1 package, return its md5"
}
