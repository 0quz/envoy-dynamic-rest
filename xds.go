package main

import (
	"envoy/dbop"
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
	var lds *dbop.ListenerRequestJson
	var cds *dbop.ClusterRequestJson
	var eds *dbop.EndpointRequestJson
	if method == "POST" {
		if xdsType == "lds" {
			lds = new(dbop.ListenerRequestJson)
			bodyParser(c, lds)
			err = dbop.AddLds(lds)
			if err != nil {
				return err
			}
		} else if xdsType == "cds" {
			cds = new(dbop.ClusterRequestJson)
			bodyParser(c, cds)
			err = dbop.AddCds(cds)
			if err != nil {
				return err
			}
		} else if xdsType == "eds" {
			eds = new(dbop.EndpointRequestJson)
			bodyParser(c, eds)
			err = dbop.AddEds(eds)
			if err != nil {
				return err
			}
		} else if xdsType == "endpoint" {
			eds = new(dbop.EndpointRequestJson)
			bodyParser(c, eds)
			err = dbop.AddEndpointAddress(eds)
			if err != nil {
				return err
			}
		}
	} else if method == "PUT" {
		if xdsType == "lds" {
			lds = new(dbop.ListenerRequestJson)
			bodyParser(c, lds)
			err = dbop.UpdateLds(lds)
			if err != nil {
				return err
			}
		} else if xdsType == "cds" {
			cds = new(dbop.ClusterRequestJson)
			bodyParser(c, cds)
			err = dbop.UpdateCds(cds)
			if err != nil {
				return err
			}
		}
	} else if method == "DELETE" {
		if xdsType == "lds" {
			lds = new(dbop.ListenerRequestJson)
			bodyParser(c, lds)
			err = dbop.DeleteLds(lds)
			if err != nil {
				return err
			}
		} else if xdsType == "cds" {
			cds = new(dbop.ClusterRequestJson)
			bodyParser(c, cds)
			err = dbop.DeleteCds(cds)
			if err != nil {
				return err
			}
		} else if xdsType == "eds" {
			eds = new(dbop.EndpointRequestJson)
			bodyParser(c, eds)
			err = dbop.DeleteEds(eds)
			if err != nil {
				return err
			}
		} else if xdsType == "endpoint" {
			eds = new(dbop.EndpointRequestJson)
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
		db := dbop.ConnectPostgresClient()
		var lds []dbop.Lds
		db.Table("lds").Where("deployed = ?", false).Find(&lds)
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
		return c.JSON(responseListener)
	} else if xds == ":clusters" {
		db := dbop.ConnectPostgresClient()
		var lds []dbop.Lds
		err := db.Model(lds).Where("deployed = ?", false).Preload("Cds").Find(&lds).Error // nested table access.
		errCheck(err)
		var responseData []ResourcesCluster
		for _, l := range lds {
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
		}
		responseCluster := ResponseCluster{
			VersionInfo: "1",
			Resources:   responseData,
		}
		return c.JSON(responseCluster)
	} else if xds == ":endpoints" {
		db := dbop.ConnectPostgresClient()
		var lds []dbop.Lds
		err := db.Model(lds).Where("deployed = ?", false).Preload("Cds.Eds").Find(&lds).Error // nested table access.
		errCheck(err)
		var lbEndpointsData []LbEndpoints
		var resourcesEndpoint []ResourcesEndpoint
		for _, l := range lds {
			var Ed []dbop.EndpointAddress
			err = db.Where("eds_name = ?", l.Cds.Eds.Name).Find(&Ed).Error
			if err != nil {
				return err
			}
			for _, e := range Ed {
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
		}

		responseEndpoint := ResponseEndpoint{
			VersionInfo: "1",
			Resources:   resourcesEndpoint,
		}
		return c.JSON(responseEndpoint)
	}
	return nil
}
