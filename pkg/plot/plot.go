package plot

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/prometheus/common/model"
	stdfnt "golang.org/x/image/font"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/palette/brewer"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Only show important part of metric name
var labelText = regexp.MustCompile(`(.*)`)

// Plot creates a plot from metric data and saves it to a temporary file.
// It's the callers' responsibility to remove the returned file when no longer needed.
func Plot(metrics model.Matrix, title, format string) (io.WriterTo, error) {
	p := plot.New()
	p.Title.Text = title
	p.Title.TextStyle.Font = font.From(font.Font{
		Typeface: "Liberation",
		Variant:  "Mono",
		Style:    stdfnt.StyleItalic,
		Weight:   stdfnt.WeightBold,
	}, 0.35*vg.Centimeter)
	p.Title.Padding = 2 * vg.Centimeter
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04:05"}
	normalFont := font.From(font.Font{
		Typeface: "Liberation",
		Variant:  "Mono",
	}, 3*vg.Millimeter)
	p.X.Tick.Label.Font = normalFont
	p.Y.Tick.Label.Font = normalFont
	p.Legend.TextStyle.Font = normalFont
	p.Legend.Top = true
	p.Legend.YOffs = 15 * vg.Millimeter

	// Color palette for drawing lines
	paletteSize := 8
	palette, err := brewer.GetPalette(brewer.TypeAny, "Dark2", paletteSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get color palette: %v", err)
	}
	colors := palette.Colors()

	for s, sample := range metrics {
		data := make(plotter.XYs, len(sample.Values))
		for i, v := range sample.Values {
			data[i].X = float64(v.Timestamp.Unix())
			f, err := strconv.ParseFloat(v.Value.String(), 64)
			if err != nil {
				return nil, fmt.Errorf("sample value not float: %s", v.Value.String())
			}
			data[i].Y = f
		}

		l, err := plotter.NewLine(data)
		if err != nil {
			return nil, fmt.Errorf("failed to create line: %v", err)
		}
		l.LineStyle.Width = vg.Points(1)
		l.LineStyle.Color = colors[s%paletteSize]

		p.Add(l)
		if len(metrics) > 1 {
			m := labelText.FindStringSubmatch(sample.Metric.String())
			if m != nil {
				p.Legend.Add(m[1], l)
			}
		}
	}

	// Draw plot in canvas with margin
	margin := 6 * vg.Millimeter
	width := 24 * vg.Centimeter
	height := 10 * vg.Centimeter
	c, err1 := draw.NewFormattedCanvas(width, height, format)
	if err1 != nil {
		return nil, fmt.Errorf("failed to create canvas: %v", err1)
	}
	p.Draw(draw.Crop(draw.New(c), margin, -margin, margin, -margin))

	return c, nil
}

// WriteToFile plots the metric and write to file
func WriteToFile(metrics model.Matrix, title string, format string, name string) error {
	w, err := Plot(metrics, title, format)
	if err != nil {
		return err
	}

	f, err1 := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0666)
	if err1 != nil {
		return err1
	}

	if _, err := w.WriteTo(f); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
