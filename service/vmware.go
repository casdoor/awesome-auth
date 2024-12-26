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

package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type MachineVmwareClient struct {
	hostname     string
	basicAuthKey string
}

type VirtualMachine struct {
	ID     string `json:"id"`
	Cpu    Cpu    `json:"cpu"`
	Memory int    `json:"memory"`
}

type VirtualMachinePath struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type Cpu struct {
	Processors int `json:"processors"`
}

// hostname is the IP address of the target host and the configured port format is {IP}:{port}
// basicAuthKey is the credential for VMware Workstation Pro REST service, in the form of {username}:{password}
func newMachineVmwareClient(hostname string, basicAuthKey string) (*MachineVmwareClient, error) {
	client := MachineVmwareClient{
		hostname,
		basicAuthKey,
	}
	return &client, nil
}

func (client MachineVmwareClient) GetMachines() ([]*Machine, error) {
	url := fmt.Sprintf("http://%s/api/vms", client.hostname)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.vmware.vmw.rest-v1+json")
	auth := client.basicAuthKey
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encodedAuth)

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var vmList []VirtualMachinePath
	err = json.Unmarshal(body, &vmList)
	if err != nil {
		return nil, err
	}

	machines := []*Machine{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(vmList))
	for _, vm := range vmList {
		wg.Add(1)
		go func(vm VirtualMachinePath) {
			defer wg.Done()
			machine, err2 := client.GetMachine(vm.ID)
			if err2 != nil {
				errChan <- err2
				return
			}

			mu.Lock()
			machines = append(machines, machine)
			mu.Unlock()
		}(vm)
	}

	wg.Wait()
	close(errChan)
	for err := range errChan {
		return nil, err
	}

	return machines, nil
}

func (client MachineVmwareClient) GetMachine(name string) (*Machine, error) {
	url := fmt.Sprintf("http://%s/api/vms/%s", client.hostname, name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.vmware.vmw.rest-v1+json")
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(client.basicAuthKey))
	req.Header.Add("Authorization", "Basic "+encodedAuth)

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var vm VirtualMachine
	err = json.Unmarshal(body, &vm)
	if err != nil {
		return nil, err
	}

	machine := &Machine{
		Name:    name,
		Id:      vm.ID,
		CpuSize: fmt.Sprintf("%d", vm.Cpu.Processors),
		MemSize: fmt.Sprintf("%d", vm.Memory),
	}
	return machine, nil
}
