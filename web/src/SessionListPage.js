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

import * as SessionBackend from "./backend/SessionBackend";
import * as Setting from "./Setting";
import {Table} from "antd";
import i18next from "i18next";
import PopconfirmModal from "./common/modal/PopconfirmModal";
import {Link} from "react-router-dom";
import BaseListPage from "./BaseListPage";
import moment from "moment";

class SessionListPage extends BaseListPage {
  deleteSession(i) {
    const session = this.state.sessions[i];
    SessionBackend.deleteSession(session).then((res) => {
      if (res.status === "ok") {
        Setting.showMessage("success", "Successfully deleted session");
        this.setState({
          sessions: Setting.deleteRow(this.state.sessions, i),
        });
      } else {
        Setting.showMessage("error", `Failed to delete session: ${res.msg}`);
      }
    });
  }

  renderTable(sessions) {
    const columns = [
      {
        title: i18next.t("general:Name"),
        dataIndex: "name",
        key: "name",
        width: "180px",
        sorter: (a, b) => a.name.localeCompare(b.name),
        render: (text, record, index) => {
          return (
            <Link to={`/sessions/${record.name}`}>{text}</Link>
          );
        },
      },
      {
        title: i18next.t("general:Protocol"),
        dataIndex: "protocol",
        key: "protocol",
        width: "50px",
        filterMultiple: false,
        filters: [
          {text: "RDP", value: "RDP"},
          {text: "VNC", value: "VNC"},
          {text: "SSH", value: "SSH"},
        ],
      },
      {
        title: i18next.t("general:IP"),
        dataIndex: "ip",
        key: "ip",
        width: "120px",
      },
      {
        title: i18next.t("general:Connected time"),
        dataIndex: "connectedTime",
        key: "connectedTime",
        width: "200px",
        sorter: (a, b) => a.connectedTime.localeCompare(b.connectedTime),
      },
      {
        title: i18next.t("general:Connected time duration"),
        dataIndex: "connectedTimeDur",
        key: "connectedTimeDur",
        width: "200px",
        render: (text, record) => {
          if (!record["connectedTime"]) {
            return "-";
          }
          const connectedTime = moment(record["connectedTime"]);
          const currentTime = moment();
          const duration = moment.duration(currentTime.diff(connectedTime));
          return `${duration.hours()}h ${duration.minutes()}m ${duration.seconds()}s`;
        },
      },
      {
        title: i18next.t("general:Actions"),
        dataIndex: "",
        key: "op",
        width: "180px",
        fixed: (Setting.isMobile()) ? "false" : "right",
        render: (text, record, index) => {
          return (
            <div>
              <PopconfirmModal
                title={i18next.t("general:Sure to delete") + `: ${record.name} ?`}
                onConfirm={() => this.deleteSession(index)}
              >
              </PopconfirmModal>
            </div>
          );
        },
      },
    ];

    const paginationProps = {
      total: this.state.pagination.total,
      showQuickJumper: true,
      showSizeChanger: true,
      showTotal: () => i18next.t("general:{total} in total").replace("{total}", this.state.pagination.total),
    };

    return (
      <Table scroll={{x: "max-content"}} columns={columns} dataSource={sessions} rowKey={(record) => `${record.owner}/${record.name}`} size="middle" bordered
        pagination={paginationProps}
        title={() => (
          <div>
            {i18next.t("general:Sessions")}&nbsp;&nbsp;&nbsp;&nbsp;
          </div>
        )}
        loading={this.state.loading}
        onChange={this.handleTableChange}
      />
    );
  }

  fetch = (params = {}) => {
    let field = params.searchedColumn, value = params.searchText;
    const sortField = params.sortField, sortOrder = params.sortOrder;
    if (params.type !== undefined && params.type !== null) {
      field = "type";
      value = params.type;
    }

    this.setState({
      loading: true,
    });

    SessionBackend.getSessions(Setting.isDefaultOrganizationSelected(this.props.account) ? "" : Setting.getRequestOrganization(this.props.account), params.pagination.current, params.pagination.pageSize, field, value, sortField, sortOrder).then((res) => {
      this.setState({
        loading: false,
      });

      if (res.status === "ok") {
        this.setState({
          data: res.data,
          pagination: {
            ...params.pagination,
            total: res.data2,
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

export default SessionListPage;
