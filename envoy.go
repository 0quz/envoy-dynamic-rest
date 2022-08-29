package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var pl = fmt.Println

// listerner config begining
type ResponseListener struct {
	VersionInfo string            `json:"version_info"`
	Resources   ResourcesListener `json:"resources"`
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
	VersionInfo string           `json:"version_info"`
	Resources   ResourcesCluster `json:"resources"`
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
	VersionInfo       string            `json:"version_info"`
	ResourcesEndpoint ResourcesEndpoint `json:"resources"`
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

// json parameters
type ListenerParams struct {
	Name    string `json:"name"`
	CdsName string `json:"cds_name"`
}

type ClusterParams struct {
	Name    string `json:"name"`
	EdsName string `json:"eds_name"`
}

type EndpointParams struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	PortValue int    `json:"port_value"`
}

// dbop

type Lds struct {
	//gorm.Model //for creating automatic id / create / update / delete date
	Name       string `gorm:"primaryKey"`
	CdsName    string
	Configured bool `gorm:"default:false"`
}

type Cds struct {
	Name       string `gorm:"primaryKey"`
	EdsName    string
	Configured bool `gorm:"default:false"`
}

type Eds struct {
	Name       string `gorm:"primaryKey"`
	Address    string
	PortValue  int
	Configured bool `gorm:"default:false"`
}

// redis connection
var ctx = context.Background()

func connectRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}

func setRedisMemcached(key string, value string) {
	rdb := connectRedisClient()
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func getRedisMemcached(key string) string {
	rdb := connectRedisClient()
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		panic(err)
	} else {
		pl("key", val)
	}
	return val
}

func ldsAddDb(l *ListenerParams) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Lds{})

	// Create
	db.Create(&Lds{Name: l.Name, CdsName: l.CdsName})
}

func cdsAddDb(c *ClusterParams) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Cds{})

	// Create
	db.Create(&Cds{Name: c.Name, EdsName: c.EdsName})
}

func edsAddDb(e *EndpointParams) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Eds{})

	// Create
	db.Create(&Eds{Name: e.Name, PortValue: e.PortValue, Address: e.Address})

	// Read
	//var product Product
	//db.First(&product, 1)                 // find product with integer primary key
	//db.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	//db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	//db.Model(&product).Updates(Product{Price: 200, Code: "KEKW"}) // non-zero fields
	//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "KEKW"})

	// Delete - delete product
	//db.Delete(&product, 1)
}

func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func xdsConfig(c *fiber.Ctx) error {
	xds := c.Params("xds")
	if xds == "lds" {
		l := new(ListenerParams)
		if err := c.BodyParser(l); err != nil {
			return err
		}
		ldsAddDb(l)
	} else if xds == "cds" {
		cd := new(ClusterParams)
		if err := c.BodyParser(cd); err != nil {
			return err
		}
		cdsAddDb(cd)
	} else if xds == "eds" {
		e := new(EndpointParams)
		if err := c.BodyParser(e); err != nil {
			return err
		}
		edsAddDb(e)
	}
	return nil
}

func xdsUpdateConfig(c *fiber.Ctx) error {
	xds := c.Params("xds")
	if xds == "lds" {
		l := new(ListenerParams)
		if err := c.BodyParser(l); err != nil {
			return err
		}
		ldsUpdateDb(l)
	} else if xds == "cds" {
		cd := new(ClusterParams)
		if err := c.BodyParser(cd); err != nil {
			return err
		}
		cdsUpdateDb(cd)
	} else if xds == "eds" {
		e := new(EndpointParams)
		if err := c.BodyParser(e); err != nil {
			return err
		}
		edsUpdateDb(e)
	}
	return nil
}

func ldsUpdateDb(l *ListenerParams) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.Model(&Lds{}).Where("name = ?", l.Name).Updates(map[string]interface{}{"cds_name": l.CdsName, "configured": false}) // I have to use interface becase of boolean field update
}

func cdsUpdateDb(c *ClusterParams) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.Model(&Cds{}).Where("name = ?", c.Name).Updates(map[string]interface{}{"eds_name": c.EdsName, "configured": false})
}

