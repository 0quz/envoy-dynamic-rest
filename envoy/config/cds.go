package config

// cds response template
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
