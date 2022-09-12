package config

// eds response template
type ResponseEndpoint struct {
	VersionInfo string              `json:"version_info"`
	Resources   []ResourcesEndpoint `json:"resources"`
}

type ResourcesEndpoint struct {
	Type        string    `json:"@type"`
	ClusterName string    `json:"cluster_name"`
	Endpoints   Endpoints `json:"endpoints"`
}

type Endpoints struct {
	LbEndpoints []LbEndpoints `json:"lb_endpoints"`
}

type LbEndpoints struct {
	Endpoint Endpoint `json:"endpoint"`
}

type Endpoint struct {
	Address EndpointsAddress `json:"address"`
}

type EndpointsAddress struct {
	SocketAddress EndpointsSocketAddress `json:"socket_address"`
}

type EndpointsSocketAddress struct {
	Address   string `json:"address"`
	PortValue int    `json:"port_value"`
}
