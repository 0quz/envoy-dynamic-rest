# Prerequisites

[Docker](https://www.docker.com/get-started)

<br/>

# Run

Run `make d.build` after that `make d.up` on the terminal when you are in the project directory.

## Suggested Example

1. Create an EDS via these URLs:

    `curl -X POST -H "Content-Type: application/json" -d '{"name":"e1"}' 0.0.0.0:8080/conf/eds`

2. Create endpoints of EDS via these URLs:

    1. `curl -X POST -H "Content-Type: application/json" -d '{"name": "e1", "address": "192.168.65.2", "port_value": 1200}' 0.0.0.0:8080/conf/endpoint`
    2. `curl -X POST -H "Content-Type: application/json" -d '{"name": "e1", "address": "192.168.65.2", "port_value": 1400}' 0.0.0.0:8080/conf/endpoint`

3. Create a CDS via this url:
  
    `curl -X POST -H "Content-Type: application/json" -d '{"name": "c1", "eds_name": "e1"}' 0.0.0.0:8080/conf/cds`

3. Create a LDS via this url:

    `curl -X POST -H "Content-Type: application/json" -d '{"name": "l1", "cds_name": "c1", "port_value": 20000}' 0.0.0.0:8080/conf/lds`
 
4. Finally, you should `curl -X GET 0.0.0.0:20000` then you will see the routing.

Feel free to change all configurations. You can use add, update delete methods on the same XDS endpoints.

## Note: 
  1. Updates are not available for EDS and endpoints. 
  2. Addresses of EDS endpoint may vary depending on your operating system. (192.168.65.2 for mac) See envoy.yaml for other container addresses.
