package main

import(
	proto "Auction/grpc"
	"context"
	"net"
	"log"
	"strconv"

	"google.golang.org/grpc"
)

type Auction_service struct{
	proto.UnimplementedAuctionServer
	error       chan error
	grpc        *grpc.Server
	serverId	int 
	highest int
	//first int is client id, second is highest bid that that client has made
	bids    map[int]int
	
}

func main(){
	server := &Auction_service{
		bids: make(map[int]int),
		highest: 0,
	}
	server.start_server(1)
	server.start_server(2)
	server.start_server(3)
}

func (server *Auction_service) start_server(serId int){
	server.grpc = grpc.NewServer()
	listener, err := net.Listen("tcp", ":5050")

	if err != nil {
		log.Fatalf("Did not work 1")
	}

	log.Println("the server has started")

	proto.RegisterAuctionServer(server.grpc, server)

	err = server.grpc.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work 2")
	}

	server.serverId = serId
	
}

func (server *Auction_service) Bid(ctx context.Context, in proto.BidIn) (*proto.BidOut, error){
	amount := in.Amount
	clientId := in.ClientId
	if clientId == 0{
		clientId = int64(len(server.bids)) + 1
	}
	
	if server.highest < int(amount) {
		server.highest = int(amount)
		server.bids[int(amount)] = int(clientId)
	}else{
		server.bids[int(amount)] = 0
		return &proto.BidOut{Ack: "Client " + strconv.FormatInt(clientId, 10) + "Your bid has been rejected", BidderId: clientId}, nil
	}

	return &proto.BidOut{Ack: "Client " + strconv.FormatInt(clientId, 10) + "Your bid has been accepted", BidderId: clientId}, nil
}

func (server *Auction_service) Result(ctx context.Context, in proto.Empty) (*proto.ResultSend, error){
	for client, price := range server.bids {
		if price == server.highest{
			return &proto.ResultSend{Message: ("the highest bidder is client: " + strconv.FormatInt(int64(client), 10) + 
												"with a bid of: " + strconv.FormatInt(int64(price), 10))}, nil
		}
	}

	return &proto.ResultSend{Message: ("the highest bid is: " + strconv.FormatInt(int64(server.highest), 10))}, nil
}