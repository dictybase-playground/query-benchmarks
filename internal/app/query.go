package app

import (
	"context"
	"fmt"
	"log"
	"time"

	anno "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func TimeQuery(c *cli.Context) error {
	// initialize anno client
	ahost := c.String("annotation-grpc-host")
	aport := c.String("annotation-grpc-port")
	aconn, err := grpc.Dial(fmt.Sprintf("%s:%s", ahost, aport), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer aconn.Close()
	ann := anno.NewTaggedAnnotationServiceClient(aconn)
	// initialize stock client
	shost := c.String("stock-grpc-host")
	sport := c.String("stock-grpc-host")
	sconn, err := grpc.Dial(fmt.Sprintf("%s:%s", shost, sport), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer sconn.Close()
	st := stock.NewStockServiceClient(sconn)
	ctx := context.Background()

	// to get non-gwdi strains "name=~GWDI"
	start := time.Now()
	avail := 0
	unavail := 0
	for v := 10; v <= 20; v += 10 {
		start := time.Now()
		strains, err := st.ListStrains(ctx, &stock.StockParameters{Cursor: 0, Limit: 10, Filter: "name!~GWDI;label!=AX4"})
		if err != nil {
			panic(err)
		}
		for _, n := range strains.Data {
			_, err := ann.GetEntryAnnotation(
				ctx,
				&anno.EntryAnnotationRequest{
					Tag:      "strain_inventory",
					Ontology: "strain_inventory",
					EntryId:  n.Id,
				},
			)
			if err != nil {
				unavail++
			} else {
				avail++
			}
		}
		elapsed := time.Since(start)
		log.Printf("query took %s", elapsed)
	}
	elapsed := time.Since(start)
	log.Printf("full execution took %s", elapsed)
	log.Printf("total number of available strains: %d", avail)
	log.Printf("total number of unavailable strains: %d", unavail)
	return nil
}
