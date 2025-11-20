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
	serverPort	string 
	highest int
	//first int is client id, second is highest bid that that client has made
	bids    map[int]int
	ports	[]string
	peers   map[string]proto.AuctionClient //client pointing to other servers
	
}

func main(){
	server := &Auction_service{
		bids: make(map[int]int),
		highest: 0,
	}
	ports := []string{
		":5050",
		":5051",
		":5052",
	} 

	

	go server.start_server(":5050",ports)
	go server.start_server(":5051",ports)
	go server.start_server(":5052",ports)

	select{}
}

func (server *Auction_service) start_server(string numberPort,[]string ports){
	server.grpc = grpc.NewServer()
	listener, err := net.Listen("tcp", numberPort)
	server.peers = make(map[string]proto.AuctionClient)


	if err != nil {
		log.Fatalf("Did not work 1")
	}

	log.Println("the server has started")

	for _, value := range server.ports {
		if(value != numberPort){
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
		server.bids[int(clientId)] = int(amount)
	}else{
		server.bids[int(clientId)] = 0
		return &proto.BidOut{Ack: "Client " + strconv.FormatInt(clientId, 10) + "Your bid has been rejected", BidderId: clientId}, nil
	}

	for _, i := range server.peers{
		
	}

	return &proto.BidOut{Ack: "Client " + strconv.FormatInt(clientId, 10) + "Your bid has been accepted", BidderId: clientId}, nil
}

func (server *Auction_service) Update(ctx context.context, in proto.UpdateSend) (*proto.Empty, error){
	amount := in.Amount
	clientId := in.ClientId

	server.highest = int(amount)
	server.bids[int(clientId)] = int(amount)

	return &proto.Empty{}, nil
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