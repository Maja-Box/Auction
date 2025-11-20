package main

import(
	proto "Auction/grpc"
	"context"
	"net"
	"log"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	counter int
	listener net.Listener
}

func main(){
	
	ports := []string{
		":5050",
		":5051",
		":5052",
	} 

	for _, port := range ports{
		server := &Auction_service{
			bids: make(map[int]int),
			highest: 0,
			counter: 0,
		}
		go server.start_server(port,ports)
	}

	/*go server.start_server(":5050",ports)
	go server.start_server(":5051",ports)
	go server.start_server(":5052",ports)*/

	select{}
}

func (server *Auction_service) start_server( numberPort string,  ports []string){
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
		log.Println("Did not work 2")
	}

}

func (server *Auction_service) Bid(ctx context.Context, in *proto.BidIn) (*proto.BidOut, error){
	if server.counter == 1{
		server.counter = 0
		log.Printf("shutting down server")
		server.listener.Close()
	}
	
	server.counter ++
	amount := in.Amount
	clientId := in.ClientId
	if clientId == 0{
		clientId = int64(len(server.bids)) + 1
		server.bids[int(clientId)] = 0
	}
	
	if server.highest < int(amount) {
		server.highest = int(amount)
		server.bids[int(clientId)] = int(amount)
	}else{
		return &proto.BidOut{Ack: "Client " + strconv.FormatInt(clientId, 10) + " Your bid has been rejected", BidderId: clientId}, nil
	}

	for value, i := range server.peers{
		_, err := i.Update(ctx, in)
		if err != nil{
			log.Printf("Could not update the info to server:", value)
		}
	}

	return &proto.BidOut{Ack: "Client " + strconv.FormatInt(clientId, 10) + " Your bid has been accepted", BidderId: clientId}, nil
}

func (server *Auction_service) Update(ctx context.Context, in *proto.BidIn) (*proto.Empty, error){
	
	amount := in.Amount
	clientId := in.ClientId

	server.highest = int(amount)
	server.bids[int(clientId)] = int(amount)

	return &proto.Empty{}, nil
}

func (server *Auction_service) Result(ctx context.Context, in *proto.Empty) (*proto.ResultSend, error){
	for client, price := range server.bids {
		if price == server.highest{
			return &proto.ResultSend{Message: ("the highest bidder is client: " + strconv.FormatInt(int64(client), 10) + 
												"with a bid of: " + strconv.FormatInt(int64(price), 10))}, nil
		}
	}

	return &proto.ResultSend{Message: ("the highest bid is: " + strconv.FormatInt(int64(server.highest), 10))}, nil
}