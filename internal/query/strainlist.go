package query

import (
	"context"
	"log"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
)

type Config struct {
	StockClient stock.StockServiceClient
	AnnoClient  annotation.TaggedAnnotationServiceClient
	Quantity    int
	Filter      string
}

type Analytics struct {
	Elapsed     time.Duration
	Available   int
	Unavailable int
}

func GetStrainList(config *Config) (*Analytics, error) {
	ctx := context.Background()
	start := time.Now()
	avail := 0
	unavail := 0
	for v := 10; v <= config.Quantity; v += 10 {
		start := time.Now()
		strains, err := config.StockClient.ListStrains(ctx, &stock.StockParameters{Cursor: 0, Limit: 10, Filter: config.Filter})
		if err != nil {
			return &Analytics{}, err
		}
		for _, n := range strains.Data {
			_, err := config.AnnoClient.GetEntryAnnotation(
				ctx,
				&annotation.EntryAnnotationRequest{
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
	return &Analytics{
		Elapsed:     elapsed,
		Available:   avail,
		Unavailable: unavail,
	}, nil
}
