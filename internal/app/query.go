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

func TimeQuery(c *cli.Context) error {
	// initialize anno client
	ahost := c.String("annotation-grpc-host")
	aport := c.String("annotation-grpc-port")
	aconn, err := grpc.Dial(fmt.Sprintf("%s:%s", ahost, aport), grpc.WithInsecure())
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("cannot connect to annotation grpc microservice %s", err),
			2,
		)
	}
	defer aconn.Close()
	ann := anno.NewTaggedAnnotationServiceClient(aconn)
	// initialize stock client
	shost := c.String("stock-grpc-host")
	sport := c.String("stock-grpc-port")
	sconn, err := grpc.Dial(fmt.Sprintf("%s:%s", shost, sport), grpc.WithInsecure())
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("cannot connect to stock grpc microservice %s", err),
			2,
		)
	}
	defer sconn.Close()
	st := stock.NewStockServiceClient(sconn)
	filter := ""
	q := c.String("query")
	switch q {
	case "reginv":
		filter = "name!~GWDI;label!=AX4"
	case "gwdi":
		filter = "name=~GWDI"
	case "inv":
		filter = "ontology==strain_inventory;tag==strain_inventory"
	}
	config := &query.Config{
		StockClient: st,
		AnnoClient:  ann,
		Quantity:    c.Int("quantity"),
		Filter:      filter,
	}
	if q == "reginv" || q == "gwdi" {
		l, err := query.GetStrainList(config)
		if err != nil {
			return err
		}
		log.Printf("full execution took %s", l.Elapsed)
		log.Printf("total number of available strains: %d", l.Available)
		log.Printf("total number of unavailable strains: %d", l.Unavailable)
	}
	if q == "inv" {
		l, err := query.GetAnnoList(config)
		if err != nil {
			return err
		}
		log.Printf("full execution took %s", l.Elapsed)
		log.Printf("total number of available strains: %d", l.Available)
		log.Printf("total number of unavailable strains: %d", l.Unavailable)
	}
	return nil
}
