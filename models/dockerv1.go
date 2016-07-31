/*
Copyright 2015 The ContainerOps Authors All rights reserved.

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

package models

import (
	"time"
)

//DockerV1 is Docker Repository V1 repository.
type DockerV1 struct {
	Id          int64      `json:"id" gorm:"primary_key"`
	Namespace   string     `json:"namespace" sql:"not null;type:varchar(255)" gorm:"unique_index:v1_repository"`
	Repository  string     `json:"repository" sql:"not null;type:varchar(255)" gorm:"unique_index:v1_repository"`
	JSON        string     `json:"json" sql:"null;type:text"`
	Manifests   string     `json:"manifests" sql:"null;type:text"`
	Agent       string     `json:"agent" sql:"null;type:text"`
	Description string     `json:"description" sql:"null;type:text"`
	Size        int64      `json:"size" sql:"default:0"`
	Locked      bool       `json:"locked" sql:"default:false"` //When create/update the repository, the locked will be true.
	CreatedAt   time.Time  `json:"created" sql:""`
	UpdatedAt   time.Time  `json:"updated" sql:""`
	DeletedAt   *time.Time `json:"deleted" sql:"index"`
}

//TableName in mysql is "docker_v1".
func (r *DockerV1) TableName() string {
	return "docker_v1"
}

//
type DockerImageV1 struct {
	Id        int64      `json:"id" gorm:"primary_key"`
	ImageId   string     `json:"imageid" sql:"not null;unique;varchar(255)"`
	JSON      string     `json:"json" sql:"null;type:text"`
	Ancestry  string     `json:"ancestry" sql:"null;type:text"`
	Checksum  string     `json:"checksum" sql:"null;unique;type:varchar(255)"`
	Payload   string     `json:"payload" sql:"null;type:varchar(255)"`
	Path      string     `json:"path" sql:"null;type:text"`
	OSS       string     `json:"oss" sql:"null;type:text"`
	Size      int64      `json:"size" sql:"default:0"`
	Locked    bool       `json:"locked" sql:"default:false"`
	CreatedAt time.Time  `json:"created" sql:""`
	UpdatedAt time.Time  `json:"updated" sql:""`
	DeletedAt *time.Time `json:"deleted" sql:"index"`
}

//TableName in mysql is "docker_image_v1".
func (*DockerImageV1) TableName() string {
	return "docker_image_v1"
}

//
type DockerTagV1 struct {
	Id        int64      `json:"id" gorm:"primary_key"`
	DockerV1  int64      `json:"dockerv1" sql:"not null"`
	Tag       string     `json:"tag" sql:"not null;varchar(255)"`
	ImageId   string     `json:"imageid" sql:"not null;varchar(255)"`
	CreatedAt time.Time  `json:"created" sql:""`
	UpdatedAt time.Time  `json:"updated" sql:""`
	DeletedAt *time.Time `json:"deleted" sql:"index"`
}

//TableName in mysql is "docker_tag_v1".
func (*DockerTagV1) TableName() string {
	return "docker_tag_v1"
}

//Save function save all properties of Docker Registry V1 repository.
func (r *DockerV1) Save(namespace, repository string) error {
	return nil
}

//Put function will create or update repository.
func (r *DockerV1) Put(namespace, repository, json, agent string) error {
	return nil
}
