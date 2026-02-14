// Copyright 2023 The Hanzo Authors. All Rights Reserved.
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
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/beego/beego/context"
	"github.com/hanzoai/vm/conf"
	"github.com/hanzoai/vm/util"
)

var (
	oldIamEndpoint     = "https://door.casdoor.com"
	newIamEndpoint     = conf.GetConfigString("iamEndpoint")
	oldClientId        = "b108dacba027db36ec26"
	newClientId        = conf.GetConfigString("clientId")
	oldIamOrganization = "casbin"
	newIamOrganization = conf.GetConfigString("iamOrganization")
	oldIamApplication  = "app-hanzo-vm"
	newIamApplication  = conf.GetConfigString("iamApplication")
)

func TransparentStatic(ctx *context.Context) {
	urlPath := ctx.Request.URL.Path
	if strings.HasPrefix(urlPath, "/api/") {
		return
	}

	path := "web/build"
	if urlPath == "/" {
		path += "/index.html"
	} else {
		path += urlPath
	}

	if !util.FileExist(path) {
		path = "web/build/index.html"
	}

	serveFileWithReplace(ctx.ResponseWriter, ctx.Request, path)
}

func serveFileWithReplace(w http.ResponseWriter, r *http.Request, name string) {
	f, err := os.Open(filepath.Clean(name))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		panic(err)
	}

	content := util.ReadStringFromPath(name)
	if oldIamEndpoint != newIamEndpoint {
		content = strings.ReplaceAll(content, fmt.Sprintf("\"%s\"", oldIamEndpoint), fmt.Sprintf("\"%s\"", newIamEndpoint))
	}
	if oldClientId != newClientId {
		content = strings.ReplaceAll(content, fmt.Sprintf("\"%s\"", oldClientId), fmt.Sprintf("\"%s\"", newClientId))
	}
	if oldIamOrganization != newIamOrganization {
		content = strings.ReplaceAll(content, fmt.Sprintf("\"%s\"", oldIamOrganization), fmt.Sprintf("\"%s\"", newIamOrganization))
	}
	if oldIamApplication != newIamApplication {
		content = strings.ReplaceAll(content, fmt.Sprintf("\"%s\"", oldIamApplication), fmt.Sprintf("\"%s\"", newIamApplication))
	}

	http.ServeContent(w, r, d.Name(), d.ModTime(), strings.NewReader(content))
}
