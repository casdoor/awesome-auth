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
	"errors"
	"strconv"
	"sync"

	"github.com/casbin/casvisor/util"
	"github.com/casbin/casvisor/util/guacamole"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
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

func GetSessionCount(owner, status, field, value string) (int64, error) {
	session := GetSession(owner, -1, -1, field, value, "", "")
	return session.Count(&Session{Status: status})
}

func GetSessions(owner string) ([]*Session, error) {
	sessions := []*Session{}
	err := adapter.engine.Desc("connected_time").Find(&sessions, &Session{Owner: owner})
	if err != nil {
		return sessions, err
	}

	return sessions, nil
}

func GetPaginationSessions(owner, status string, offset, limit int, field, value, sortField, sortOrder string) ([]*Session, error) {
	var sessions []*Session
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&sessions, &Session{Status: status})
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

func UpdateSession(id string, session *Session, columns ...string) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	if s, err := getSession(owner, name); err != nil {
		return false, err
	} else if s == nil {
		return false, nil
	}

	_, err := adapter.engine.ID(core.PK{owner, name}).Cols(columns...).Update(session)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteSession(session *Session) (bool, error) {
	affected, err := adapter.engine.ID(core.PK{session.Owner, session.Name}).Delete(&Session{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func AddSession(session *Session) (bool, error) {
	affected, err := adapter.engine.Insert(session)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func CreateSession(ip, assetId, mode string, user *casdoorsdk.User) (*Session, error) {
	asset, err := GetAsset(assetId)
	if err != nil {
		return nil, err
	}

	if asset == nil {
		return nil, nil
	}

	if !user.IsAdmin && user.Owner != "built-in" {
		return nil, errors.New("you have no permission to access this asset")
	}

	session := Session{
		Owner:       asset.Owner,
		Name:        util.GenerateId(),
		CreatedTime: util.GetCurrentTime(),
		Protocol:    asset.Protocol,
		IP:          asset.Ip,
		Port:        asset.Port,
		AssetId:     assetId,
		Username:    asset.Username,
		Password:    asset.Password,
		Creator:     user.Name,
		ClientIP:    ip,
		Width:       1280,
		Height:      720,
		Status:      NoConnect,
		Mode:        mode,
		FileSystem:  true,
		Upload:      true,
		Download:    true,
		Delete:      true,
		Rename:      true,
		Edit:        true,
		CreateDir:   true,
		Copy:        true,
		Paste:       true,
		Reviewed:    false,
	}

	_, err = AddSession(&session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func CloseDBSession(id string, code int, msg string) error {
	s, err := GetConnSession(id)
	if err != nil {
		return err
	}
	if s == nil {
		return nil
	}

	if s.Status == Disconnected {
		return nil
	}

	if s.Status == Connecting {
		// The session has not been established successfully, so you do not need to save data
		_, err := DeleteSession(s)
		if err != nil {
			return err
		}
		return nil
	}

	s.Status = Disconnected
	s.Code = code
	s.Message = msg
	s.DisconnectedTime = util.GetCurrentTime()
	s.Password = "-"
	s.PrivateKey = "-"
	s.Passphrase = "-"

	_, err = UpdateSession(id, s, "status", "code", "message", "disconnected_time", "password", "private_key", "passphrase")
	if err != nil {
		return err
	}
	return nil
}

func WriteCloseMessage(sess *guacamole.Session, mode string, code int, msg string) {
	err := guacamole.NewInstruction("error", "", strconv.Itoa(code))
	_ = sess.WriteString(err.String())
	disconnect := guacamole.NewInstruction("disconnect")
	_ = sess.WriteString(disconnect.String())
}

var mutex sync.Mutex

func CloseSession(id string, code int, msg string) error {
	mutex.Lock()
	defer mutex.Unlock()
	guacSession := guacamole.GlobalSessionManager.Get(id)

	if guacSession != nil {
		WriteCloseMessage(guacSession, guacSession.Mode, code, msg)

		if guacSession.Observer != nil {
			guacSession.Observer.Range(func(key string, ob *guacamole.Session) {
				WriteCloseMessage(ob, ob.Mode, code, msg)
			})
		}
	}
	guacamole.GlobalSessionManager.Delete(id)

	err := CloseDBSession(id, code, msg)
	if err != nil {
		return err
	}

	return nil
}
