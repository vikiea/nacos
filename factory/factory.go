package factory

import (
	"errors"
	"net/url"

	"github.com/vikiea/registry"
	"github.com/vikiea/registry/nacos"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
)

func NewRegistrar(uri *url.URL, groupName string) (registry.Registrar, error) {
	switch uri.Scheme {
	case nacos.Scheme:
		return nacos.NewRegistrar(func() (naming_client.INamingClient, error) {
			return nacos.NewNacosNamingClient(uri)
		}, nacos.GroupName(groupName))
	default:
		return nil, errors.New("not support this scheme " + uri.Scheme)
	}
}

func NewDiscovery(uri *url.URL, groupName string) (registry.Discovery, error) {
	switch uri.Scheme {
	case nacos.Scheme:
		return nacos.NewDiscovery(func() (naming_client.INamingClient, error) {
			return nacos.NewNacosNamingClient(uri)
		}, nacos.GroupName(groupName))
	default:
		return nil, errors.New("not support this scheme " + uri.Scheme)
	}
}
