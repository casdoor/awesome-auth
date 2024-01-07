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

package controllers

import (
	"net/http"
	"strconv"

	"github.com/beego/beego"
	"github.com/casbin/casvisor/object"
	"github.com/casbin/casvisor/util"
	"github.com/casbin/casvisor/util/guacamole"
	"github.com/gorilla/websocket"
)

const (
	TunnelClosed             int = -1
	Normal                   int = 0
	NotFoundSession          int = 800
	NewTunnelError           int = 801
	ForcedDisconnect         int = 802
	AccessGatewayUnAvailable int = 803
	AccessGatewayCreateError int = 804
	AssetNotActive           int = 805
	NewSshClientError        int = 806
)

var UpGrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"guacamole"},
}

func (c *ApiController) GetAssetTunnel() {
	ctx := c.Ctx
	ws, err := UpGrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
	if err != nil {
		c.ResponseError("WebSocket upgrade failed:", err)
		return
	}

	owner := c.Input().Get("owner")
	name := c.Input().Get("name")
	width := c.Input().Get("width")
	height := c.Input().Get("height")
	dpi := c.Input().Get("dpi")

	intWidth, err := strconv.Atoi(width)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	intHeight, err := strconv.Atoi(height)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	remoteAppName := c.Input().Get("remoteApp")
	remoteAppDir := c.Input().Get("remoteAppDir")
	remoteAppArgs := c.Input().Get("remoteAppArgs")

	asset, err := object.GetAsset(util.GetIdFromOwnerAndName(owner, name))
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	configuration := guacamole.NewConfiguration()
	configuration.Protocol = asset.Protocol
	propertyMap := configuration.LoadConfig()

	setConfig(propertyMap, configuration)

	configuration.SetParameter("width", width)
	configuration.SetParameter("height", height)
	configuration.SetParameter("dpi", dpi)

	configuration.SetParameter("hostname", asset.Ip)
	configuration.SetParameter("port", strconv.Itoa(asset.Port))
	configuration.SetParameter("username", asset.Username)
	configuration.SetParameter("password", asset.Password)

	if asset.Protocol == "rdp" && asset.EnableRemoteApp {
		configuration.SetParameter("remote-app", "||"+remoteAppName)
		configuration.SetParameter("remote-app-dir", remoteAppDir)
		configuration.SetParameter("remote-app-args", remoteAppArgs)
	}

	addr := beego.AppConfig.String("guacamoleEndpoint")
	tunnel, err := guacamole.NewTunnel(addr, configuration)
	if err != nil {
		guacamole.Disconnect(ws, NewTunnelError, err.Error())
		panic(err)
	}

	session := object.Session{
		ConnectionId: tunnel.ConnectionID,
		Width:        intWidth,
		Height:       intHeight,
		Status:       object.Connecting,
		Recording:    configuration.GetParameter(guacamole.RecordingPath),
	}
	if session.Recording == "" {
		// No audit is required when no screen is recorded
		session.Reviewed = true
	}

	_, err = object.AddSession(&session)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	guacamoleHandler := NewGuacamoleHandler(ws, tunnel)
	guacamoleHandler.Start()
	defer guacamoleHandler.Stop()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			_ = tunnel.Close()
			return
		}
		_, err = tunnel.WriteAndFlush(message)
		if err != nil {
			return
		}
	}
}

func setConfig(propertyMap map[string]string, configuration *guacamole.Configuration) {
	switch configuration.Protocol {
	case "rdp":
		configuration.SetParameter("security", "any")
		configuration.SetParameter("ignore-cert", "true")
		configuration.SetParameter("create-drive-path", "true")
		//configuration.SetParameter("resize-method", "reconnect")
		configuration.SetParameter("resize-method", "display-update")
		configuration.SetParameter(guacamole.EnableWallpaper, propertyMap[guacamole.EnableWallpaper])
		configuration.SetParameter(guacamole.EnableTheming, propertyMap[guacamole.EnableTheming])
		configuration.SetParameter(guacamole.EnableFontSmoothing, propertyMap[guacamole.EnableFontSmoothing])
		configuration.SetParameter(guacamole.EnableFullWindowDrag, propertyMap[guacamole.EnableFullWindowDrag])
		configuration.SetParameter(guacamole.EnableDesktopComposition, propertyMap[guacamole.EnableDesktopComposition])
		configuration.SetParameter(guacamole.EnableMenuAnimations, propertyMap[guacamole.EnableMenuAnimations])
		configuration.SetParameter(guacamole.DisableBitmapCaching, propertyMap[guacamole.DisableBitmapCaching])
		configuration.SetParameter(guacamole.DisableOffscreenCaching, propertyMap[guacamole.DisableOffscreenCaching])
		configuration.SetParameter(guacamole.ColorDepth, propertyMap[guacamole.ColorDepth])
		configuration.SetParameter(guacamole.ForceLossless, propertyMap[guacamole.ForceLossless])
		configuration.SetParameter(guacamole.PreConnectionId, propertyMap[guacamole.PreConnectionId])
		configuration.SetParameter(guacamole.PreConnectionBlob, propertyMap[guacamole.PreConnectionBlob])
	case "ssh":
		configuration.SetParameter(guacamole.FontSize, propertyMap[guacamole.FontSize])
		//configuration.SetParameter(guacamole.FontName, propertyMap[guacamole.FontName])
		configuration.SetParameter(guacamole.ColorScheme, propertyMap[guacamole.ColorScheme])
		configuration.SetParameter(guacamole.Backspace, propertyMap[guacamole.Backspace])
		configuration.SetParameter(guacamole.TerminalType, propertyMap[guacamole.TerminalType])
	case "telnet":
		configuration.SetParameter(guacamole.FontSize, propertyMap[guacamole.FontSize])
		//configuration.SetParameter(guacamole.FontName, propertyMap[guacamole.FontName])
		configuration.SetParameter(guacamole.ColorScheme, propertyMap[guacamole.ColorScheme])
		configuration.SetParameter(guacamole.Backspace, propertyMap[guacamole.Backspace])
		configuration.SetParameter(guacamole.TerminalType, propertyMap[guacamole.TerminalType])
	default:
	}
}
