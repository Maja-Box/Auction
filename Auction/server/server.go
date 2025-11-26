package main

import (
	proto "Auction/grpc"
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Auction_service struct {
	proto.UnimplementedAuctionServer
	error      chan error
	grpc       *grpc.Server
	serverPort string
	highest    int
	//first int is client id, second is highest bid that that client has made
	bids          map[int]int
	ports         []string
	peers         map[string]proto.AuctionClient //client pointing to other servers
	listener      net.Listener
	timeIsStarted bool
	auctionOver   bool
}

func main() {
	ports := []string{
		":5050",
		":5051",
		":5052",
	}

	server := &Auction_service{
		bids:          make(map[int]int),
		highest:       0,
		timeIsStarted: false,
		auctionOver:   false,
	}

	log.Println("Enter the port of the server (A number from 0 to 2)")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	Text := scanner.Text()

	if Text == "0" || Text == "1" || Text == "2" {
		port, _ := strconv.ParseInt(Text, 10, 64)
		go server.start_server(ports[port], ports)
		log.Println("Port selected: " + ports[port])
	} else {
		log.Println("Enter the correct port of the server (A number from 0 to 2)")
	}

	/*go server.start_server(":5050",ports)
	go server.start_server(":5051",ports)
	go server.start_server(":5052",ports)*/

	select {}
}

func (server *Auction_service) start_server(numberPort string, ports []string) {
	server.grpc = grpc.NewServer()
	listener, err := net.Listen("tcp", numberPort)
	if err != nil {
		log.Fatalf("Did not work 1")
	}
	server.peers = make(map[string]proto.AuctionClient)
	server.listener = listener

	log.Println("the server has started")
	server.ports = ports

	for _, value := range server.ports {
		if value != numberPort {
			connection := "localhost" + value
			conn, err := grpc.NewClient(connection, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("connection failed")
			}

			client := proto.NewAuctionClient(conn)
			server.peers[value] = client
		}
	}

	proto.RegisterAuctionServer(server.grpc, server)

	err = server.grpc.Serve(listener)
	if err != nil {
		log.Println("Could not call grpc.Serve on server")
	}

}

func (server *Auction_service) Bid(ctx context.Context, in *proto.BidIn) (*proto.BidOut, error) {

	if !server.timeIsStarted {
		server.timeIsStarted = true
		go server.Timer()
	}

	if server.auctionOver {
		log.Println("Auction is over")
		return &proto.BidOut{Ack: "The auction is over, your bid has been rejected. If you want to check the winner of the auction, type check."}, nil

	}

	amount := in.Amount
	clientId := in.ClientId
	if clientId == 0 {
		clientId = int64(len(server.bids)) + 1
		server.bids[int(clientId)] = 0
	}

	if server.highest < int(amount) {
		server.highest = int(amount)
		server.bids[int(clientId)] = int(amount)
	} else {
		log.Println("Client " + strconv.FormatInt(clientId, 10) + " Your bid has been rejected")
		return &proto.BidOut{Ack: "Client " + strconv.FormatInt(clientId, 10) + " Your bid has been rejected ", BidderId: clientId}, nil
	}

	for value, i := range server.peers {
		_, err := i.Update(ctx, in)
		if err != nil {

			log.Printf("Could not update the info to server:", value)
			server.UpdateServer(context.Background(), &proto.Crash{Port : value,})
		}
	}

	log.Println("Client " + strconv.FormatInt(clientId, 10) + " Your bid has been accepted")
	return &proto.BidOut{Ack: "Client " + strconv.FormatInt(clientId, 10) + " Your bid has been accepted ", BidderId: clientId}, nil
}

func (server *Auction_service) Update(ctx context.Context, in *proto.BidIn) (*proto.Empty, error) {

	amount := in.Amount
	clientId := in.ClientId

	if !server.timeIsStarted {
		server.timeIsStarted = true
		go server.Timer()
	}

	server.highest = int(amount)
	server.bids[int(clientId)] = int(amount)

	return &proto.Empty{}, nil
}

func (server *Auction_service) UpdateServer(ctx context.Context, in *proto.Crash) (*proto.Empty, error) {

	bad := in.Port
	for i, value := range server.ports {
		if value == bad {
			delete_index(server.ports, i)
			delete(server.peers, value)
		}
	}

	for value, i := range server.peers {
		_, err := i.ReplicateCrash(ctx, in)
		if err != nil {
			log.Printf("Could not update the info to server: ", value)

		}
	}

	return &proto.Empty{}, nil
}

func (server *Auction_service) ReplicateCrash(ctx context.Context, in *proto.Crash) (*proto.Empty, error) {
	bad := in.Port
	for i, value := range server.ports {
		if value == bad {
			delete_index(server.ports, i)
			delete(server.peers, value)
		}
	}
	return &proto.Empty{}, nil
}

func delete_index(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func (server *Auction_service) Result(ctx context.Context, in *proto.Empty) (*proto.ResultSend, error) {
	for client, price := range server.bids {
		if price == server.highest {
			if server.auctionOver {
				return &proto.ResultSend{Message: ("the winner of the auction is client: " + strconv.FormatInt(int64(client), 10) +
					" with a bid of: " + strconv.FormatInt(int64(price), 10))}, nil
			}

			return &proto.ResultSend{Message: ("the highest bidder is client: " + strconv.FormatInt(int64(client), 10) +
				" with a bid of: " + strconv.FormatInt(int64(price), 10))}, nil
		}
	}

	return &proto.ResultSend{Message: ("the highest bid is: " + strconv.FormatInt(int64(server.highest), 10))}, nil
}

func (server *Auction_service) Timer() {
	start := time.Now()

	for {
		if time.Since(start) >= (120 * time.Second) {
			server.auctionOver = true
			log.Println("Time finished")
			break
		}
	}
}
