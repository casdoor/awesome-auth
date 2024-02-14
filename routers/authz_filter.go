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

package routers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/beego/beego/context"
	"github.com/casbin/casvisor/authz"
	"github.com/casbin/casvisor/util"
)

type Object struct {
	Owner        string `json:"owner"`
	Name         string `json:"name"`
	AccessKey    string `json:"accessKey"`
	AccessSecret string `json:"accessSecret"`
}

func getUsername(ctx *context.Context) (username string) {
	user := GetSessionUser(ctx)
	if user != nil {
		username = util.GetIdFromOwnerAndName(user.Owner, user.Name)
	} else {
		username, _ = getUsernameByClientIdSecret(ctx)
	}
	return
}

func getSubject(ctx *context.Context) (string, string) {
	username := getUsername(ctx)
	if username == "" {
		return "anonymous", "anonymous"
	}

	// username == "built-in/admin"
	return util.GetOwnerAndNameFromId(username)
}

func getObject(ctx *context.Context) (string, string) {
	method := ctx.Request.Method

	if method == http.MethodGet {
		// query == "?id=built-in/admin"
		id := ctx.Input.Query("id")
		if id != "" {
			return util.GetOwnerAndNameFromIdNoCheck(id)
		}

		owner := ctx.Input.Query("owner")
		if owner != "" {
			return owner, ""
		}

		return "", ""
	} else {
		id := ctx.Input.Query("id")
		if id != "" {
			return util.GetOwnerAndNameFromIdNoCheck(id)
		}

		body := ctx.Input.RequestBody
		if len(body) == 0 {
			id := ctx.Request.Form.Get("id")
			if id != "" {
				return util.GetOwnerAndNameFromIdNoCheck(id)
			}

			return ctx.Request.Form.Get("owner"), ctx.Request.Form.Get("name")
		}

		var obj Object
		err := json.Unmarshal(body, &obj)
		if err != nil {
			return "", ""
		}

		return obj.Owner, obj.Name
	}
}

func willLog(subOwner string, subName string, method string, urlPath string, objOwner string, objName string) bool {
	if subOwner == "anonymous" && subName == "anonymous" && method == "GET" && (urlPath == "/api/get-account") && objOwner == "" && objName == "" {
		return false
	}
	return true
}

func getUrlPath(urlPath string) string {
	return urlPath
}

func ApiFilter(ctx *context.Context) {
	subOwner, subName := getSubject(ctx)
	method := ctx.Request.Method
	urlPath := getUrlPath(ctx.Request.URL.Path)
	objOwner, objName := getObject(ctx)

	user := GetSessionUser(ctx)
	isAllowed := authz.IsAllowed(user, subOwner, subName, method, urlPath, objOwner, objName)

	result := "deny"
	if isAllowed {
		result = "allow"
	}

	if willLog(subOwner, subName, method, urlPath, objOwner, objName) {
		logLine := fmt.Sprintf("subOwner = %s, subName = %s, method = %s, urlPath = %s, obj.Owner = %s, obj.Name = %s, result = %s",
			subOwner, subName, method, urlPath, objOwner, objName, result)
		fmt.Println(logLine)
		util.LogInfo(ctx, logLine)
	}

	if !isAllowed {
		requestDeny(ctx)
	}
}
