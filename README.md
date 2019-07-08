# **gRPC Documentation**


## **What is gRPC ?**
[gRPC](https://grpc.io/) is a modern open-source high-performance RPC framework that can run in any environment. It can efficiently connect services in and across data centers with pluggable support for load balancing, tracing, health checking and authentication. It is also applicable in the last mile of distributed computing to connect devices, mobile applications, and browsers to backend services.

## **[gRPC Example](https://github.com/nandakumar111/grpc-example)**

## **Types of RPC**

* Unary
* Server streaming
* Client streaming
* Bi-Directional streaming

### **Proto File**

[Proto](https://developers.google.com/protocol-buffers/docs/proto3) file contains message(s) and service(s). Message(s), It contains message attributes.

```
...
message Nodes{
    string attr1 = 1;
    int32 attr2 = 2;
} 
message NestedNodes {
    message NestedOne{
        message NestedTwo{
           string attr1 = 1;
        }
        NestedTwo nested_two = 1;
    }
    repeated int32 attr1 = 1;
    NestedOne nested_one = 2;
}
...
```
Proto service(s), It contains number of rpc services, Here we can define whether it's Unary or streaming like that.
```
service Services { 
    // Request and Response both are messages
    rpc Unary (Request) returns (Response);
    
    rpc ServerStreaming (Request) returns (Response);

    rpc ClientStreaming (Request) returns (Response);

    rpc BiDirectionalStreaming (Request) returns (Response);
   
}
```
Proto Generation

```$ protoc {source_file_path}/{file_name}.proto --go_out=plugins=grpc:{destination_file_path}```

## **Service Declaration**
```
//server
type Server {}

func main(){
    ...
    s := grpc.NewServer()
    proto.RegisterServicesServer(s, &Server{})
    ...
}
```
```
//client
...
func main(){
    client, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
	log.Fatalln(err)
    }
    defer client.Close()
    
   conn := proto.NewServicesClient(client)
   ...
}
```


## **Unary Call**

```server.go```

```
...
func (*Server) Unary (ctx context.Context, req *proto.Request) (*proto.Response, error){
    ...
    return &Response{},nil
}
...
```
```client.go```
```
...
res, err := conn.Unary(context.Background(), &proto.Request{})
   if err != nil {
	log.Fatalf("Error while calling RPC : %v", err)
   }
...
```


## **Server Streaming Call**

```server.go```

```
...
func (*Server) ServerStreaming (req *proto.Request, stream proto.Services_ServerStreamingServer) error{
    ...
    // streaming loop start
        _ = stream.Send(&proto.Response{})
    // streaming loop end
    return nil
}
...
```
```client.go```
```
...
resStream, err := conn.ServerStreaming(context.Background(), &proto.Request{})
if err != nil {
		log.Fatalf("Error while calling ServerStreaming RPC : %v", err)
	}
	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming result from server : %v", err)
		}
		log.Printf("Response from ServerStreaming : %v", res)
	}
...
```


## **Client Streaming Call**

```server.go```

```
...
func (*Server) ClientStreaming (reqStream proto.Services_ServerStreamingServer) error{
   
    for{
       req, err := reqStream.Recv()
       		if err == io.EOF {
       			_ = reqStream.SendAndClose(&proto.Response{})
       			return nil
       		}
       		if err != nil {
       			return err
       		}
    } 
    return nil
}
...
```
```client.go```
```
...
reqStream := []*proto.Request{
   &proto.Request{},
   &proto.Request{},
   ...
}

stream, err := conn.ClientStreaming(context.Background())
if err != nil {
		log.Fatalf("Error while calling ClientStreaming function : %v", err)
	}

	for _, req := range reqStream {
		log.Printf("Sending request : %v", req)
		_ = stream.Send(req)
		//time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("Error : %v", err)
	}
	log.Printf("Response from ClientStreaming : %v", res)
...
```

## **Bi-Directional Stream**

```server.go```
```
func (*server) BiDirectionalStreaming(stream proto.Services_BiDirectionalStreamingServer) error {
   ...
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error while calling BiDirectionalStreaming function : %v", err)
			return err
		}
		_ = stream.Send(&proto.Response{})
	}
}
```
```client.go```
```
...
reqStream := []*proto.Request{
   &proto.Request{},
   &proto.Request{},
   ...
}
	stream, err := client.BiDirectionalStreaming(context.Background())
	if err != nil {
		log.Fatalf("Error while calling BiDirectionalStreaming function : %v", err)
	}
	waitc := make(chan struct{})

	go func() {
		for _, req:= range reqStream {
			log.Printf("Sending request : %v", req)
			_ = stream.Send(req)
			//time.Sleep(1000 * time.Millisecond)
		}
		_ = stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while Receiving : %v", err)
			}
			log.Printf("Response : %v", res)
		}
		close(waitc)
	}()

	<-waitc
	...
```

## **[Transcoding HTTP/JSON to gRPC](https://cloud.google.com/endpoints/docs/grpc/transcoding)**


Cloud Endpoints supports protocol transcoding so that clients can access your gRPC API by using HTTP/JSON. The Extensible Service Proxy (ESP) transcodes HTTP/JSON to gRPC.

**Reference:**

* Basic **[Cloud Endpoints for gRPC APIS](https://cloud.google.com/endpoints/docs/grpc/about-grpc)**

```<proto_file_name>.proto```
```
import "google/api/annotations.proto";

service Services { 

      rpc Service1 (PostRequest) returns (Response){
          option (google.api.http) = {
              post : "/api/v1/service1"
              body : "*"
           };
      }
      
      rpc Service2 (PutRequest) returns (Response){
            option (google.api.http) = {
                put : "/api/v1/service2"
                body : "*"
             };
      }
      
      rpc Service3 (Request) returns (Response){
            option (google.api.http) = {
                get : "/api/v1/service3/{msg}"
             };
       }
       
       rpc Service4 (Request) returns (Response){
            option (google.api.http) = {
                delete : "/api/v1/service4/{msg}"
             };
       }
}

message PostRequest {
    ...
}

message PutRequest {
    ...
}

message Request {
    string msg = 1;
}
...
```
Run command to generate grpc files and gateway of transcoding of HTTP/JSON
```
// proto file and gateway file generation
protoc -I. -I. \
    -I${GOPATH}/src \
    -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --go_out=plugins=grpc:. \
    <proto_file_name>.proto

protoc -I. -I. \
    -I${GOPATH}/src \
    -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --grpc-gateway_out=logtostderr=true:. \
    <proto_file_name>.proto

protoc -I. -I. \
    -I${GOPATH}/src \
    -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --swagger_out=logtostderr=true:. \
    <proto_file_name>.proto
```
# gRPC and gRPC-web connectivity via [Envoy Proxy](https://www.envoyproxy.io/) 

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
                            cluster: auth_service
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
    - name: auth_service
      connect_timeout: 0.25s
      type: logical_dns
      http2_protocol_options: {}
      lb_policy: round_robin
      hosts: [{ socket_address: { address: host.docker.internal, port_value: 8080 }}]
```
**Command to run Envoy Proxy server in docker**
```
docker build -t nandakumar111/grpc-web-server-envoy-proxy
docker run -d -p 9090:9090 nandakumar111/grpc-web-server-envoy-proxy
```
