package main

// listerner config begining
type ResponseListener struct {
	VersionInfo string              `json:"version_info"`
	Resources   []ResourcesListener `json:"resources"`
}

type ResourcesListener struct {
	Type         string        `json:"@type"`
	Name         string        `json:"name"`
	Address      Address       `json:"address"`
	FilterChains []FilterChain `json:"filter_chains"`
}

type Address struct {
	SocketAddress SocketAddress `json:"socket_address"`
}

type SocketAddress struct {
	Address   string `json:"address"`
	PortValue int    `json:"port_value"`
}

type FilterChain struct {
	Filters []Filter `json:"filters"`
}

type Filter struct {
	Name        string      `json:"name"`
	TypedConfig TypedConfig `json:"typed_config"`
}

type TypedConfig struct {
	Type        string      `json:"@type"`
	CodecType   string      `json:"codec_type"`
	StatPrefix  string      `json:"stat_prefix"`
	HttpFilters HttpFilters `json:"http_filters"`
	RouteConfig RouteConfig `json:"route_config"`
}

type HttpFilters struct {
	Name        string       `json:"name"`
	TypedConfig TypedConfig2 `json:"typed_config"`
}

type TypedConfig2 struct {
	Type string `json:"@type"`
}

type RouteConfig struct {
	Name         string       `json:"name"`
	VirtualHosts VirtualHosts `json:"virtual_hosts"`
}

type VirtualHosts struct {
	Name    string   `json:"name"`
	Domains []string `json:"domains"`
	Routes  Routes   `json:"routes"`
}

type Routes struct {
	Match Match `json:"match"`
	Route Route `json:"route"`
}

type Match struct {
	Prefix string `json:"prefix"`
}

type Route struct {
	Cluster string `json:"cluster"`
}

// cluster config begining
type ResponseCluster struct {
	VersionInfo string             `json:"version_info"`
	Resources   []ResourcesCluster `json:"resources"`
}

type ResourcesCluster struct {
	ClusterType      string           `json:"@type"`
	Name             string           `json:"name"`
	Type             string           `json:"type"`
	LbPolicy         string           `json:"lb_policy"`
	ConnectTimeout   string           `json:"connect_timeout"`
	DnsLookupFamily  string           `json:"dns_lookup_family"`
	EdsClusterConfig EdsClusterConfig `json:"eds_cluster_config"`
}

type EdsClusterConfig struct {
	ServiceName string    `json:"service_name"`
	EdsConfig   EdsConfig `json:"eds_config"`
}

type EdsConfig struct {
	ResourceApiVersion string          `json:"resource_api_version"`
	ApiConfigSource    ApiConfigSource `json:"api_config_source"`
}

type ApiConfigSource struct {
	ApiType             string   `json:"api_type"`
	TransportApiVersion string   `json:"transport_api_version"`
	ClusterNames        []string `json:"cluster_names"`
	RefreshDelay        string   `json:"refresh_delay"`
}

// endpoint config begining
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
