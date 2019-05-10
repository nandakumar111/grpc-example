package main

import (
	"context"
	"io"
	"log"
	"time"

	"../greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	creds, err := credentials.NewClientTLSFromFile("greet/ssl/server.crt", "")
	if err != nil {
		log.Fatalln(err)
	}
	_ = grpc.WithTransportCredentials(creds)
	client, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// Greetings
	conn := greetpb.NewGreetServiceClient(client)
	unaryGreet(conn)
	//serverStreamGreet(conn)
	//serverStreamPrimeNumberDecomposition(conn)
	//clienStreamLongGreet(conn)
	//doGreetEveryOne(conn)
	//doFindMaximum(conn)

}

func unaryGreet(client greetpb.GreetServiceClient) {
	log.Print("Starting to do a unary RPC..!!")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Nandakumar",
			LastName:  "R",
		},
	}

	res, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Erro while calling Greet RPC : %v", err)
	}
	log.Printf("Response from Greet : %v", res)
}

func serverStreamGreet(client greetpb.GreetServiceClient) {
	log.Println("Staring to do server streaming RPC..!!")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Nandakumar",
			LastName:  "R",
		},
	}

	resStream, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC : %v", err)
	}
	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming result from server : %v", err)
		}
		log.Printf("Response from GreetingManyTimes : %v", res.GetResult())
	}

}

func clienStreamLongGreet(client greetpb.GreetServiceClient) {
	log.Print("Starting to do client streaming RPC..!!")

	reqStream := []*greetpb.GreetRequest{
		&greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "nanda",
			},
		}, &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "nanda1",
			},
		}, &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "nanda2",
			},
		}, &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "nanda3",
			},
		}, &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "nanda4",
			},
		}, &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "nanda5",
			},
		}, &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "nanda6",
			},
		},
	}
	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet function : %v", err)
	}

	for _, req := range reqStream {
		log.Printf("Sending request : %v", req)
		_ = stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("Error : %v", err)
	}
	log.Printf("Response from Longgreet : %v", res)
}

func doGreetEveryOne(client greetpb.GreetServiceClient) {
	log.Print("Starting to do client and streaming RPC..!!")

	reqStream := []*greetpb.GreetRequest{
		&greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "nanda",
			},
		}, &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "nanda1",
			},
		}, &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "nanda2",
			},
		}, &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "nanda3",
			},
		}, &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "nanda4",
			},
		}, &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "nanda5",
			},
		}, &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "nanda6",
			},
		},
	}
	stream, err := client.GreetEveryOne(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetEveryOne function : %v", err)
	}
	waitc := make(chan struct{})

	go func() {
		for _, req := range reqStream {
			log.Printf("Sending request : %v", req)
			_ = stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
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
			log.Printf("Response : %v", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}

func doFindMaximum(client greetpb.GreetServiceClient) {
	log.Print("Starting to do client and streaming RPC..!!")

	reqStream := []*greetpb.Numbers{
		&greetpb.Numbers{Number: 1},
		&greetpb.Numbers{Number: 5},
		&greetpb.Numbers{Number: 3},
		&greetpb.Numbers{Number: 6},
		&greetpb.Numbers{Number: 2},
		&greetpb.Numbers{Number: 20},
	}
	stream, err := client.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while calling FindMaximum function : %v", err)
	}
	waitc := make(chan struct{})

	go func() {
		for _, req := range reqStream {
			log.Printf("Sending request : %v", req)
			_ = stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
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
			log.Printf("Response : %v", res.GetNumber())
		}
		close(waitc)
	}()

	<-waitc
}
