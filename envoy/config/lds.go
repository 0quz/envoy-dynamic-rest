package config

// lds response template
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
	Name            string          `json:"name"`
	HttpTypedConfig HttpTypedConfig `json:"typed_config"`
}

type HttpTypedConfig struct {
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
