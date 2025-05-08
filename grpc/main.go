package main

import (
	"fmt"
	pb "grpc/proto/addressbook"
)


func main() {
	fmt.Println("Hello from ./grpc")

	person := &pb.Person{
		Id: 1234,
		Name: "John Doe",
		Email: "jdoe@example.com",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "555-4321", Type: pb.PhoneType_PHONE_TYPE_HOME},
		},
	}

	fmt.Println("Person:", person)
}