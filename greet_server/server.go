package main

import (
	"context"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"../greetpb"
	"google.golang.org/grpc"
)

type server struct{}

// Unary
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("Greet function invoked  with %v", req)
	return &greetpb.GreetResponse{Result: "Hi " + req.GetGreeting().GetFirstName() + " " + req.GetGreeting().GetLastName()}, nil
}

// Server Stream
func (*server) GreetManyTimes(req *greetpb.GreetRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Printf("GreetManyTimes function invoked  with %v", req)
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetResponse{
			Result: "Hello " + req.GetGreeting().GetFirstName() + " " + req.GetGreeting().GetLastName() + " : number = " + strconv.Itoa(i),
		}
		_ = stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

// Client Streaming
func (*server) LongGreet(reqStream greetpb.GreetService_LongGreetServer) error{
	log.Printf("LongGreet function invoked  with client streaming..!!")
	result := ""
	for {
		req, err := reqStream.Recv()
		if err == io.EOF {
			_ = reqStream.SendAndClose(&greetpb.GreetResponse{Result: result})
			return nil
		}
		if err != nil {
			return err
		}
		result += "Hello " + req.GetGreeting().GetFirstName() + "! "
	}
}

// Bi-Directional Streaming
func (*server) GreetEveryOne(stream greetpb.GreetService_GreetEveryOneServer) error {
	log.Printf("GreetEveryOne function invoked  with client and server streaming..!!")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error while calling LongGreet function : %v", err)
			return err
		}
		result := "Hello " + req.GetGreeting().GetFirstName() + "! "
		_ = stream.Send(&greetpb.GreetResponse{Result: result})
	}
}


func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}
	var options []grpc.ServerOption
	tls := false
	if tls{
		certFile := "greet/ssl/server.crt"
		keyFile := "greet/ssl/server.pem"

		creds, err:= credentials.NewServerTLSFromFile(certFile,keyFile)
		if err != nil{
			log.Fatalf("SSL file err : %v", err)
		}
		options = append(options, grpc.Creds(creds))
	}

	s := grpc.NewServer()

	// Greeting
	greetpb.RegisterGreetServiceServer(s, &server{})

	log.Print("Greeting service run successfully..!!")
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Server Error: %v", err)
	}
}
