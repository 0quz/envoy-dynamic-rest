package products

import (
	"envoy/config"
	"envoy/dbop"
	"envoy/redis"

	"github.com/gofiber/fiber/v2"
)

// Binds the request body to a struct.
func bodyParser(c *fiber.Ctx, out interface{}) error {
	if err := c.BodyParser(out); err != nil {
		return err
	}
	return nil
}

// Envoy config route
func (h handler) XdsConfig(c *fiber.Ctx) error {
	//c.Accepts("application/json")
	// Get request method type
	method := string(c.Request().Header.Method())
	// get xds value
	xdsType := c.Params("xds")
	var err error
	// Add xds configs
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
	} else if method == "PUT" { // Update xds configs
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
	} else if method == "DELETE" { // Delete xds configs
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
	if xds == ":listeners" { // Configuration source: https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/listener/v3/listener.proto#config-listener-v3-listener
		// If LDS was deployed return 304
		if redis.GetRedisMemcached("ldsDeployed") == "yes" {
			c.Status(fiber.StatusNotModified)
			return nil
		}
		// Get all LDS
		var lds []dbop.Lds
		err := h.DB.Model(lds).Find(&lds).Error
		if err != nil {
			return err
		}
		// LDS configuration part
		var responseData []config.ResourcesListener
		for _, l := range lds {
			var domains []string
			resources := config.ResourcesListener{
				Type: "type.googleapis.com/envoy.config.listener.v3.Listener",
				Name: l.Name, // LDS "name": "l1",
				Address: config.Address{
					SocketAddress: config.SocketAddress{
						Address:   l.Address,   // Domain address "address": 0.0.0.0
						PortValue: l.PortValue, // Domain port "port_value": 20000
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
													Cluster: l.CdsName, // Bind CDS "cds_name": "c1",
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
		// Set EDS deployed status yes to prevent unnecessary DB operation if EDS is updated successfully.
		redis.SetRedisMemcached("ldsDeployed", "yes")
		return c.JSON(responseListener)
	} else if xds == ":clusters" { // Configuration source: https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/cluster/v3/cluster.proto#config-cluster-v3-cluster
		// If Cds was deployed return 304
		if redis.GetRedisMemcached("cdsDeployed") == "yes" {
			c.Status(fiber.StatusNotModified)
			return nil
		}
		// Get all CDS from parent LDS
		var lds []dbop.Lds
		err := h.DB.Model(lds).Preload("Cds").Find(&lds).Error // nested table access.
		if err != nil {
			return nil
		}
		// CDS configuration part
		var responseData []config.ResourcesCluster
		for _, l := range lds {
			var clusterNames []string
			resources := config.ResourcesCluster{
				ClusterType:     "type.googleapis.com/envoy.config.cluster.v3.Cluster",
				Name:            l.CdsName, // CDS "name": "c1",
				Type:            "EDS",
				LbPolicy:        "ROUND_ROBIN",
				ConnectTimeout:  "0.25s",
				DnsLookupFamily: "V4_ONLY",
				EdsClusterConfig: config.EdsClusterConfig{
					ServiceName: l.Cds.EdsName, // Bind EDS "eds_name": "e1"
					EdsConfig: config.EdsConfig{
						ResourceApiVersion: "V3",
						ApiConfigSource: config.ApiConfigSource{
							ApiType:             "REST",
							TransportApiVersion: "V3",
							ClusterNames:        append(clusterNames, "xds_cluster"), // Static bind cluster source: "envoy.yaml"
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
		// Set CDS deployed status yes to prevent unnecessary DB operation if CDS is updated successfully.
		redis.SetRedisMemcached("cdsDeployed", "yes")
		// When you update CDS. Envoy needs to reconfigure EDS to CDS
		redis.SetRedisMemcached("edsDeployed", "no")
		return c.JSON(responseCluster)
	} else if xds == ":endpoints" { // Configuration source: https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/endpoint/v3/endpoint.proto#config-endpoint-v3-clusterloadassignment
		// If EDS was deployed return 304
		if redis.GetRedisMemcached("edsDeployed") == "yes" {
			c.Status(fiber.StatusNotModified)
			return nil
		}
		// Get all EDS from grandparent LDS
		var lds []dbop.Lds
		err := h.DB.Model(lds).Preload("Cds.Eds").Find(&lds).Error // nested table access.
		if err != nil {
			return nil
		}
		// EDS configuration part
		var lbEndpointsData []config.LbEndpoints
		var resourcesEndpoint []config.ResourcesEndpoint
		for _, l := range lds {
			// Get EndpointAddress from matching EDS cluster.
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
								Address:   e.Address,   // routing address "address": "192.168.65.2",
								PortValue: e.PortValue, // routing port "port_value": 1200,1400
							},
						},
					},
				}
				lbEndpointsData = append(lbEndpointsData, lbEndpoints)
			}
			resourcesData := config.ResourcesEndpoint{
				Type:        "type.googleapis.com/envoy.config.endpoint.v3.ClusterLoadAssignment",
				ClusterName: l.Cds.Eds.Name, // EDS "name": "e1",
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
		// Set EDS deployed status yes to prevent unnecessary DB operation if EDS is updated successfully.
		redis.SetRedisMemcached("edsDeployed", "yes")
		return c.JSON(responseEndpoint)
	}
	return nil
}
