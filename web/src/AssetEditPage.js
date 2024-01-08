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

// http://localhost:18001/assets/dccb0b3e-aa59-443f-996b-c69a98b21ea9

import React from "react";
import {Button, Card, Col, Input, Row, Select, Switch} from "antd";
import * as AssetBackend from "./backend/AssetBackend";
import * as Setting from "./Setting";
import i18next from "i18next";
import ServiceTable from "./ServiceTable";
import RemoteAppTable from "./RemoteAppTable";

const {Option} = Select;

class AssetEditPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      assetOwner: props.match.params.organizationName,
      assetName: props.match.params.assetName,
      asset: null,
      organizations: [],
      mode: props.location.mode !== undefined ? props.location.mode : "edit",
    };

    this.timer = null;
  }

  UNSAFE_componentWillMount() {
    this.getAsset();
    this.startTimer();
  }

  componentWillUnmount() {
    this.stopTimer();
  }

  getAsset(updateFromRemote = false) {
    AssetBackend.getAsset(this.props.account.owner, this.state.assetName)
      .then((res) => {
        if (res.status === "ok") {
          // Clone the fetched asset data using the spread operator
          const newAsset = {...res.data};

          if (!updateFromRemote && this.state.asset !== null) {
            newAsset.autoQuery = this.state.asset.autoQuery;

            // Update or add remote Apps to the asset data
            newAsset.remoteApps = newAsset.remoteApps.map((newApp, i) => {
              if (i < this.state.asset.remoteApps.length) {
                const oldApp = this.state.asset.remoteApps[i];
                return {
                  ...oldApp, // Preserve old attributes
                  ...newApp, // Override with new attributes
                };
              } else {
                return newApp; // Add new service
              }
            }).slice(0, this.state.asset.remoteApps.length);

            // Update or add services to the asset data
            newAsset.services = newAsset.services.map((newService, i) => {
              if (i < this.state.asset.services.length) {
                const oldService = this.state.asset.services[i];
                return {
                  ...oldService, // Preserve old attributes
                  ...newService, // Override with new attributes
                };
              } else {
                return newService; // Add new service
              }
            }).slice(0, this.state.asset.services.length);
          }

          this.setState({
            asset: newAsset,
          });
        } else {
          Setting.showMessage("error", `Failed to get asset: ${res.msg}`);
        }
      });
  }

  startTimer() {
    if (this.timer === null) {
      this.timer = window.setInterval(this.doTimer.bind(this), 3000);
    }
  }

  stopTimer() {
    if (this.timer !== null) {
      clearInterval(this.timer);
      this.timer = null;
    }
  }

  doTimer() {
    if (this.state.asset?.autoQuery) {
      this.getAsset(false);
    }
  }

  parseAssetField(key, value) {
    if (["port"].includes(key)) {
      value = Setting.myParseInt(value);
    }
    return value;
  }

  // parseAssetField(key, value) {
  //   if ([""].includes(key)) {
  //     value = Setting.myParseInt(value);
  //   }
  //   return value;
  // }

  updateAssetField(key, value) {
    value = this.parseAssetField(key, value);

    const asset = this.state.asset;
    asset[key] = value;
    this.setState({
      asset: asset,
    });
  }

  getDefaultPort(protocol) {
    if (protocol === "rdp") {
      return 3389;
    } else if (protocol === "vnc") {
      return 5900;
    } else if (protocol === "ssh") {
      return 22;
    } else if (protocol === "telnet") {
      return 23;
    } else {
      return 0;
    }
  }

  renderAsset() {
    return (
      <Card size="small" title={
        <div>
          {this.state.mode === "add" ? i18next.t("asset:New Asset") : i18next.t("asset:Edit Asset")}&nbsp;&nbsp;&nbsp;&nbsp;
          <Button onClick={() => this.submitAssetEdit(false)}>{i18next.t("general:Save")}</Button>
          <Button style={{marginLeft: "20px"}} type="primary" onClick={() => this.submitAssetEdit(true)}>{i18next.t("general:Save & Exit")}</Button>
          {this.state.mode === "add" ? <Button style={{marginLeft: "20px"}} onClick={() => this.deleteAsset()}>{i18next.t("general:Cancel")}</Button> : null}
        </div>
      } style={{marginLeft: "5px"}} type="inner">
        <Row style={{marginTop: "10px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Organization"), i18next.t("general:Organization - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.asset.owner} onChange={e => {
              this.updateAssetField("owner", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Name"), i18next.t("general:Name - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.asset.name} onChange={e => {
              this.updateAssetField("name", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Description"), i18next.t("general:Description - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.asset.description} onChange={e => {
              this.updateAssetField("description", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Protocol"), i18next.t("general:Protocol - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Select virtual={false} style={{width: "100%"}} value={this.state.asset.protocol} onChange={value => {
              this.updateAssetField("protocol", value);
              this.updateAssetField("port", this.getDefaultPort(value));
            }}>
              {
                [
                  {id: "rdp", name: "RDP"},
                  {id: "vnc", name: "VNC"},
                  {id: "ssh", name: "SSH"},
                  {id: "telnet", name: "Telnet"},
                ].map((item, index) => <Option key={index} value={item.id}>{item.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:IP"), i18next.t("general:IP - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.asset.ip} onChange={e => {
              this.updateAssetField("ip", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Port"), i18next.t("general:Port - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input
              value={this.state.asset.port}
              defaultValue={this.getDefaultPort(this.state.asset.protocol)}
              onChange={e => {
                this.updateAssetField("port", e.target.value);
              }}
            />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Username"), i18next.t("general:Username - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.asset.username} onChange={e => {
              this.updateAssetField("username", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Password"), i18next.t("general:Password - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.asset.password} onChange={e => {
              this.updateAssetField("password", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Language"), i18next.t("general:Language - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.asset.language} onChange={e => {
              this.updateAssetField("language", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Auto query"), i18next.t("general:Auto query - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Switch checked={this.state.asset.autoQuery} onChange={checked => {
              this.updateAssetField("autoQuery", checked);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Is Permanent"), i18next.t("application:Is Permanent - Tooltip"))} :
          </Col>
          <Col span={1} >
            <Switch checked={this.state.asset.isPermanent} onChange={checked => {
              this.updateAssetField("isPermanent", checked);
            }} />
          </Col>
        </Row>
        {this.state.asset.protocol === "rdp" && (
          <div>
            <Row style={{marginTop: "20px"}} >
              <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
                {Setting.getLabel(i18next.t("general:Enable Remote App"), i18next.t("general:Enable Remote App - Tooltip"))} :
              </Col>
              <Col span={22}>
                <Switch checked={this.state.asset.enableRemoteApp} onChange={checked => {
                  if (checked && this.state.asset.remoteApps.length === 0) {
                    Setting.showMessage("error", i18next.t("asset:Cannot enable Remote App when Remote Apps are empty. Please add at least one Remote App in below table first, then enable again"));
                    return;
                  }
                  this.updateAssetField("enableRemoteApp", checked);
                }} />
              </Col>
            </Row>
            {this.state.asset.enableRemoteApp && (
              <Row style={{marginTop: "20px"}} >
                <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2} >
                  {Setting.getLabel(i18next.t("general:Remote Apps"), i18next.t("general:Remote Apps - Tooltip"))} :
                </Col>
                <Col span={22} >
                  <RemoteAppTable title={"Remote Apps"} table={this.state.asset.remoteApps} onUpdateTable={(value) => {
                    this.updateAssetField("remoteApps", value);
                  }} />
                </Col>
              </Row>
            )}
          </div>
        )}
        {this.state.asset.protocol === "ssh" && (
          <div>
          </div>
        )}
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Services"), i18next.t("general:Services - Tooltip"))} :
          </Col>
          <Col span={22} >
            <ServiceTable title={"Services"} table={this.state.asset.services} onUpdateTable={(value) => {
              this.updateAssetField("services", value);
            }} />
          </Col>
        </Row>
      </Card>
    );
  }

  submitAssetEdit(willExist) {
    const asset = Setting.deepCopy(this.state.asset);
    AssetBackend.updateAsset(this.state.asset.owner, this.state.assetName, asset)
      .then((res) => {
        if (res.status === "ok") {
          if (res.data) {
            Setting.showMessage("success", "Successfully saved");
            this.setState({
              assetName: this.state.asset.name,
            });
            if (willExist) {
              this.props.history.push("/assets");
            } else {
              this.props.history.push(`/assets/${this.state.asset.owner}/${encodeURIComponent(this.state.asset.name)}`);
            }
            // this.getAsset(true);
          } else {
            Setting.showMessage("error", "failed to save: server side failure");
            this.updateAssetField("name", this.state.assetName);
          }
        } else {
          Setting.showMessage("error", `failed to save: ${res.msg}`);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `failed to save: ${error}`);
      });
  }

  deleteAsset() {
    AssetBackend.deleteAsset(this.state.asset)
      .then((res) => {
        if (res.status === "ok") {
          this.props.history.push("/assets");
        } else {
          Setting.showMessage("error", `${i18next.t("general:Failed to delete")}: ${res.msg}`);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `${i18next.t("general:Failed to connect to server")}: ${error}`);
      });
  }

  render() {
    return (
      <div>
        {
          this.state.asset !== null ? this.renderAsset() : null
        }
        <div style={{marginTop: "20px", marginLeft: "40px"}}>
          <Button size="large" onClick={() => this.submitAssetEdit(false)}>{i18next.t("general:Save")}</Button>
          <Button style={{marginLeft: "20px"}} type="primary" size="large" onClick={() => this.submitAssetEdit(true)}>{i18next.t("general:Save & Exit")}</Button>
          {this.state.mode === "add" ? <Button style={{marginLeft: "20px"}} size="large" onClick={() => this.deleteAsset()}>{i18next.t("general:Cancel")}</Button> : null}
        </div>
      </div>
    );
  }
}

export default AssetEditPage;
