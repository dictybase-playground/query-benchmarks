package app

import (
	"fmt"
	"log"

	anno "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictybase-playground/query-benchmarks/internal/query"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

type Conn struct {
	aconn *grpc.ClientConn
	sconn *grpc.ClientConn
}

func TimeAnnoQuery(c *cli.Context) error {
	clients, err := initializeClients(c)
	if err != nil {
		return err
	}
	defer clients.aconn.Close()
	defer clients.sconn.Close()
	config := &query.Config{
		AnnoClient:  anno.NewTaggedAnnotationServiceClient(clients.aconn),
		StockClient: stock.NewStockServiceClient(clients.sconn),
		Quantity:    c.Int("qty"),
		Filter:      "ontology==strain_inventory;tag==strain_inventory",
	}
	l, err := query.GetAnnoList(config)
	if err != nil {
		return err
	}
	log.Printf("full execution took %s", l.Elapsed)
	log.Printf("total number of available strains: %d", l.Available)
	log.Printf("total number of unavailable strains: %d", l.Unavailable)
	return nil
}

func initializeClients(c *cli.Context) (*Conn, error) {
	cn := &Conn{}
	// initialize anno client
	ahost := c.String("annotation-grpc-host")
	aport := c.String("annotation-grpc-port")
	aconn, err := grpc.Dial(fmt.Sprintf("%s:%s", ahost, aport), grpc.WithInsecure())
	if err != nil {
		return cn, err
	}
	// initialize stock client
	shost := c.String("stock-grpc-host")
	sport := c.String("stock-grpc-port")
	sconn, err := grpc.Dial(fmt.Sprintf("%s:%s", shost, sport), grpc.WithInsecure())
	if err != nil {
		return cn, err
	}
	return &Conn{
		aconn: aconn,
		sconn: sconn,
	}, nil
}
