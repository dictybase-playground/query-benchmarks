package app

import (
	"log"

	anno "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictybase-playground/query-benchmarks/internal/query"
	"github.com/urfave/cli"
)

func TimeRegQuery(c *cli.Context) error {
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
		Filter:      "name!~GWDI;label!=AX4",
	}
	l, err := query.GetStrainList(config)
	if err != nil {
		return err
	}
	log.Printf("full execution took %s", l.Elapsed)
	log.Printf("total number of available strains: %d", l.Available)
	log.Printf("total number of unavailable strains: %d", l.Unavailable)
	return nil
}
