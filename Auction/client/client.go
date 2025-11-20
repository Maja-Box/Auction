package main

import (
	proto "Auction/gRPC"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var clientId int = 0

var counter int

var client proto.AuctionClient

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working client 1")
	}
	log.Println("Client has connected to server")
	if err != nil {
		log.Fatalf(err.Error())
	}

	client = proto.NewAuctionClient(conn)
	clientId = 0

	bid()
	query(clientId)
}

func query(clientId int) {

	for {
		fmt.Println("If you want to check the state of the auction, type check. If you want to bid, type bid.")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		Text := scanner.Text()

		if Text == "check" {
			checkFunc()
		} else if Text == "bid" {
			bid()
		}
	}

}

func checkFunc(){
	check, err := client.Result(context.Background(), &proto.Empty{})
	if err != nil{
		log.Println("error while checking state of auction")
		connect()
		checkFunc()
	}
	log.Println(check.Message)
}

func bid() {
	bidAmount := 0
	fmt.Println("Enter your bid")
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		bid := scanner.Text()

		if _, err := strconv.ParseInt(bid, 10, 64); err != nil {
			fmt.Printf("Please input a valid number")
		} else {
			bidAmount, _ = strconv.Atoi(bid)
			break
		}
	}
	log.Println(clientId)
	send, err := client.Bid(context.Background(),
		&proto.BidIn{
			Amount:   int64(bidAmount),
			ClientId: int64(clientId),
		})
		if err != nil{
			log.Println("error while sending client bid")
			connect()
			bid()
		}
	if clientId == 0{
		clientId = int(send.BidderId)
		log.Println(send.BidderId)
	}
	log.Println(send.Ack)
}

func connect(){
	counter++
	if counter == 1{
		conn, err := grpc.NewClient("localhost:5051", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
			log.Fatalf("Not working client 1")
		}
		log.Println("Client has connected to server")
		if err != nil {
			log.Fatalf(err.Error())
		}
		client = proto.NewAuctionClient(conn)
		log.Println("i connected to server 5051")
	}else{
		conn, err := grpc.NewClient("localhost:5052", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
			log.Fatalf("Not working client 1")
		}
		log.Println("Client has connected to server")
		if err != nil {
			log.Fatalf(err.Error())
		}
		client = proto.NewAuctionClient(conn)
		log.Println("i connected to server 5052")
	}
}
