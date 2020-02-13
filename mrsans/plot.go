package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

const PLOT_WIDTH = 2048
const PLOT_HEIGHT = 2048
const CONTENT_HEIGHT_PERCENT = 0.75
const CONTENT_WIDTH_PERCENT = 0.85
const PLOT_FILL_COLOR = "2196F3"

func PositionSeries(id int, total int, height float64) float64 {
	content := CONTENT_HEIGHT_PERCENT
	margin := 1 - content
	all_percent := content*float64(total) + margin*(float64(total)+1)
	return height / all_percent * (content*float64(id) + margin*(float64(id)+1))
}

func HeightSeries(total int, height float64) float64 {
	content := CONTENT_HEIGHT_PERCENT
	margin := 1 - content
	all_percent := content*float64(total) + margin*(float64(total)+1)
	return height / all_percent * content
}

func RangeOfSeries(series *[]model.SamplePair) (float64, float64) {
	first := true
	min := 0.0
	max := 1.0
	for _, data := range *series {
		val := float64(data.Value)
		if !math.IsNaN(val) {
			if first {
				min = val
				max = val
				first = false
			}
			min = math.Min(min, val)
			max = math.Max(max, val)
		}
	}
	return min, max
}

func PercentageOf(data float64, min float64, max float64, margin float64) float64 {
	if math.IsNaN(data) {
		return 0
	}
	margin = (max - min) * margin
	max += margin
	min -= margin
	return (data - min) / (max - min)
}

func PositionXOffset(width float64) float64 {
	return width * (1 - CONTENT_WIDTH_PERCENT) / 2
}

func SenseConvertString(x float64) string {
	if x <= 1000 {
		return fmt.Sprintf("%.2f", x)
	} else {
		return fmt.Sprintf("%.0f", x)
	}
}

func PlotSeries(chunkSize time.Duration, chunkOffset time.Duration, series *[]model.SamplePair, ctx *gg.Context, y float64, height float64) {
	width := float64(ctx.Width()) / float64(len(*series)) * CONTENT_WIDTH_PERCENT
	x_offset := PositionXOffset(float64(ctx.Width()))
	x_offset_right := float64(ctx.Width()) - x_offset
	min, max := RangeOfSeries(series)
	margin := 0.05

	// Text
	ctx.SetLineWidth(1)
	ctx.SetDash(16, 16)
	{
		percent := PercentageOf(min, min, max, margin)
		ypos := height*(1-percent) + y
		ctx.SetHexColor("000000")
		ctx.DrawStringAnchored(SenseConvertString(min), x_offset-10, ypos, 1, 0.5)
		ctx.SetHexColor("aaaaaa")
		ctx.DrawLine(x_offset, ypos, x_offset_right, ypos)
		ctx.Stroke()
	}
	{
		percent := PercentageOf(max, min, max, margin)
		ypos := height*(1-percent) + y
		ctx.SetHexColor("000000")
		ctx.DrawStringAnchored(SenseConvertString(max), x_offset-10, ypos, 1, 0.5)
		ctx.SetHexColor("aaaaaa")
		ctx.DrawLine(x_offset, ypos, x_offset_right, ypos)
		ctx.Stroke()
	}

	// Chunk
	firstChunk := true
	lstChunk := time.Now()
	for idx, data := range *series {
		xpos := float64(idx)*width + x_offset
		currentChunk := data.Timestamp.Time().Add(-chunkOffset).Truncate(chunkSize).Add(chunkOffset)
		if !currentChunk.Equal(lstChunk) || firstChunk {
			if !firstChunk {
				ctx.DrawLine(xpos, y, xpos, y+height)
				ctx.Stroke()
			}
			firstChunk = false
			lstChunk = currentChunk
		}
	}

	ctx.SetDash(1)

	// Fill
	ctx.SetRGBA255(33, 150, 243, int(math.Floor(255*0.3)))
	ctx.MoveTo(x_offset, y+height)
	for idx, data := range *series {
		percent := PercentageOf(float64(data.Value), min, max, margin)
		if percent == 0 {
			continue
		}
		ypos := height*(1-percent) + y
		xpos := float64(idx)*width + x_offset
		ctx.LineTo(xpos, ypos)
	}
	ctx.LineTo(width*float64(len(*series))+x_offset, y+height)
	ctx.LineTo(x_offset, y+height)
	ctx.Fill()

	first := true
	lstX, lstY := float64(0), float64(0)

	// Draw Line
	ctx.SetHexColor(PLOT_FILL_COLOR)
	for idx, data := range *series {
		percent := PercentageOf(float64(data.Value), min, max, margin)
		ypos := height*(1-percent) + y
		xpos := float64(idx)*width + x_offset
		if percent == 0 {
			continue
		}
		if first {
			first = false
		} else {
			ctx.DrawLine(lstX, lstY, xpos, ypos)
			ctx.SetLineWidth(3)
			ctx.Stroke()
		}
		lstX = xpos
		lstY = ypos
	}
}

