package query

import (
	"context"
	"log"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
)

func GetAnnoList(config *Config) (*Analytics, error) {
	ctx := context.Background()
	start := time.Now()
	avail := 0
	unavail := 0
	for v := 10; v <= config.Quantity; v += 10 {
		start := time.Now()
		inv, err := config.AnnoClient.ListAnnotations(ctx, &annotation.ListParameters{Cursor: 0, Limit: 10, Filter: config.Filter})
		if err != nil {
			return &Analytics{}, err
		}
		for _, n := range inv.Data {
			_, err := config.StockClient.GetStrain(
				ctx,
				&stock.StockId{
					Id: n.Attributes.EntryId,
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
