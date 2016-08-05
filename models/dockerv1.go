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
	CreatedAt   time.Time  `json:"create_at" sql:""`
	UpdatedAt   time.Time  `json:"update_at" sql:""`
	DeletedAt   *time.Time `json:"delete_at" sql:"index"`
}

//TableName in mysql is "docker_v1".
func (r *DockerV1) TableName() string {
	return "docker_v1"
}

//
type DockerImageV1 struct {
	Id        int64      `json:"id" gorm:"primary_key"`
	ImageId   string     `json:"image_id" sql:"not null;unique;varchar(255)"`
	JSON      string     `json:"json" sql:"null;type:text"`
	Ancestry  string     `json:"ancestry" sql:"null;type:text"`
	Checksum  string     `json:"checksum" sql:"null;type:varchar(255)"`
	Payload   string     `json:"payload" sql:"null;type:varchar(255)"`
	Path      string     `json:"path" sql:"null;type:text"`
	OSS       string     `json:"oss" sql:"null;type:text"`
	Size      int64      `json:"size" sql:"default:0"`
	Locked    bool       `json:"locked" sql:"default:false"`
	CreatedAt time.Time  `json:"create_at" sql:""`
	UpdatedAt time.Time  `json:"update_at" sql:""`
	DeletedAt *time.Time `json:"delete_at" sql:"index"`
}

//TableName in mysql is "docker_image_v1".
func (i *DockerImageV1) TableName() string {
	return "docker_image_v1"
}

//
type DockerTagV1 struct {
	Id        int64      `json:"id" gorm:"primary_key"`
	DockerV1  int64      `json:"docker_v1" sql:"not null"`
	Tag       string     `json:"tag" sql:"not null;varchar(255)"`
	ImageId   string     `json:"image_id" sql:"not null;varchar(255)"`
	CreatedAt time.Time  `json:"create_at" sql:""`
	UpdatedAt time.Time  `json:"update_at" sql:""`
	DeletedAt *time.Time `json:"delete_at" sql:"index"`
}

//TableName in mysql is "docker_tag_v1".
func (t *DockerTagV1) TableName() string {
	return "docker_tag_v1"
}

//Put function will create or update repository.
func (r *DockerV1) Put(namespace, repository, json, agent string) error {
	r.Namespace, r.Repository, r.JSON, r.Agent, r.Locked = namespace, repository, json, agent, true

	tx := db.Begin()

	if err := db.Debug().Where("namespace = ? AND repository = ? ", namespace, repository).FirstOrCreate(&r).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Debug().Model(&r).Updates(map[string]interface{}{"json": json, "agent": agent, "locked": true}).Error; err != nil {
		tx.Rollback()
		return err
	} else if err == nil {
		tx.Commit()
		return nil
	}

	tx.Commit()
	return nil
}

//Get is search image by ImageID.
func (i *DockerImageV1) Get(imageID string) (DockerImageV1, error) {
	if err := db.Debug().Where("image_id = ?", imageID).First(&i).Error; err != nil {
		return *i, err
	} else {
		return *i, nil
	}
}

//Put image json by ImageID.
func (i *DockerImageV1) PutJSON(imageID, json string) error {
	i.ImageId = imageID

	tx := db.Begin()

	if err := db.Debug().Where("image_id = ?", imageID).FirstOrCreate(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Debug().Model(&i).Updates(map[string]interface{}{"json": json}).Error; err != nil {
		tx.Rollback()
		return err
	} else if err == nil {
		tx.Commit()
		return nil
	}

	tx.Commit()
	return nil
}
