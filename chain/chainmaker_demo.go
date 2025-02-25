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

package chain

import (
	"fmt"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tbaas "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tbaas/v20180416"
)

type ChainTencentChainmakerDemoClient struct {
	ClientId     string
	ClientSecret string
	Region       string
	NetworkId    string
	ChainId      string
	Client       *tbaas.Client
}

func newChainTencentChainmakerDemoClient(clientId, clientSecret, region, networkId, chainId string) (*ChainTencentChainmakerDemoClient, error) {
	credential := common.NewCredential(clientId, clientSecret)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tbaas.tencentcloudapi.com"

	client, err := tbaas.NewClient(credential, region, cpf)
	if err != nil {
		return nil, fmt.Errorf("newChainTencentChainmakerClient() error: %v", err)
	}

	return &ChainTencentChainmakerDemoClient{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Region:       region,
		NetworkId:    networkId,
		ChainId:      chainId,
		Client:       client,
	}, nil
}

func (client *ChainTencentChainmakerDemoClient) Commit(data string) (string, error) {
	request := tbaas.NewInvokeChainMakerDemoContractRequest()
	request.ClusterId = common.StringPtr(client.NetworkId)
	request.ChainId = common.StringPtr(client.ChainId)
	request.ContractName = common.StringPtr("ChainMakerDemo")
	request.FuncName = common.StringPtr("save")
	request.FuncParam = common.StringPtr(data)

	response, err := client.Client.InvokeChainMakerDemoContract(request)
	if err != nil {
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
			return "", fmt.Errorf("TencentCloudSDKError: %v", sdkErr)
		}

		return "", fmt.Errorf("ChainTencentChainmakerDemoClient.Client.Invoke() error: %v", err)
	}

	txId := *(response.Response.Result.TxId)
	return txId, nil
}

func (client ChainTencentChainmakerDemoClient) Query(blockId string, data map[string]string) (string, error) {
	// simulate the situation that error occurs
	if strings.HasSuffix(data["id"], "0") {
		return "", fmt.Errorf("some error occurred in the ChainTencentChainmakerDemoClient::Commit operation")
	}

	// Query the data from the blockchain
	// Write some code... (if error occurred, handle it as above)

	// assume the chain data are retrieved from the blockchain, here we just generate it statically
	chainData := map[string]string{"organization": "casbin"}

	// Check if the data are matched with the chain data
	res := "Matched"
	if chainData["organization"] != data["organization"] {
		res = "Mismatched"
	}

	// simulate the situation that mismatch occurs
	if strings.HasSuffix(blockId, "2") || strings.HasSuffix(blockId, "4") || strings.HasSuffix(blockId, "6") || strings.HasSuffix(blockId, "8") || strings.HasSuffix(blockId, "0") {
		res = "Mismatched"
	}

	return fmt.Sprintf("The query result for block [%s] is: %s", blockId, res), nil
}
