package main

import (
	"fmt"
	"os"

	"github.com/wcharczuk/go-chart"
)

func main() {
	// Sample data as an array of maps
	mapList := []map[string]int{
		{
			"CSS":         2883,
			"Dockerfile": 27100,
			"Go":          1797643,
			"HTML":        24153,
			"Java":        14593494,
			"JavaScript":  3423720,
			"Makefile":    47721,
			"Mustache":    9534,
			"Perl":        951,
			"Shell":       147851,
		},
		{
			"Dockerfile": 3221,
			"Go":         696651,
			"Makefile":   3435,
			"Shell":      1530,
		},
	}

	// Initialize a map to store counts for each programming language
	languageCounts := make(map[string]int)

	// Sum up counts for each language across all data
	for _, item := range mapList {
		for lang, count := range item {
			// Increment the count for the current language
			languageCounts[lang] += count
		}
	}

	// Create data for the pie chart
	var values []chart.Value
	for lang, count := range languageCounts {
		values = append(values, chart.Value{Label: lang, Value: float64(count)})
	}

	// Create a pie chart
	pie := chart.PieChart{
		Width:  4096,  // Increased width
		Height: 4096,  // Increased height
		Values: values,
	}

	// Save the chart as an image file
	f, err := os.Create("pie_chart.png")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	err = pie.Render(chart.PNG, f)
	if err != nil {
		fmt.Println("Error rendering pie chart:", err)
		return
	}

	fmt.Println("Pie chart saved as 'pie_chart.png'")
}
