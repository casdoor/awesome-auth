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

import React from "react";
import {Link} from "react-router-dom";
import {Button, Switch, Table} from "antd";
import BaseListPage from "./BaseListPage";
import moment from "moment";
import * as Setting from "./Setting";
import * as AssetBackend from "./backend/AssetBackend";
import i18next from "i18next";
import PopconfirmModal from "./common/modal/PopconfirmModal";

class AssetListPage extends BaseListPage {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      assets: null,
    };
  }

  UNSAFE_componentWillMount() {
    this.getAssets();
  }

  getAssets() {
    AssetBackend.getAssets(this.props.account.owner)
      .then((res) => {
        if (res.status === "ok") {
          this.setState({
            assets: res.data,
          });
        } else {
          Setting.showMessage("error", `Failed to get assets: ${res.msg}`);
        }
      });
  }

  newAsset() {
    return {
      owner: this.props.account.owner,
      name: `machine_${this.state.assets.length}`,
      createdTime: moment().format(),
      description: `New Machine - ${this.state.assets.length}`,
      protocol: "rdp",
      ip: "127.0.0.1",
      port: 3389,
      username: "Administrator",
      password: "123",
      language: "zh",
      autoQuery: false,
      isPermanent: true,
      // remoteAppName:"",
      // remoteAppDir:"",
      // remoteAppArgs:"",
      remoteApps: [],
      services: [],
    };
  }

  addAsset() {
    const newAsset = this.newAsset();
    AssetBackend.addAsset(newAsset)
      .then((res) => {
        if (res.status === "ok") {
          this.props.history.push({pathname: `/assets/${newAsset.owner}/${newAsset.name}`, mode: "add"});
          Setting.showMessage("success", "Asset added successfully");
        } else {
          Setting.showMessage("error", `Failed to add Asset: ${res.msg}`);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `Asset failed to add: ${error}`);
      });
  }

  deleteAsset(i) {
    AssetBackend.deleteAsset(this.state.assets[i])
      .then((res) => {
        if (res.status === "ok") {
          Setting.showMessage("success", "Asset deleted successfully");
          this.setState({
            assets: Setting.deleteRow(this.state.assets, i),
            // pagination: {total: this.state.pagination.total - 1},
          });
        } else {
          Setting.showMessage("error", `Failed to delete Asset: ${res.msg}`);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `Asset failed to delete: ${error}`);
      });
  }

  renderTable(assets) {
    const columns = [
      {
        title: i18next.t("general:Organization"),
        dataIndex: "owner",
        key: "owner",
        width: "110px",
        sorter: true,
        ...this.getColumnSearchProps("owner"),
        render: (text, asset, index) => {
          return (
            <a target="_blank" rel="noreferrer" href={Setting.getMyProfileUrl(this.props.account).replace("/account", `/organizations/${text}`)}>
              {text}
            </a>
          );
        },
      },
      {
        title: i18next.t("general:Name"),
        dataIndex: "name",
        key: "name",
        width: "120px",
        sorter: true,
        ...this.getColumnSearchProps("name"),
        render: (text, record, index) => {
          return (
            <Link to={`/assets/${record.owner}/${record.name}`}>{text}</Link>
          );
        },
      },
      {
        title: i18next.t("general:Created time"),
        dataIndex: "createdTime",
        key: "createdTime",
        width: "160px",
        // sorter: true,
        sorter: (a, b) => a.createdTime.localeCompare(b.createdTime),
        render: (text, asset, index) => {
          return Setting.getFormattedDate(text);
        },
      },
      {
        title: i18next.t("general:Description"),
        dataIndex: "description",
        key: "description",
        // width: '200px',
        sorter: (a, b) => a.description.localeCompare(b.description),
      },
      {
        title: i18next.t("general:Protocol"),
        dataIndex: "protocol",
        key: "protocol",
        width: "50px",
        sorter: true,
        filterMultiple: false,
        filters: [
          {text: "RDP", value: "RDP"},
          {text: "VNC", value: "VNC"},
          {text: "SSH", value: "SSH"},
          {text: "", value: "-"},
        ],
      },
      {
        title: i18next.t("general:IP"),
        dataIndex: "ip",
        key: "ip",
        width: "120px",
        sorter: (a, b) => a.ip.localeCompare(b.ip),
      },
      {
        title: i18next.t("general:Port"),
        dataIndex: "port",
        key: "port",
        width: "90px",
        sorter: (a, b) => a.port - b.port,
      },
      {
        title: i18next.t("general:Username"),
        dataIndex: "username",
        key: "username",
        width: "130px",
        sorter: (a, b) => a.username.localeCompare(b.username),
      },
      {
        title: i18next.t("general:Language"),
        dataIndex: "language",
        key: "language",
        width: "90px",
        sorter: (a, b) => a.language.localeCompare(b.language),
      },
      {
        title: i18next.t("general:Auto query"),
        dataIndex: "autoQuery",
        key: "autoQuery",
        width: "100px",
        render: (text, record, index) => {
          return (
            <Switch disabled checked={text} />
          );
        },
      },
      {
        title: i18next.t("general:Is permanent"),
        dataIndex: "isPermanent",
        key: "isPermanent",
        width: "110px",
        render: (text, record, index) => {
          return (
            <Switch disabled checked={text} />
          );
        },
      },
      {
        title: i18next.t("general:Enable Remote App"),
        dataIndex: "enableRemoteApp",
        key: "enableRemoteApp",
        width: "150px",
        render: (text, record, index) => {
          return (
            <Switch disabled checked={text} />
          );
        },
      },
      {
        title: i18next.t("general:Remote Apps"),
        dataIndex: "remoteApps",
        key: "remoteApps",
        width: "120px",
        // todo: fix filter
        render: (text, record, index) => {
          return `${record.enableRemoteApp ? 1 : 0}  / ${record.remoteApps === null ? 0 : record.remoteApps.length}`;
        },
      },
      {
        title: i18next.t("general:Services"),
        dataIndex: "services",
        key: "services",
        width: "90px",
        // todo: fix filter
        render: (text, record, index) => {
          return `${record.services.filter(service => service.status === "Running").length} / ${record.services.length}`;
        },
      },
      {
        title: i18next.t("general:Action"),
        dataIndex: "action",
        key: "action",
        width: "260px",
        fixed: (Setting.isMobile()) ? "false" : "right",
        render: (text, asset, index) => {
          return (
            <div>
              <Button
                disabled={!Setting.isAdminUser(this.props.account) && (asset.owner !== this.props.account.owner)}
                style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}}
                type="primary"
                onClick={() => {
                  let link = `access?owner=${asset.owner}&name=${asset.name}&protocol=${asset.protocol}`;
                  if (asset.enableRemoteApp) {
                    link += `&remoteApp=${asset.remoteApps[0].remoteAppName}&remoteAppDir=${asset.remoteApps[0].remoteAppDir}&remoteAppArgs=${asset.remoteApps[0].remoteAppArgs}`;
                  }
                  Setting.openLink(link);
                }}
              >
                {i18next.t("general:Connect")}
              </Button>
              <Button
                disabled={!Setting.isAdminUser(this.props.account) && (asset.owner !== this.props.account.owner)}
                style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}}
                onClick={() => this.props.history.push(`/assets/${asset.owner}/${asset.name}`)}
              >{i18next.t("general:Edit")}
              </Button>
              <PopconfirmModal
                disabled={!Setting.isAdminUser(this.props.account) && (asset.owner !== this.props.account.owner)}
                title={i18next.t("general:Sure to delete") + `: ${asset.name} ?`}
                onConfirm={() => this.deleteAsset(index)}
              >
              </PopconfirmModal>
            </div>
          );
        },
      },
    ];

    const paginationProps = {
      showQuickJumper: true,
      showSizeChanger: true,
    };

    return (
      <div>
        <Table scroll={{x: "max-content"}} columns={columns} dataSource={assets} rowKey={(asset) => `${asset.owner}/${asset.name}`} size="middle" bordered pagination={paginationProps}
          title={() => (
            <div>
              {i18next.t("general:Assets")}&nbsp;&nbsp;&nbsp;&nbsp;
              <Button type="primary" size="small" disabled={!Setting.isAdminUser(this.props.account)} onClick={this.addAsset.bind(this)}>{i18next.t("general:Add")}</Button>
            </div>
          )}
          loading={this.state.loading}
          onChange={this.handleTableChange}
        />
      </div>
    );
  }

  render() {
    return (
      <div>
        {
          this.renderTable(this.state.assets)
        }
      </div>
    );
  }

  fetch = (params = {}) => {
    let field = params.searchedColumn, value = params.searchText;
    const sortField = params.sortField, sortOrder = params.sortOrder;
    if (params.type !== undefined && params.type !== null) {
      field = "type";
      value = params.type;
    }
    this.setState({loading: true});
    AssetBackend.getAssets(Setting.getRequestOrganization(this.props.account), params.pagination.current, params.pagination.pageSize, field, value, sortField, sortOrder)
      .then((res) => {
        this.setState({
          loading: false,
        });
        if (res.status === "ok") {
          this.setState({
            data: res.data,
            pagination: {
              ...params.pagination,
              // total: res.data2,
            },
            searchText: params.searchText,
            searchedColumn: params.searchedColumn,
          });
        } else {
          if (Setting.isResponseDenied(res)) {
            this.setState({
              isAuthorized: false,
            });
          } else {
            Setting.showMessage("error", res.msg);
          }
        }
      });
  };
}

export default AssetListPage;
