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

package authz

import (
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casvisor/conf"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	xormadapter "github.com/casdoor/xorm-adapter/v3"
	stringadapter "github.com/qiangmzsx/string-adapter/v2"
)

var Enforcer *casbin.Enforcer

func InitAuthz() {
	var err error

	tableNamePrefix := conf.GetConfigString("tableNamePrefix")
	driverName := conf.GetConfigString("driverName")
	dataSourceName := conf.GetConfigDataSourceName()
	if conf.GetConfigString("driverName") == "mysql" {
		dataSourceName = dataSourceName + conf.GetConfigString("dbName")
	}

	a, err := xormadapter.NewAdapterWithTableName(driverName, dataSourceName, "api_rule", tableNamePrefix, true)
	if err != nil {
		panic(err)
	}

	modelText := `
[request_definition]
r = subOwner, subName, method, urlPath, objOwner, objName

[policy_definition]
p = subOwner, subName, method, urlPath, objOwner, objName

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (r.subOwner == p.subOwner || p.subOwner == "*") && (r.subName == p.subName || p.subName == "*") && \
  (r.method == p.method || p.method == "*") && (r.urlPath == p.urlPath || p.urlPath == "*") && \
  (r.objOwner == p.objOwner || p.objOwner == "*") && (r.objName == p.objName || p.objName == "*")
`

	m, err := model.NewModelFromString(modelText)
	if err != nil {
		panic(err)
	}

	Enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		panic(err)
	}

	if true {
		ruleText := `
p, built-in, *, *, *, *, *
p, app, *, *, *, *, *
p, *, *, POST, /api/signin, *, *
p, *, *, POST, /api/signout, *, *
p, *, *, GET, /api/get-account, *, *
p, *, *, GET, /api/get-asset-tunnel, *, *
p, *, *, POST, /api/add-asset-tunnel, *, *
p, *, *, POST, /api/start-session, *, *
`

		sa := stringadapter.NewAdapter(ruleText)
		// load all rules from string adapter to enforcer's memory
		err := sa.LoadPolicy(Enforcer.GetModel())
		if err != nil {
			panic(err)
		}

		// save all rules from enforcer's memory to Xorm adapter (DB)
		// same as:
		// a.SavePolicy(Enforcer.GetModel())
		err = Enforcer.SavePolicy()
		if err != nil {
			panic(err)
		}
	}
}

func IsAllowed(user *casdoorsdk.User, subOwner string, subName string, method string, urlPath string, objOwner string, objName string) bool {
	if conf.GetConfigBool("IsDemoMode") {
		if !isAllowedInDemoMode(method, urlPath) {
			return false
		}
	}

	if subOwner == "app" {
		return true
	}

	if user != nil {
		if user.IsDeleted {
			return false
		}

		if user.IsAdmin && (subOwner == objOwner || (objOwner == "admin")) {
			return true
		}
	}

	res, err := Enforcer.Enforce(subOwner, subName, method, urlPath, objOwner, objName)
	if err != nil {
		panic(err)
	}

	return res
}

func isAllowedInDemoMode(method string, urlPath string) bool {
	if method == "POST" {
		if strings.HasPrefix(urlPath, "/api/signin") || urlPath == "/api/signout" || urlPath == "/api/add-asset-tunnel" || urlPath == "/api/start-session" || urlPath == "/api/stop-session" {
			return true
		} else {
			return false
		}
	}

	// If method equals GET
	return true
}
