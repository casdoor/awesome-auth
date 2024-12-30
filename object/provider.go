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

	"github.com/casvisor/casvisor/util"
	"xorm.io/core"
)

type Provider struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	UpdatedTime string `xorm:"varchar(100)" json:"updatedTime"`
	DisplayName string `xorm:"varchar(100)" json:"displayName"`

	Category string `xorm:"varchar(100)" json:"category"`
	Type     string `xorm:"varchar(100)" json:"type"`

	ClientId     string `xorm:"varchar(100)" json:"clientId"`
	ClientSecret string `xorm:"varchar(100)" json:"clientSecret"`
	Region       string `xorm:"varchar(100)" json:"region"`
	Network      string `xorm:"varchar(100)" json:"network"`
	Chain        string `xorm:"varchar(100)" json:"chain"`

	State string `xorm:"varchar(100)" json:"state"`
}

func GetProviderCount(owner, field, value string) (int64, error) {
	session := GetSession(owner, -1, -1, field, value, "", "")
	return session.Count(&Provider{})
}

func GetProviders(owner string) ([]*Provider, error) {
	providers := []*Provider{}
	err := adapter.engine.Desc("created_time").Find(&providers, &Provider{Owner: owner})
	if err != nil {
		return providers, err
	}

	return providers, nil
}

func getActiveProviders(owner string) ([]*Provider, error) {
	providers, err := GetProviders(owner)
	if err != nil {
		return nil, err
	}

	res := []*Provider{}
	for _, provider := range providers {
		if provider.ClientId != "" && provider.ClientSecret != "" && provider.State == "Active" {
			res = append(res, provider)
		}
	}
	return res, nil
}

func GetPaginationProviders(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*Provider, error) {
	providers := []*Provider{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&providers)
	if err != nil {
		return providers, err
	}

	return providers, nil
}

func getProvider(owner string, name string) (*Provider, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	provider := Provider{Owner: owner, Name: name}
	existed, err := adapter.engine.Get(&provider)
	if err != nil {
		return &provider, err
	}

	if existed {
		return &provider, nil
	} else {
		return nil, nil
	}
}

func GetProvider(id string) (*Provider, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getProvider(owner, name)
}

func GetMaskedProvider(provider *Provider, errs ...error) (*Provider, error) {
	if len(errs) > 0 && errs[0] != nil {
		return nil, errs[0]
	}

	if provider == nil {
		return nil, nil
	}

	if provider.ClientSecret != "" {
		provider.ClientSecret = "***"
	}
	return provider, nil
}

func GetMaskedProviders(providers []*Provider, errs ...error) ([]*Provider, error) {
	if len(errs) > 0 && errs[0] != nil {
		return nil, errs[0]
	}

	var err error
	for _, provider := range providers {
		provider, err = GetMaskedProvider(provider)
		if err != nil {
			return nil, err
		}
	}

	return providers, nil
}

func UpdateProvider(id string, provider *Provider) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	p, err := getProvider(owner, name)
	if err != nil {
		return false, err
	} else if p == nil {
		return false, nil
	}

	if provider.ClientSecret == "***" {
		provider.ClientSecret = p.ClientSecret
	}

	affected, err := adapter.engine.ID(core.PK{owner, name}).AllCols().Update(provider)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func AddProvider(provider *Provider) (bool, error) {
	affected, err := adapter.engine.Insert(provider)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func DeleteProvider(provider *Provider) (bool, error) {
	affected, err := adapter.engine.ID(core.PK{provider.Owner, provider.Name}).Delete(&Provider{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func (provider *Provider) getId() string {
	return fmt.Sprintf("%s/%s", provider.Owner, provider.Name)
}
