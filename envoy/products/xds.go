package products

import (
	"envoy/config"
	"envoy/dbop"
	"envoy/redis"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var pl = fmt.Println

func bodyParser(c *fiber.Ctx, out interface{}) error {
	if err := c.BodyParser(out); err != nil {
		return err
	}
	return nil
}

func (h handler) XdsConfig(c *fiber.Ctx) error {
	//c.Accepts("application/json")
	method := string(c.Request().Header.Method())
	xdsType := c.Params("xds")
	var err error
	if method == "POST" {
		if xdsType == "lds" {
			lds := new(dbop.ListenerRequestJson)
			bodyParser(c, lds)
			err = dbop.AddLds(lds, h.DB)
			if err != nil {
				return err
			}
		} else if xdsType == "cds" {
			cds := new(dbop.ClusterRequestJson)
			bodyParser(c, cds)
			err = dbop.AddCds(cds, h.DB)
			if err != nil {
				return err
			}
		} else if xdsType == "eds" {
			eds := new(dbop.EndpointRequestJson)
			bodyParser(c, eds)
			err = dbop.AddEds(eds, h.DB)
			if err != nil {
				return err
			}
		} else if xdsType == "endpoint" {
			eds := new(dbop.EndpointRequestJson)
			bodyParser(c, eds)
			err = dbop.AddEndpointAddress(eds, h.DB)
			if err != nil {
				return err
			}
		}
	} else if method == "PUT" {
		if xdsType == "lds" {
			lds := new(dbop.ListenerRequestJson)
			bodyParser(c, lds)
			err = dbop.UpdateLds(lds, h.DB)
			if err != nil {
				return err
			}
		} else if xdsType == "cds" {
			cds := new(dbop.ClusterRequestJson)
			bodyParser(c, cds)
			err = dbop.UpdateCds(cds, h.DB)
			if err != nil {
				return err
			}
		}
	} else if method == "DELETE" {
		if xdsType == "lds" {
			lds := new(dbop.ListenerRequestJson)
			bodyParser(c, lds)
			err = dbop.DeleteLds(lds, h.DB)
			if err != nil {
				return err
			}
		} else if xdsType == "cds" {
			cds := new(dbop.ClusterRequestJson)
			bodyParser(c, cds)
			err = dbop.DeleteCds(cds, h.DB)
			if err != nil {
				return err
			}
		} else if xdsType == "eds" {
			eds := new(dbop.EndpointRequestJson)
			bodyParser(c, eds)
			err = dbop.DeleteEds(eds, h.DB)
			if err != nil {
				return err
			}
		} else if xdsType == "endpoint" {
			eds := new(dbop.EndpointRequestJson)
			bodyParser(c, eds)
			err = dbop.DeleteEndpointAddress(eds, h.DB)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h handler) Xds(c *fiber.Ctx) error {
	xds := c.Params("xds")
	if xds == ":listeners" {
		deployed := redis.GetRedisMemcached("ldsDeployed")
		if deployed == "yes" {
			c.Status(fiber.StatusNotModified)
			return nil
		}
		var lds []dbop.Lds
		err := h.DB.Model(lds).Find(&lds).Error
		if err != nil {
			c.Status(fiber.StatusNoContent)
			return nil
		}
		var responseData []config.ResourcesListener
		for _, l := range lds {
			var domains []string
			resources := config.ResourcesListener{
				Type: "type.googleapis.com/envoy.config.listener.v3.Listener",
				Name: l.Name,
				Address: config.Address{
					SocketAddress: config.SocketAddress{
						Address:   l.Address,
						PortValue: l.PortValue,
					},
				},
				FilterChains: []config.FilterChain{
					{
						Filters: []config.Filter{
							{
								Name: "envoy.filters.network.http_connection_manager",
								TypedConfig: config.TypedConfig{
									Type:       "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager",
									CodecType:  "AUTO",
									StatPrefix: "ingress_http",
									HttpFilters: config.HttpFilters{
										Name: "envoy.filters.http.router",
										HttpTypedConfig: config.HttpTypedConfig{
											Type: "type.googleapis.com/envoy.extensions.filters.http.router.v3.Router",
										},
									},
									RouteConfig: config.RouteConfig{
										Name: "local_route",
										VirtualHosts: config.VirtualHosts{
											Name:    "local_route",
											Domains: append(domains, "*"),
											Routes: config.Routes{
												Match: config.Match{
													Prefix: "/",
												},
												Route: config.Route{
													Cluster: l.CdsName,
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
			responseData = append(responseData, resources)
		}
		responseListener := config.ResponseListener{
			VersionInfo: "1",
			Resources:   responseData,
		}
		redis.SetRedisMemcached("ldsDeployed", "yes")
		return c.JSON(responseListener)
	} else if xds == ":clusters" {
		deployed := redis.GetRedisMemcached("cdsDeployed")
		if deployed == "yes" {
			c.Status(fiber.StatusNotModified)
			return nil
		}
		var lds []dbop.Lds
		err := h.DB.Model(lds).Preload("Cds").Find(&lds).Error // nested table access.
		if err != nil {
			c.Status(fiber.StatusNoContent)
			return nil
		}
		var responseData []config.ResourcesCluster
		for _, l := range lds {
			pl(l.Cds.Name)
			var clusterNames []string
			resources := config.ResourcesCluster{
				ClusterType:     "type.googleapis.com/envoy.config.cluster.v3.Cluster",
				Name:            l.CdsName,
				Type:            "EDS",
				LbPolicy:        "ROUND_ROBIN",
				ConnectTimeout:  "0.25s",
				DnsLookupFamily: "V4_ONLY",
				EdsClusterConfig: config.EdsClusterConfig{
					ServiceName: l.Cds.EdsName,
					EdsConfig: config.EdsConfig{
						ResourceApiVersion: "V3",
						ApiConfigSource: config.ApiConfigSource{
							ApiType:             "REST",
							TransportApiVersion: "V3",
							ClusterNames:        append(clusterNames, "xds_cluster"),
							RefreshDelay:        "3s",
						},
					},
				},
			}
			responseData = append(responseData, resources)
		}
		responseCluster := config.ResponseCluster{
			VersionInfo: "1",
			Resources:   responseData,
		}
		redis.SetRedisMemcached("cdsDeployed", "yes")
		redis.SetRedisMemcached("edsDeployed", "no")
		return c.JSON(responseCluster)
	} else if xds == ":endpoints" {
		deployed := redis.GetRedisMemcached("edsDeployed")
		if deployed == "yes" {
			c.Status(fiber.StatusNotModified)
			return nil
		}
		var lds []dbop.Lds
		err := h.DB.Model(lds).Preload("Cds.Eds").Find(&lds).Error // nested table access.
		if err != nil {
			c.Status(fiber.StatusNoContent)
			return nil
		}
		var lbEndpointsData []config.LbEndpoints
		var resourcesEndpoint []config.ResourcesEndpoint
		for _, l := range lds {
			var eA []dbop.EndpointAddress
			err = h.DB.Where("eds_name = ?", l.Cds.Eds.Name).Find(&eA).Error
			if err != nil {
				c.Status(fiber.StatusNoContent)
				return nil
			}
			for _, e := range eA {
				lbEndpoints := config.LbEndpoints{
					Endpoint: config.Endpoint{
						Address: config.EndpointsAddress{
							SocketAddress: config.EndpointsSocketAddress{
								Address:   e.Address,
								PortValue: e.PortValue,
							},
						},
					},
				}
				lbEndpointsData = append(lbEndpointsData, lbEndpoints)
			}
			resourcesData := config.ResourcesEndpoint{
				Type:        "type.googleapis.com/envoy.config.endpoint.v3.ClusterLoadAssignment",
				ClusterName: l.Cds.Eds.Name,
				Endpoints: config.Endpoints{
					LbEndpoints: lbEndpointsData,
				},
			}
			resourcesEndpoint = append(resourcesEndpoint, resourcesData)
		}
		responseEndpoint := config.ResponseEndpoint{
			VersionInfo: "1",
			Resources:   resourcesEndpoint,
		}
		redis.SetRedisMemcached("edsDeployed", "yes")
		return c.JSON(responseEndpoint)
	}
	return nil
}
