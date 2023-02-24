// Copyright (c) 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package access

import (
	"github.com/gardener/gardener-extension-provider-openstack/pkg/openstack/client"
	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/sharenetworks"
)

// ShareFileSystemAccess provides methods for managing share networks
type ShareFileSystemAccess interface {
	// ShareNetworks
	CreateShareNetwork(desired *sharenetworks.ShareNetwork) (*sharenetworks.ShareNetwork, error)
	GetShareNetworkByID(id string) (*sharenetworks.ShareNetwork, error)
	GetShareNetworksByName(name string) ([]*sharenetworks.ShareNetwork, error)
	UpdateShareNetwork(desired, current *sharenetworks.ShareNetwork) (modified bool, err error)
	DeleteShareNetwork(id string) error
}

type shareFileSystemAccess struct {
	shareFileSystem client.SharedFileSystem
	log             logr.Logger
}

var _ ShareFileSystemAccess = &shareFileSystemAccess{}

// NewShareFileSystemAccess creates a new access object
func NewShareFileSystemAccess(shareFileSystem client.SharedFileSystem, log logr.Logger) (ShareFileSystemAccess, error) {
	return &shareFileSystemAccess{
		shareFileSystem: shareFileSystem,
		log:             log,
	}, nil
}

func (a *shareFileSystemAccess) CreateShareNetwork(desired *sharenetworks.ShareNetwork) (*sharenetworks.ShareNetwork, error) {
	createOpts := sharenetworks.CreateOpts{
		Name:            desired.Name,
		Description:     desired.Description,
		NeutronNetID:    desired.NeutronNetID,
		NeutronSubnetID: desired.NeutronSubnetID,
	}

	return a.shareFileSystem.CreateShareNetwork(createOpts)
}

func (a *shareFileSystemAccess) GetShareNetworkByID(id string) (*sharenetworks.ShareNetwork, error) {
	sg, err := a.shareFileSystem.GetShareNetwork(id)
	return sg, client.IgnoreNotFoundError(err)
}

func (a *shareFileSystemAccess) GetShareNetworksByName(name string) ([]*sharenetworks.ShareNetwork, error) {
	raw, err := a.shareFileSystem.GetShareNetworksByName(name)
	if err != nil {
		return nil, client.IgnoreNotFoundError(err)
	}
	var results []*sharenetworks.ShareNetwork
	for i := range raw {
		item := raw[i]
		results = append(results, &item)
	}
	return results, nil
}

func (a *shareFileSystemAccess) UpdateShareNetwork(desired, current *sharenetworks.ShareNetwork) (modified bool, err error) {
	if current.Description == desired.Description {
		return
	}
	modified = true
	updateOpts := sharenetworks.UpdateOpts{
		Description: &desired.Description,
	}
	_, err = a.shareFileSystem.UpdateShareNetwork(current.ID, updateOpts)
	return
}

func (a *shareFileSystemAccess) DeleteShareNetwork(id string) error {
	return client.IgnoreNotFoundError(a.shareFileSystem.DeleteShareNetwork(id))
}
