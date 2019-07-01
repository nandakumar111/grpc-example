# gRPC and gRPC-web connectivity via [Envoy Proxy](https://www.envoyproxy.io/) - #1831

##Why is envoy proxy required??
[grpc-web](https://github.com/grpc/grpc-web)leverages on http/2. A comparison with grpc will be something like REST and web-sockets.

gRPC is not faster than REST over HTTP/2 by default, but it gives the tools to make it faster. There are some things that would be difficult or impossible to do with REST.

- `Selective message compression.` In gRPC, a streaming RPC can decide to compress or not compress messages. For example, if we are streaming mixed text and images over a single stream (or really any mixed compressible content), we can turn off compression for the images. This saves from compressing already compressed data which won't get any smaller but will burn up our server CPU.
- `First class load balancing.` While not an improvement in point to point connections, gRPC can intelligently pick which backend to send traffic to. (this is a library feature, not a wire protocol feature). This means we can send our requests to the least loaded backend server without resorting to using a proxy. This is a latency win.
- `Heavily optimized.` gRPC (the library) is under continuous benchmarks to ensure that there are no speed regressions. Those benchmarks are improving constantly. Again, this doesn't have anything to do with gRPC the protocol, but our program will be faster for having used gRPC.
- We will see most of our performance improvement just from using Protobuf. While we could be used proto with REST, it is very nicely integrated with gRPC. Technically, we could use JSON with gRPC, but most people don't want to pay the performance cost after getting used to protos.

Basically, grpc and http2 are not the same that's why we need a proxy.
For more details,
- https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-WEB.md
- https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md

[**Reference**](https://stackoverflow.com/questions/53051648/why-is-envoy-proxy-required-for-grpc-web)

> Example envoy proxy server(`envoy.yaml`) creation in Docker 

```
admin:
  access_log_path: /tmp/admin_access.log
  address:
    socket_address: { address: 0.0.0.0, port_value: 9901 }
## grpc_backend_server_port : 8080, envoy_proxy_port : 9090
static_resources:
  listeners:
    - name: listener_0
      address:
        socket_address: { address: 0.0.0.0, port_value: 9090 }
      filter_chains:
        - filters:
            - name: envoy.http_connection_manager
              config:
                codec_type: auto
                stat_prefix: ingress_http
                route_config:
                  name: local_route
                  virtual_hosts:
                    - name: local_service
                      domains: ["*"]
                      routes:
                        - match: { prefix: "/" }
                          route:
                            cluster: ping_pong_service
                            max_grpc_timeout: 0s
                      cors:
                        allow_origin:
                          - "*"
                        allow_headers: keep-alive,user-agent,cache-control,content-type,content-transfer-encoding,x-accept-content-transfer-encoding,x-accept-response-streaming,x-user-agent,x-grpc-web,grpc-timeout
                        expose_headers: grpc-status,grpc-message
                http_filters:
                  - name: envoy.grpc_web
                  - name: envoy.cors
                  - name: envoy.router
  clusters:
    - name: fogfind_service
      connect_timeout: 0.25s
      type: logical_dns
      http2_protocol_options: {}
      lb_policy: round_robin
      hosts: [{ socket_address: { address: host.docker.internal, port_value: 8080 }}]
```
**Command to run Envoy Proxy server in docker**
```
docker build -t fogfind/grpc-web-server-envoy-proxy
docker run -d -p 9090:9090 fogfind/grpc-web-server-envoy-proxy
```
