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

package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	dyutils "github.com/containerops/dockyard/utils"
)

var (
	ErrorsDUCConfigExist  = errors.New("dockyard update client configuration is already exist")
	ErrorsDUCEmptyURL     = errors.New("invalid repository url")
	ErrorsDUCRepoExist    = errors.New("repository is already exist")
	ErrorsDUCRepoNotExist = errors.New("repository is not exist")
)

const (
	topDir     = ".dockyard"
	configName = "config.json"
	cacheDir   = "cache"
)

type DyUpdaterClientConfig struct {
	DefaultServer string
	CacheDir      string
	Repos         []string
}

func (dyc *DyUpdaterClientConfig) exist() bool {
	configFile := filepath.Join(os.Getenv("HOME"), topDir, configName)
	return dyutils.IsFileExist(configFile)
}

func (dyc *DyUpdaterClientConfig) Init() error {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return errors.New("Cannot get home directory")
	}

	topURL := filepath.Join(homeDir, topDir)
	cacheURL := filepath.Join(topURL, cacheDir)
	if !dyutils.IsDirExist(cacheURL) {
		if err := os.MkdirAll(cacheURL, os.ModePerm); err != nil {
			return err
		}
	}

	dyc.CacheDir = cacheURL

	if !dyc.exist() {
		return dyc.save()
	}
	return nil
}

func (dyc *DyUpdaterClientConfig) save() error {
	data, err := json.MarshalIndent(dyc, "", "\t")
	if err != nil {
		return err
	}

	configFile := filepath.Join(os.Getenv("HOME"), topDir, configName)
	if err := ioutil.WriteFile(configFile, data, 0666); err != nil {
		return err
	}

	return nil
}

func (dyc *DyUpdaterClientConfig) Load() error {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return errors.New("Cannot get home directory")
	}

	content, err := ioutil.ReadFile(filepath.Join(homeDir, topDir, configName))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, &dyc); err != nil {
		return err
	}

	if dyc.CacheDir == "" {
		dyc.CacheDir = filepath.Join(homeDir, topDir, cacheDir)
	}

	return nil
}

func (dyc *DyUpdaterClientConfig) Add(url string) error {
	if url == "" {
		return ErrorsDUCEmptyURL
	}

	var err error
	if !dyc.exist() {
		err = dyc.Init()
	} else {
		err = dyc.Load()
	}
	if err != nil {
		return err
	}

	for _, repo := range dyc.Repos {
		if repo == url {
			return ErrorsDUCRepoExist
		}
	}
	dyc.Repos = append(dyc.Repos, url)

	return dyc.save()
}

func (dyc *DyUpdaterClientConfig) Remove(url string) error {
	if url == "" {
		return ErrorsDUCEmptyURL
	}

	if !dyc.exist() {
		return ErrorsDUCRepoNotExist
	}

	if err := dyc.Load(); err != nil {
		return err
	}
	found := false
	for i, _ := range dyc.Repos {
		if dyc.Repos[i] == url {
			found = true
			dyc.Repos = append(dyc.Repos[:i], dyc.Repos[i+1:]...)
			break
		}
	}
	if !found {
		return ErrorsDUCRepoNotExist
	}

	return dyc.save()
}
