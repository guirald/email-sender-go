package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/guirald/email-sender-go/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUsersStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "João",
		Email: "j@j.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "j@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}
		fmt.Println("Status: ", stream.Status)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "g1",
			Name:  "guirald",
			Email: "guira@guirald.com",
		},
		&pb.User{
			Id:    "g2",
			Name:  "guirald 2",
			Email: "guira2@guirald.com",
		},
		&pb.User{
			Id:    "g3",
			Name:  "guirald 3",
			Email: "guira3@guirald.com",
		},
		&pb.User{
			Id:    "g4",
			Name:  "guirald 4",
			Email: "guira4@guirald.com",
		},
		&pb.User{
			Id:    "g5",
			Name:  "guirald 5",
			Email: "guira5@guirald.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUsersStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUsersStreamBoth((context.Background()))

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "g1",
			Name:  "guirald",
			Email: "guira@guirald.com",
		},
		&pb.User{
			Id:    "g2",
			Name:  "guirald 2",
			Email: "guira2@guirald.com",
		},
		&pb.User{
			Id:    "g3",
			Name:  "guirald 3",
			Email: "guira3@guirald.com",
		},
		&pb.User{
			Id:    "g4",
			Name:  "guirald 4",
			Email: "guira4@guirald.com",
		},
		&pb.User{
			Id:    "g5",
			Name:  "guirald 5",
			Email: "guira5@guirald.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	// em paralelo/concorrente
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Printf("Receiving user %v with status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
