package main

import (
	"envoy/dbop"
	"envoy/redis"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var pl = fmt.Println

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func bodyParser(c *fiber.Ctx, out interface{}) error {
	if err := c.BodyParser(out); err != nil {
		return err
	}
	return nil
}

func xdsConfig(c *fiber.Ctx) error {
	//c.Accepts("application/json")
	method := string(c.Request().Header.Method())
	xdsType := c.Params("xds")
	var err error
	if method == "POST" {
		if xdsType == "lds" {
			lds := new(dbop.ListenerRequestJson)
			bodyParser(c, lds)
			err = dbop.AddLds(lds)
			if err != nil {
				return err
			}
		} else if xdsType == "cds" {
			cds := new(dbop.ClusterRequestJson)
			bodyParser(c, cds)
			err = dbop.AddCds(cds)
			if err != nil {
				return err
			}
		} else if xdsType == "eds" {
			eds := new(dbop.EndpointRequestJson)
			bodyParser(c, eds)
			err = dbop.AddEds(eds)
			if err != nil {
				return err
			}
		} else if xdsType == "endpoint" {
			eds := new(dbop.EndpointRequestJson)
			bodyParser(c, eds)
			err = dbop.AddEndpointAddress(eds)
			if err != nil {
				return err
			}
		}
	} else if method == "PUT" {
		if xdsType == "lds" {
			lds := new(dbop.ListenerRequestJson)
			bodyParser(c, lds)
			err = dbop.UpdateLds(lds)
			if err != nil {
				return err
			}
		} else if xdsType == "cds" {
			cds := new(dbop.ClusterRequestJson)
			bodyParser(c, cds)
			err = dbop.UpdateCds(cds)
			if err != nil {
				return err
			}
		}
	} else if method == "DELETE" {
		if xdsType == "lds" {
			lds := new(dbop.ListenerRequestJson)
			bodyParser(c, lds)
			err = dbop.DeleteLds(lds)
			if err != nil {
				return err
			}
		} else if xdsType == "cds" {
			cds := new(dbop.ClusterRequestJson)
			bodyParser(c, cds)
			err = dbop.DeleteCds(cds)
			if err != nil {
				return err
			}
		} else if xdsType == "eds" {
			eds := new(dbop.EndpointRequestJson)
			bodyParser(c, eds)
			err = dbop.DeleteEds(eds)
			if err != nil {
				return err
			}
		} else if xdsType == "endpoint" {
			eds := new(dbop.EndpointRequestJson)
			bodyParser(c, eds)
			err = dbop.DeleteEndpointAddress(eds)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func xds(c *fiber.Ctx) error {
	xds := c.Params("xds")
	if xds == ":listeners" {
		deployed := redis.GetRedisMemcached("ldsDeployed")
		if deployed == "yes" {
			c.Status(304)
			return nil
		}
		db := dbop.ConnectPostgresClient()
		var lds []dbop.Lds
		err := db.Model(lds).Preload("Cds.Eds").Find(&lds).Error
		if err != nil {
			c.Status(204)
			return nil
		}
		var responseData []ResourcesListener
		for _, l := range lds {
			var domains []string
			resources := ResourcesListener{
				Type: "type.googleapis.com/envoy.config.listener.v3.Listener",
				Name: l.Name,
				Address: Address{
					SocketAddress: SocketAddress{
						Address:   l.Address,
						PortValue: l.PortValue,
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
		responseListener := ResponseListener{
			VersionInfo: "1",
			Resources:   responseData,
		}
		redis.SetRedisMemcached("ldsDeployed", "yes")
		redis.SetRedisMemcached("cdsDeployed", "no")
		return c.JSON(responseListener)
	} else if xds == ":clusters" {
		deployed := redis.GetRedisMemcached("cdsDeployed")
		if deployed == "yes" {
			c.Status(304)
			return nil
		}
		db := dbop.ConnectPostgresClient()
		var lds []dbop.Lds
		err := db.Model(lds).Preload("Cds.Eds").Find(&lds).Error // nested table access.
		if err != nil {
			c.Status(204)
			return nil
		}
		var responseData []ResourcesCluster
		for _, l := range lds {
			if l.Cds.Name != "" {
				var clusterNames []string
				resources := ResourcesCluster{
					ClusterType:     "type.googleapis.com/envoy.config.cluster.v3.Cluster",
					Name:            l.CdsName,
					Type:            "EDS",
					LbPolicy:        "ROUND_ROBIN",
					ConnectTimeout:  "0.25s",
					DnsLookupFamily: "V4_ONLY",
					EdsClusterConfig: EdsClusterConfig{
						ServiceName: l.Cds.EdsName,
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
				}
				responseData = append(responseData, resources)
			} else {
				c.Status(204)
				return nil
			}
		}
		responseCluster := ResponseCluster{
			VersionInfo: "1",
			Resources:   responseData,
		}
		redis.SetRedisMemcached("cdsDeployed", "yes")
		redis.SetRedisMemcached("edsDeployed", "no")
		return c.JSON(responseCluster)
	} else if xds == ":endpoints" {
		deployed := redis.GetRedisMemcached("edsDeployed")
		if deployed == "yes" {
			c.Status(304)
			return nil
		}
		db := dbop.ConnectPostgresClient()
		var lds []dbop.Lds
		err := db.Model(lds).Preload("Cds.Eds").Find(&lds).Error // nested table access.
		if err != nil {
			c.Status(204)
			return nil
		}
		var lbEndpointsData []LbEndpoints
		var resourcesEndpoint []ResourcesEndpoint
		for _, l := range lds {
			if l.Cds.Eds.Name != "" {
				var eA []dbop.EndpointAddress
				err = db.Where("eds_name = ?", l.Cds.Eds.Name).Find(&eA).Error
				if err != nil {
					c.Status(204)
					return nil
				}
				for _, e := range eA {
					lbEndpoints := LbEndpoints{
						Endpoint: Endpoint{
							Address: EndpointsAddress{
								SocketAddress: EndpointsSocketAddress{
									Address:   e.Address,
									PortValue: e.PortValue,
								},
							},
						},
					}
					lbEndpointsData = append(lbEndpointsData, lbEndpoints)
				}
				resourcesData := ResourcesEndpoint{
					Type:        "type.googleapis.com/envoy.config.endpoint.v3.ClusterLoadAssignment",
					ClusterName: l.Cds.Eds.Name,
					Endpoints: Endpoints{
						LbEndpoints: lbEndpointsData,
					},
				}
				resourcesEndpoint = append(resourcesEndpoint, resourcesData)
			} else {
				c.Status(204)
				return nil
			}
		}
		responseEndpoint := ResponseEndpoint{
			VersionInfo: "1",
			Resources:   resourcesEndpoint,
		}
		redis.SetRedisMemcached("edsDeployed", "yes")
		return c.JSON(responseEndpoint)
	}
	return nil
}
