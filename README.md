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
