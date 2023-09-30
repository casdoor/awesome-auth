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
	"github.com/casbin/casvisor/util/tunnel"
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

func (c *ApiController) GetAssetTunnel() error {
	ctx := c.Ctx
	ws, err := UpGrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
	if err != nil {
		beego.Error("WebSocket upgrade failed:", err)
		return err
	}

	owner := c.Input().Get("owner")
	name := c.Input().Get("name")
	width := c.Input().Get("width")
	height := c.Input().Get("height")
	dpi := c.Input().Get("dpi")

	remoteAppName := c.Input().Get("remoteApp")
	remoteAppDir := c.Input().Get("remoteAppDir")
	remoteAppArgs := c.Input().Get("remoteAppArgs")

	asset, err := object.GetAsset(util.GetIdFromOwnerAndName(owner, name))
	if err != nil {
		return err
	}

	configuration := tunnel.NewConfiguration()
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

	// Todo: Support ssh via privateKey
	//if asset.Protocol == "ssh"  {
	//	if len(asset.PrivateKey) > 0 && asset.PrivateKey != "-" {
	//		configuration.SetParameter("username", asset.Username)
	//		configuration.SetParameter("private-key", asset.PrivateKey)
	//		configuration.SetParameter("passphrase", asset.Passphrase)
	//	} else {
	//		configuration.SetParameter("username", asset.Username)
	//		configuration.SetParameter("password", asset.Password)
	//	}
	//}

	addr := beego.AppConfig.String("guacamoleEndpoint")
	// fmt.Sprintf("%s:%s", configuration.GetParameter("hostname"), configuration.GetParameter("port"))
	// log.Debug("Intializing guacd.go session", log.String("sessionId", sessionId), log.String("addr", addr), log.String("asset", asset))

	guacdTunnel, err := tunnel.NewTunnel(addr, configuration)
	if err != nil {
		tunnel.Disconnect(ws, NewTunnelError, err.Error())
		// log.Error("Failed to start session", log.String("sessionId", sessionId), log.NamedError("err", err))
		panic(err)
	}

	guacamoleHandler := NewGuacamoleHandler(ws, guacdTunnel)
	guacamoleHandler.Start()
	defer guacamoleHandler.Stop()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			// log.Debug("WebSocket shutdown", log.String("sessionId", sessionId), log.NamedError("err", err))
			_ = guacdTunnel.Close()

			return nil
		}
		_, err = guacdTunnel.WriteAndFlush(message)
		if err != nil {
			//service.SessionService.CloseSessionById(sessionId, TunnelClosed, "Remote connection shut down")
			//panic(err)
			return nil
		}
	}
}

func setConfig(propertyMap map[string]string, configuration *tunnel.Configuration) {
	switch configuration.Protocol {
	case "rdp":
		configuration.SetParameter("security", "any")
		configuration.SetParameter("ignore-cert", "true")
		configuration.SetParameter("create-drive-path", "true")
		//configuration.SetParameter("resize-method", "reconnect")
		configuration.SetParameter("resize-method", "display-update")
		configuration.SetParameter(tunnel.EnableWallpaper, propertyMap[tunnel.EnableWallpaper])
		configuration.SetParameter(tunnel.EnableTheming, propertyMap[tunnel.EnableTheming])
		configuration.SetParameter(tunnel.EnableFontSmoothing, propertyMap[tunnel.EnableFontSmoothing])
		configuration.SetParameter(tunnel.EnableFullWindowDrag, propertyMap[tunnel.EnableFullWindowDrag])
		configuration.SetParameter(tunnel.EnableDesktopComposition, propertyMap[tunnel.EnableDesktopComposition])
		configuration.SetParameter(tunnel.EnableMenuAnimations, propertyMap[tunnel.EnableMenuAnimations])
		configuration.SetParameter(tunnel.DisableBitmapCaching, propertyMap[tunnel.DisableBitmapCaching])
		configuration.SetParameter(tunnel.DisableOffscreenCaching, propertyMap[tunnel.DisableOffscreenCaching])
		configuration.SetParameter(tunnel.ColorDepth, propertyMap[tunnel.ColorDepth])
		configuration.SetParameter(tunnel.ForceLossless, propertyMap[tunnel.ForceLossless])
		configuration.SetParameter(tunnel.PreConnectionId, propertyMap[tunnel.PreConnectionId])
		configuration.SetParameter(tunnel.PreConnectionBlob, propertyMap[tunnel.PreConnectionBlob])
	case "ssh":
		configuration.SetParameter(tunnel.FontSize, propertyMap[tunnel.FontSize])
		//configuration.SetParameter(tunnel.FontName, propertyMap[tunnel.FontName])
		configuration.SetParameter(tunnel.ColorScheme, propertyMap[tunnel.ColorScheme])
		configuration.SetParameter(tunnel.Backspace, propertyMap[tunnel.Backspace])
		configuration.SetParameter(tunnel.TerminalType, propertyMap[tunnel.TerminalType])
	case "telnet":
		configuration.SetParameter(tunnel.FontSize, propertyMap[tunnel.FontSize])
		//configuration.SetParameter(tunnel.FontName, propertyMap[tunnel.FontName])
		configuration.SetParameter(tunnel.ColorScheme, propertyMap[tunnel.ColorScheme])
		configuration.SetParameter(tunnel.Backspace, propertyMap[tunnel.Backspace])
		configuration.SetParameter(tunnel.TerminalType, propertyMap[tunnel.TerminalType])
	default:
	}
}