func edsUpdateDb(e *EndpointParams) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.Model(&Eds{}).Where("name = ?", e.Name).Updates(map[string]interface{}{"address": e.Address, "port_value": e.PortValue, "configured": false})
}

func xds(c *fiber.Ctx) error {
	xds := c.Params("xds")
	if xds == ":listeners" {
		db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		var lds Lds
		db.Last(&lds)
		if lds.Configured {
			c.Status(304)
			return nil
		}
		var domains []string
		responseListener := &ResponseListener{
			VersionInfo: "0",
			Resources: ResourcesListener{
				Type: "type.googleapis.com/envoy.config.listener.v3.Listener",
				Name: lds.Name,
				Address: Address{
					SocketAddress{
						Address:   "0.0.0.0",
						PortValue: 10000,
					},
				},
				FilterChains: []FilterChain{
					{
						Filters: []Filter{
							{
								Name: "envoy.filters.network.http_connection_manager",
								TypedConfig: TypedConfig{
									Type:       "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager",
									CodecType:  "AUTO",
									StatPrefix: "ingress_http",
									HttpFilters: HttpFilters{
										Name: "envoy.filters.http.router",
										TypedConfig: TypedConfig2{
											Type: "type.googleapis.com/envoy.extensions.filters.http.router.v3.Router",
										},
									},
									RouteConfig: RouteConfig{
										Name: "local_route",
										VirtualHosts: VirtualHosts{
											Name:    "local_route",
											Domains: append(domains, "*"),
											Routes: Routes{
												Match: Match{
													Prefix: "/",
												},
												Route: Route{
													Cluster: lds.CdsName,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		lds.Configured = true
		db.Save(&lds)
		return c.JSON(responseListener)
	} else if xds == ":clusters" {
		db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		var cds Cds
		db.Last(&cds)
		if cds.Configured {
			c.Status(304)
			return nil
		}
		var clusterNames []string
		responseCluster := &ResponseCluster{
			VersionInfo: "0",
			Resources: ResourcesCluster{
				ClusterType:     "type.googleapis.com/envoy.config.cluster.v3.Cluster",
				Name:            cds.Name,
				Type:            "EDS",
				LbPolicy:        "ROUND_ROBIN",
				ConnectTimeout:  "0.25s",
				DnsLookupFamily: "V4_ONLY",
				EdsClusterConfig: EdsClusterConfig{
					ServiceName: "eds_service",
					EdsConfig: EdsConfig{
						ResourceApiVersion: "V3",
						ApiConfigSource: ApiConfigSource{
							ApiType:             "REST",
							TransportApiVersion: "V3",
							ClusterNames:        append(clusterNames, "xds_cluster"),
							RefreshDelay:        "3s",
						},
					},
				},
			},
		}
		cds.Configured = true
		db.Save(&cds)
		return c.JSON(responseCluster)
	} else if xds == ":endpoints" {
		db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		var eds Eds
		db.Last(&eds)
		if eds.Configured {
			c.Status(304)
			return nil
		}
		responseEndpoint := &ResponseEndpoint{
			VersionInfo: "0",
			ResourcesEndpoint: ResourcesEndpoint{
				Type:        "type.googleapis.com/envoy.config.endpoint.v3.ClusterLoadAssignment",
				ClusterName: "eds_service",
				Endpoints: Endpoints{
					LbEndpoints: []LbEndpoints{
						{
							Endpoint: Endpoint{
								Address: EndpointsAddress{
									SocketAddress: EndpointsSocketAddress{
										Address:   eds.Address,
										PortValue: eds.PortValue,
									},
								},
							},
						},
					},
				},
			},
		}
		eds.Configured = true
		db.Save(&eds)
		return c.JSON(responseEndpoint)
	}
	return nil
}

func main() {
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Post("/config/:xds", xdsConfig)
	app.Put("/config/:xds", xdsUpdateConfig)
	app.Post("/v3/discovery:xds", xds)
	app.Listen(":8080")
}
