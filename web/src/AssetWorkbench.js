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

class AssetWorkbench extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      selectedAsset: null,
      assetTreeWidth: 240,
      isFullScreen: false,
    };
  }

  handleAssetSelect = (assetId) => {
    const arr = assetId.split("/");
    const asset = {owner: arr[0], name: arr[1]};
    this.setState({
      selectedAsset: asset,
    });
  };

  toggleFullScreen = () => {
    this.setState({
      isFullScreen: !this.state.isFullScreen,
    });
  };

  render() {
    const {assetTreeWidth} = this.state;

    return (
      <div
        style={{
          display: "flex",
          height: "100vh",
          background: "#f5f5f5",
        }}>
        <div style={{width: assetTreeWidth, background: "#fff"}}>
          <AssetTree onSelect={this.handleAssetSelect} account={this.props.account} />
        </div>
        <div style={{width: "100%", overflow: "hidden"}} >
          <RemoteDesktop
            asset={this.state.selectedAsset}
            toggleFullscreen={this.toggleFullScreen}
          />
        </div>
      </div>
    );
  }
}

export default AssetWorkbench;
