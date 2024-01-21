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

import React from "react";
import AssetTree from "./component/access/AssetTree";
import RemoteDesktop from "./component/access/RemoteDesktop";
import * as AssetBackend from "./backend/AssetBackend";
import * as Setting from "./Setting";

class AssetWorkbench extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      selectedAsset: null,
      width: 0,
      height: 0,
      ref: React.createRef(),
    };
  }

  componentDidMount() {
    this.setState({
      width: this.state.ref.current.offsetWidth,
      height: this.state.ref.current.offsetHeight,
    });
  }

  handleAssetSelect = (assetId) => {
    const arr = assetId.split("/");
    AssetBackend.getAsset(arr[0], arr[1]).then((res) => {
      if (res.status === "ok") {
        this.setState({
          selectedAsset: res.data,
        });
      } else {
        Setting.showMessage("error", `Failed to get asset: ${res.msg}`);
      }
    });
  };

  render() {
    return (
      <div style={{
        display: "flex",
        height: "100vh",
        background: "#f5f5f5",
      }}>
        <div style={{flex: "0 0 15%",
          background: "#fff",
          overflow: "auto"}}>
          <AssetTree onSelect={this.handleAssetSelect} account={this.props.account} />
        </div>
        <div style={{
          flex: "1"}}
        ref={this.state.ref}>
          <RemoteDesktop
            asset={this.state.selectedAsset}
            width={this.state.width}
            height={this.state.height}
          />
        </div>
      </div>
    );
  }
}

export default AssetWorkbench;
