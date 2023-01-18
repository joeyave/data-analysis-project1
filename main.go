package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/joeyave/data-analysis-project1/kmeans"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"
)

//var observations = []kmeans.Node{
//	{20.0, 20.0},
//	{21.0, 21.0},
//	{100.5, 100.5},
//	{50.1, 50.1},
//	{64.2, 64.2},
//}

func main() {

	f, err := os.Open("dataset.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var observations []kmeans.Node
	for _, record := range records {
		x, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			continue
		}
		y, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			continue
		}

		observations = append(observations, kmeans.Node{x, y})
	}

	p := plot.New()
	p.Add(plotter.NewGrid())

	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	p.X.Min = Min(observations, 1)
	p.X.Max = Max(observations, 1)

	p.Y.Min = Min(observations, 1)
	p.Y.Max = Max(observations, 1)

	// Get a list of centroids and output the values
	if success, centroids := kmeans.Train(observations, 2, 50); success {
		// Show the centroids
		fmt.Println("The centroids are")
		for _, centroid := range centroids {
			fmt.Println(centroid)
		}

		centroidToObservations := make(map[*kmeans.Node][]kmeans.Node)

		// Output the clusters
		fmt.Println("...")
		for _, observation := range observations {
			index := kmeans.Nearest(observation, centroids)
			centroidToObservations[&centroids[index]] = append(centroidToObservations[&centroids[index]], observation)
			//fmt.Println(observation, "belongs in cluster", index+1, ".")
		}

		for centroid, observations := range centroidToObservations {
			colr := color.RGBA{
				R: uint8(rand.Intn(127)),
				G: uint8(rand.Intn(127)),
				B: uint8(rand.Intn(127)),
				A: 127,
			}

			dots := plotter.XYs{}

			for _, o := range observations {
				dots = append(dots, plotter.XY{X: o[0], Y: o[1]})
			}

			scatter, err := plotter.NewScatter(dots)
			if err != nil {
				log.Fatal(err)
			}
			scatter.GlyphStyle = draw.GlyphStyle{
				Color:  colr,
				Radius: vg.Points(3),
				Shape:  draw.CircleGlyph{},
			}
			p.Add(scatter)

			centroidDots := plotter.XYs{{X: (*centroid)[0], Y: (*centroid)[1]}}
			centroidScatter, err := plotter.NewScatter(centroidDots)
			if err != nil {
				log.Fatal(err)
			}
			centroidScatter.GlyphStyle = draw.GlyphStyle{
				Color:  colr,
				Radius: vg.Points(5),
				Shape:  draw.CrossGlyph{},
			}

			p.Add(centroidScatter)
		}
	}

	writerTo, err := p.WriterTo(360, 360, "svg")
	if err != nil {
		return
	}

	buf := new(bytes.Buffer)
	writerTo.WriteTo(buf)

	err = os.WriteFile("res.svg", buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func Max(nodes []kmeans.Node, index int) float64 {
	max := nodes[0][index]
	for i := range nodes {
		if nodes[i][index] > max {
			max = nodes[i][index]
		}
	}
	return max
}

func Min(nodes []kmeans.Node, index int) float64 {
	min := nodes[0][index]
	for i := range nodes {
		if nodes[i][index] < min {
			min = nodes[i][index]
		}
	}
	return min
}
