package main

import (
	proto "Auction/gRPC"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var id int = 0

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working client 1")
	}
	log.Println("Client has connected to server")
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	client := proto.NewAuctionClient(conn)
	clientId := 0

	bid(clientId)
	go query(clientId)
}

func query(clientId int) {

	for {
		fmt.Println("If you want to check the state of the auction, type in check. If you want to bid, type bid.")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		Text := scanner.Text()

		if Text == "check" {
			send, err := client.Result(context.Background())
		} else if Text == "bid" {
			bid(clientId)
		}
	}

}

func bid(clientId int) {
	bidAmount := 0
	fmt.Println("Enter your bid")
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		bid := scanner.Text()

		if _, err := strconv.ParseInt(bid, 10, 64); err != nil {
			fmt.Printf("Please input a valid number")
		} else {
			bidAmount = strconv.ParseInt(bid, 10, 64)
			break
		}
	}

	send, err := client.BidIn(context.Background(),
		&proto.BidIn{
			Amount:   bidAmount,
			ClientId: clientId,
		})

}
