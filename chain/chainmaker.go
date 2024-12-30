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

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

type ChainTencentChainmakerClient struct {
	Client *ecs.Client
}

func newChainTencentChainmakerClient(accessKeyId string, accessKeySecret string, region string) (ChainTencentChainmakerClient, error) {
	client, err := ecs.NewClientWithAccessKey(
		region,
		accessKeyId,
		accessKeySecret,
	)
	if err != nil {
		return ChainTencentChainmakerClient{}, err
	}

	return ChainTencentChainmakerClient{Client: client}, nil
}

func (client ChainTencentChainmakerClient) Commit(data string) (*Response, error) {
	return nil, fmt.Errorf("not implemented")
}

func (client ChainTencentChainmakerClient) Query(blockId string) (*Response, error) {
	return nil, fmt.Errorf("not implemented")
}
