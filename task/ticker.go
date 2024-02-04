// Copyright 2024 The Casbin Authors. All Rights Reserved.
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

package task

import (
	"time"

	"github.com/beego/beego/logs"
	"github.com/casbin/casvisor/object"
)

type Ticker struct{}

func NewTicker() *Ticker {
	return &Ticker{}
}

func (t *Ticker) SetupTicker() {
	// delete unused session every hour
	unUsedSessionTicker := time.NewTicker(time.Hour)
	go func() {
		for range unUsedSessionTicker.C {
			t.deleteUnUsedSession()
		}
	}()
}

func (t *Ticker) deleteUnUsedSession() {
	sessions, err := object.GetSessionsByStatus([]string{object.NoConnect, object.Connecting})
	if err != nil {
		logs.Info("failed to get unused sessions: ", err)
		return
	}

	now := time.Now()
	for _, session := range sessions {
		if session.StartTime != "" {
			startTime, err := time.ParseInLocation(time.RFC3339, session.StartTime, time.Local)
			if err != nil {
				continue
			}
			if now.Sub(startTime).Hours() > 1 {
				_, err := object.DeleteSessionById(session.GetId())
				if err != nil {
					logs.Info("delete session failed: ", err)
					return
				}
			}
		} else {
			_, err := object.DeleteSessionById(session.GetId())
			if err != nil {
				logs.Info("delete session failed: ", err)
				return
			}
		}
	}
}
