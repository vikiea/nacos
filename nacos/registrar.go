package nacos

import (
	"context"
	"errors"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/vikiea/registry"
	"github.com/vikiea/registry/pkg/stringx"
)

const instanceID = "InstanceID"

type Registrar struct {
	namingClient naming_client.INamingClient
	nacosOptions *nacosOptions
}

func (r *Registrar) Register(ctx context.Context, instance registry.ServiceInstance) error {
	if stringx.IsBlank(instance.Name()) {
		return errors.New("no service to register for nacos client")
	}
	return r.register(ctx, instance)
}

func (r *Registrar) Deregister(ctx context.Context, instance registry.ServiceInstance) error {
	if stringx.IsBlank(instance.Name()) {
		return errors.New("no service to register for nacos client")
	}
	return r.deregister(ctx, instance)
}

func (r *Registrar) register(_ context.Context, instance registry.ServiceInstance) error {
	metadata := instance.Metadata()
	if stringx.IsNotBlank(instance.ID()) {
		metadata[instanceID] = instance.ID()
	}
	param := vo.RegisterInstanceParam{
		Ip:          instance.IP(),           // 服务实例IP
		Port:        uint64(instance.Port()), // 服务实例port
		Weight:      r.nacosOptions.Weight,   // 权重
		Enable:      true,                    // 是否上线
		Healthy:     true,                    // 是否健康
		Metadata:    metadata,                // 扩展信息
		ClusterName: r.nacosOptions.ClusterName,
		ServiceName: instance.Name(),
		GroupName:   groupName(r.nacosOptions, instance),
		Ephemeral:   true, // 是否临时实例
	}
	ok, err := r.namingClient.RegisterInstance(param)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("failed to register instance")
	}
	return nil
}

func (r *Registrar) deregister(_ context.Context, instance registry.ServiceInstance) error {
	param := vo.DeregisterInstanceParam{
		Ip:          instance.IP(),           // 服务实例IP
		Port:        uint64(instance.Port()), // 服务实例port
		Cluster:     r.nacosOptions.ClusterName,
		ServiceName: instance.Name(),
		GroupName:   groupName(r.nacosOptions, instance),
		Ephemeral:   true, // 是否临时实例
	}
	ok, err := r.namingClient.DeregisterInstance(param)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("failed to deregister instance")
	}
	return nil
}

func NewRegistrar(factory NamingClientFactoryFunc, opts ...Option) (*Registrar, error) {
	namingClient, err := factory.Create()
	if err != nil {
		return nil, err
	}
	r := &Registrar{
		namingClient: namingClient,
		nacosOptions: &nacosOptions{
			Clusters:    nil,
			ClusterName: "",
			Weight:      10,
			Namespace:   "",
			GroupName:   "",
		},
	}
	for _, opt := range opts {
		opt(r.nacosOptions)
	}
	return r, nil
}
