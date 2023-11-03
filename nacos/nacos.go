package nacos

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/vikiea/registry"
)

type NamingClientFactory interface {
	Create() (naming_client.INamingClient, error)
}

type NamingClientFactoryFunc func() (naming_client.INamingClient, error)

func (f NamingClientFactoryFunc) Create() (naming_client.INamingClient, error) {
	return f()
}

const Scheme = "nacos"

type nacosOptions struct {
	Clusters    []string
	ClusterName string
	Weight      float64
	Namespace   string
	GroupName   string
}

type Option func(r *nacosOptions)

func ClusterName(clusterName string) Option {
	return func(r *nacosOptions) {
		r.ClusterName = clusterName
	}
}

func Clusters(clusters []string) Option {
	return func(r *nacosOptions) {
		r.Clusters = clusters
	}
}

func Weight(weight float64) Option {
	return func(r *nacosOptions) {
		r.Weight = weight
	}
}

func NameSpace(nameSpace string) Option {
	return func(r *nacosOptions) {
		r.Namespace = nameSpace
	}
}

func GroupName(name string) Option {
	return func(r *nacosOptions) {
		r.GroupName = name
	}
}

func groupName(o *nacosOptions, instance registry.ServiceInstance) string {
	groupName := o.GroupName
	if len(groupName) <= 0 {
		groupName = instance.Scheme()
	}
	return groupName
}

// NewNacosNamingClient 实例化客户端,可根据自身需求进行自定义
func NewNacosNamingClient(uri *url.URL) (naming_client.INamingClient, error) {
	hosts := strings.Split(uri.Host, ",")
	var sc = make([]constant.ServerConfig, 0, len(hosts))
	for _, host := range hosts {
		ipAddr, portStr, err := net.SplitHostPort(host)
		if err != nil {
			return nil, fmt.Errorf("failed split host port %s, %w", host, err)
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("failed convert port to int %s, %w", portStr, err)
		}
		sc = append(sc, *constant.NewServerConfig(ipAddr, uint64(port)))
	}

	namespaceID := uri.Query().Get("namespaceId")
	cc := constant.NewClientConfig(
		constant.WithNotLoadCacheAtStart(true),
		constant.WithNamespaceId(namespaceID),
		constant.WithTimeoutMs(5000),
	)

	clientParam := vo.NacosClientParam{ClientConfig: cc, ServerConfigs: sc}
	client, err := clients.NewNamingClient(clientParam)
	if err != nil {
		return nil, fmt.Errorf("failed new nacos client, %w", err)
	}
	return client, err
}
