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

package object

import (
	"fmt"

	"github.com/casvisor/casvisor/chain"
	"github.com/casvisor/casvisor/util"
)

func (record *Record) getRecordProvider() (*Provider, error) {
	if record.Provider != "" {
		provider, err := getProvider(record.Owner, record.Provider)
		if err != nil {
			return nil, err
		}

		if provider != nil {
			return provider, nil
		}
	}

	provider, err := getActiveBlockchainProvider(record.Owner)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (record *Record) getRecordChainClient() (chain.ChainClientInterface, error) {
	provider, err := record.getRecordProvider()
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, fmt.Errorf("there is no active blockchain provider")
	}

	client, err2 := chain.NewChainClient(provider.Type, provider.ClientId, provider.ClientSecret, provider.Region)
	if err2 != nil {
		return nil, err2
	}

	return client, nil
}

func (record *Record) toString() string {
	return util.StructToJson(record)
}

func CommitRecord(record *Record) (bool, error) {
	client, err := record.getRecordChainClient()
	if err != nil {
		return false, err
	}

	resp, err := client.Commit(record.toString())
	if err != nil {
		return false, err
	}

	if resp.Status != "ok" {
		return false, fmt.Errorf(resp.Msg)
	}

	record.Block = resp.Data
	return UpdateRecord(record.getId(), record)
}

func QueryRecord(id string) (string, error) {
	record, err := GetRecord(id)
	if err != nil {
		return "", err
	}
	if record == nil {
		return "", fmt.Errorf("the record: %s does not exist", id)
	}

	if record.Block == "" {
		return "", fmt.Errorf("the record: %s's block ID should not be empty", record.getId())
	}

	client, err := record.getRecordChainClient()
	if err != nil {
		return "", err
	}

	resp, err := client.Query(record.Block)
	if err != nil {
		return "", err
	}

	if resp.Status != "ok" {
		return "", fmt.Errorf(resp.Msg)
	}

	res := resp.Data
	return res, nil
}
