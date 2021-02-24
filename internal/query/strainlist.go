package query

import (
	"context"
	"image/color"
	"log"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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

type PlotData struct {
	QueryCount       float64
	UnavailableCount float64
}

func GetStrainList(config *Config) (*Analytics, error) {
	a := &Analytics{}
	ctx := context.Background()
	start := time.Now()
	avail := 0
	unavail := 0
	qc := 0
	uarr := []*PlotData{}
	var cursor int64
	for v := 10; v <= config.Quantity; v += 10 {
		uc := 0
		start := time.Now()
		qc++
		strains, err := config.StockClient.ListStrains(ctx, &stock.StockParameters{Cursor: cursor, Limit: 10, Filter: config.Filter})
		if err != nil {
			return a, err
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
				uc++
			} else {
				avail++
			}
		}
		elapsed := time.Since(start)
		log.Printf("query took %s", elapsed)
		log.Printf("%d available, %d unavailable", ac, uc)
		uarr = append(uarr, &PlotData{
			QueryCount:       float64(qc),
			UnavailableCount: float64(uc),
		})
		cursor = strains.Meta.NextCursor
	}
	if err := savePlot(uarr); err != nil {
		return a, err
	}
	elapsed := time.Since(start)
	return &Analytics{
		Elapsed:     elapsed,
		Available:   avail,
		Unavailable: unavail,
	}, nil
}

func savePlot(d []*PlotData) error {
	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = "Unavailable Strains Per Query"
	p.X.Label.Text = "Query #"
	p.Y.Label.Text = "Unavailable Strains"
	pts := make(plotter.XYs, len(d))
	for i, v := range d {
		pts[i].X = v.QueryCount
		pts[i].Y = v.UnavailableCount
	}

	s, err := plotter.NewScatter(pts)
	if err != nil {
		return err
	}
	s.Color = color.RGBA{R: 255, B: 128, A: 255}

	l, err := plotter.NewLine(pts)
	if err != nil {
		return err
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	p.Add(s, l)
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		return err
	}
	return nil
}