func Plot(msg string, chunkSize time.Duration, chunkOffset time.Duration, temp *[]model.SamplePair, hum *[]model.SamplePair, pa *[]model.SamplePair, pm25 *[]model.SamplePair, pm10 *[]model.SamplePair, filename string) {
	plot_total := 5
	ctx := gg.NewContext(PLOT_WIDTH, PLOT_HEIGHT)
	ctx.SetHexColor("ffffff")
	ctx.Clear()

	// Plot Text
	ctx.SetRGB(0, 0, 0)
	x_offset := PositionXOffset(float64(ctx.Width()))
	x_offset_right := float64(ctx.Width()) - PositionXOffset(float64(ctx.Width()))
	if err := ctx.LoadFontFace(Config.plot_fontface, 40); err != nil {
		log.Fatalf("failed to load font face: %v", err)
	}
	ctx.DrawStringAnchored(fmt.Sprintf("Mr. Sans reporting at %s", time.Now().Format("Mon Jan 2 15:04 MST 2006")), x_offset, 50, 0, 0)
	if err := ctx.LoadFontFace(Config.plot_fontface, 30); err != nil {
		log.Fatalf("failed to load font face: %v", err)
	}

	ctx.DrawString("Temperature  °C", x_offset, PositionSeries(0, plot_total, PLOT_HEIGHT))
	ctx.DrawString("Humidity  %", x_offset, PositionSeries(1, plot_total, PLOT_HEIGHT))
	ctx.DrawString("Pressure  kPa", x_offset, PositionSeries(2, plot_total, PLOT_HEIGHT))
	ctx.DrawString("PM2.5  µg/m^3", x_offset, PositionSeries(3, plot_total, PLOT_HEIGHT))
	ctx.DrawString("PM10  µg/m^3", x_offset, PositionSeries(4, plot_total, PLOT_HEIGHT))

	start_time := (*temp)[0].Timestamp
	end_time := (*temp)[len(*temp)-1].Timestamp

	ctx.DrawStringAnchored(start_time.Time().Format("Mon Jan 2 15:04"), x_offset, PLOT_HEIGHT-60, 0, 1)
	ctx.DrawStringAnchored(end_time.Time().Format("Mon Jan 2 15:04"), x_offset_right, PLOT_HEIGHT-60, 1, 1)
	ctx.DrawStringAnchored(msg, PLOT_WIDTH/2, PLOT_HEIGHT-60, 0.5, 1)

	// Plot Series
	PlotSeries(chunkSize, chunkOffset, temp, ctx, PositionSeries(0, plot_total, PLOT_HEIGHT), HeightSeries(plot_total, PLOT_HEIGHT))
	PlotSeries(chunkSize, chunkOffset, hum, ctx, PositionSeries(1, plot_total, PLOT_HEIGHT), HeightSeries(plot_total, PLOT_HEIGHT))
	PlotSeries(chunkSize, chunkOffset, pa, ctx, PositionSeries(2, plot_total, PLOT_HEIGHT), HeightSeries(plot_total, PLOT_HEIGHT))
	PlotSeries(chunkSize, chunkOffset, pm25, ctx, PositionSeries(3, plot_total, PLOT_HEIGHT), HeightSeries(plot_total, PLOT_HEIGHT))
	PlotSeries(chunkSize, chunkOffset, pm10, ctx, PositionSeries(4, plot_total, PLOT_HEIGHT), HeightSeries(plot_total, PLOT_HEIGHT))

	ctx.SavePNG(filename)
}
