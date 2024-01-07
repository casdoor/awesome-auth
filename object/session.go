// Copyright 2024 The casbin Authors. All Rights Reserved.
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
	"github.com/casbin/casvisor/util"
	"xorm.io/core"
)

const (
	NoConnect    = "no_connect"
	Connecting   = "connecting"
	Connected    = "connected"
	Disconnected = "disconnected"
)

type Session struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	ConnectedTime    string `xorm:"varchar(100)" json:"connectedTime"`
	DisconnectedTime string `xorm:"varchar(100)" json:"disconnectedTime"`

	Protocol     string `xorm:"varchar(20)" json:"protocol"`
	IP           string `xorm:"varchar(200)" json:"ip"`
	Port         int    `json:"port"`
	ConnectionId string `xorm:"varchar(50)" json:"connectionId"`
	AssetId      string `xorm:"varchar(200) index" json:"assetId"`
	Username     string `xorm:"varchar(200)" json:"username"`
	Password     string `xorm:"varchar(500)" json:"password"`
	Creator      string `xorm:"varchar(36) index" json:"creator"`
	ClientIP     string `xorm:"varchar(200)" json:"clientIp"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Status       string `xorm:"varchar(20) index" json:"status"`
	Recording    string `xorm:"varchar(1000)" json:"recording"`
	PrivateKey   string `xorm:"mediumtext" json:"privateKey"`
	Passphrase   string `xorm:"varchar(500)" json:"passphrase"`
	Code         int    `json:"code"`
	Message      string `json:"message"`

	Mode       string `xorm:"varchar(10)" json:"mode"`
	FileSystem bool   `json:"fileSystem"`
	Upload     bool   `json:"upload"`
	Download   bool   `json:"download"`
	Delete     bool   `json:"delete"`
	Rename     bool   `json:"rename"`
	Edit       bool   `json:"edit"`
	CreateDir  bool   `json:"createDir"`
	Copy       bool   `json:"copy"`
	Paste      bool   `json:"paste"`

	Reviewed     bool  `json:"reviewed"`
	CommandCount int64 `json:"commandCount"`
}

func GetSessionCount(owner, field, value string) (int64, error) {
	session := GetSession(owner, -1, -1, field, value, "", "")
	return session.Count(&Session{})
}

func GetSessions(owner string) ([]*Session, error) {
	sessions := []*Session{}
	err := adapter.engine.Desc("connected_time").Find(&sessions, &Session{Owner: owner})
	if err != nil {
		return sessions, err
	}

	return sessions, nil
}

func GetPaginationSessions(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*Session, error) {
	var sessions []*Session
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&sessions)
	if err != nil {
		return sessions, err
	}

	return sessions, nil
}

func getSession(owner string, name string) (*Session, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	session := Session{Owner: owner, Name: name}
	existed, err := adapter.engine.Get(&session)
	if err != nil {
		return &session, err
	}

	if existed {
		return &session, nil
	} else {
		return nil, nil
	}
}

func GetConnSession(id string) (*Session, error) {
	owner, name := util.GetOwnerAndNameFromIdNoCheck(id)
	return getSession(owner, name)
}

func UpdateSession(id string, session *Session) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	if s, err := getSession(owner, name); err != nil {
		return false, err
	} else if s == nil {
		return false, nil
	}

	_, err := adapter.engine.ID(core.PK{owner, name}).AllCols().Update(session)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteSession(id string) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	if s, err := getSession(owner, name); err != nil {
		return false, err
	} else if s == nil {
		return false, nil
	}

	_, err := adapter.engine.ID(core.PK{owner, name}).Delete(&Session{})
	if err != nil {
		return false, err
	}

	return true, nil
}

func AddSession(session *Session) (bool, error) {
	affected, err := adapter.engine.Insert(session)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}
