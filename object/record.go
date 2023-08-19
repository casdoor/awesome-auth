// Copyright 2023 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"strings"

	"github.com/beego/beego/context"
	"github.com/casbin/casvisor/util"
)

type Record struct {
	Id int `xorm:"int notnull pk autoincr" json:"id"`

	Owner       string `xorm:"varchar(100) index" json:"owner"`
	Name        string `xorm:"varchar(100) index" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	Organization string `xorm:"varchar(100)" json:"organization"`
	ClientIp     string `xorm:"varchar(100)" json:"clientIp"`
	User         string `xorm:"varchar(100)" json:"user"`
	Method       string `xorm:"varchar(100)" json:"method"`
	RequestUri   string `xorm:"varchar(1000)" json:"requestUri"`
	Action       string `xorm:"varchar(1000)" json:"action"`

	Object string `xorm:"-" json:"object"`
	// ExtendedUser *User  `xorm:"-" json:"extendedUser"`

	IsTriggered bool `json:"isTriggered"`
}

func GetRecordCount(field, value string, filterRecord *Record) (int64, error) {
	session := GetSession("", -1, -1, field, value, "", "")
	return session.Count(filterRecord)
}

func GetRecords() ([]*Record, error) {
	records := []*Record{}
	err := adapter.engine.Desc("id").Find(&records)
	if err != nil {
		return records, err
	}

	return records, nil
}

func GetPaginationRecords(offset, limit int, field, value, sortField, sortOrder string, filterRecord *Record) ([]*Record, error) {
	records := []*Record{}
	session := GetSession("", offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&records, filterRecord)
	if err != nil {
		return records, err
	}

	return records, nil
}

func GetRecordsByField(record *Record) ([]*Record, error) {
	records := []*Record{}
	err := adapter.engine.Find(&records, record)
	if err != nil {
		return records, err
	}

	return records, nil
}

func getRecord(owner string, name string) (*Record, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	record := Record{Name: name}
	existed, err := adapter.engine.Get(&record)
	if err != nil {
		return &record, err
	}

	if existed {
		return &record, nil
	} else {
		return nil, nil
	}
}

func GetRecord(id string) (*Record, error) {
	owner, name := util.GetOwnerAndNameFromIdNoCheck(id)
	return getRecord(owner, name)
}

func UpdateRecord(id string, record *Record) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	if p, err := getRecord(owner, name); err != nil {
		return false, err
	} else if p == nil {
		return false, nil
	}

	affected, err := adapter.engine.Where("name = ?", name).AllCols().Update(record)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func NewRecord(ctx *context.Context) *Record {
	ip := strings.Replace(util.GetIPFromRequest(ctx.Request), ": ", "", -1)
	action := strings.Replace(ctx.Request.URL.Path, "/api/", "", -1)
	requestUri := util.FilterQuery(ctx.Request.RequestURI, []string{"accessToken"})
	if len(requestUri) > 1000 {
		requestUri = requestUri[0:1000]
	}

	object := ""
	if ctx.Input.RequestBody != nil && len(ctx.Input.RequestBody) != 0 {
		object = string(ctx.Input.RequestBody)
	}

	record := Record{
		Name:        util.GenerateId(),
		CreatedTime: util.GetCurrentTime(),
		ClientIp:    ip,
		User:        "",
		Method:      ctx.Request.Method,
		RequestUri:  requestUri,
		Action:      action,
		Object:      object,
		IsTriggered: false,
	}
	return &record
}

func AddRecord(record *Record) bool {
	if record.Organization == "app" {
		return false
	}

	record.Owner = record.Organization

	affected, err := adapter.engine.Insert(record)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func DeleteRecord(record *Record) (bool, error) {
	affected, err := adapter.engine.Where("name = ?", record.Name).Delete(&Record{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}
