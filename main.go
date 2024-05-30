package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/vincy-william-alida/sigtest/testdata"
	"gonum.org/v1/gonum/stat/distuv"
)

func main() {

	var dimensions []int

	crossMap := make(map[string]float64)
	crossMapPropotion := make(map[string]float64)
	crossMapZscores := make(map[string]float64)
	crossMapIsSignificant := make(map[string]bool)

	//Read response cross tab
	err := json.Unmarshal([]byte(testdata.JsonDataResponded), &crossMap)

	if err != nil {
		fmt.Println("Error:", err)
	}

	//Read dimensions
	err = json.Unmarshal([]byte(testdata.JsonDataDimensions), &dimensions)

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(crossMap)
	rowTotal := make([]float64, dimensions[0])
	columnTotal := make([]float64, dimensions[1])

	calculateCrossTabRowTotal(crossMap, rowTotal)

	fmt.Println("RowTotal:")
	fmt.Println(rowTotal)

	calculateCrossTabColumnTotal(crossMap, columnTotal)

	fmt.Println("ColumnTotal:")
	fmt.Println(columnTotal)

	//Calculate propotion based on columnTotal
	calculatePropotionFromCrossTab(crossMap, crossMapPropotion, columnTotal)

	fmt.Println(crossMapPropotion)

	//Calculate z-test
	calculateZtest(crossMapPropotion, crossMapZscores, columnTotal, dimensions)

	fmt.Println(crossMapZscores)

	calculateSignificance(crossMapZscores, crossMapIsSignificant)

	fmt.Println(crossMapIsSignificant)
}

func calculatePropotionFromCrossTab(crossMap map[string]float64, crossMapPropotion map[string]float64, columnTotal []float64) {
	//get Propotion
	for key, value := range crossMap {
		intKey := strings.Split(key, "-")
		i, _ := strconv.Atoi(intKey[1])
		crossMapPropotion[key] = value / columnTotal[i]
	}

}

func calculateCrossTabRowTotal(crossMap map[string]float64, rowTotal []float64) {
	for key, value := range crossMap {
		intKey := strings.Split(key, "-")
		i, _ := strconv.Atoi(intKey[0])
		rowTotal[i] = rowTotal[i] + value
	}
}

func calculateCrossTabColumnTotal(crossMap map[string]float64, columnTotal []float64) {
	for key, value := range crossMap {
		intKey := strings.Split(key, "-")
		i, _ := strconv.Atoi(intKey[1])
		columnTotal[i] = columnTotal[i] + value
	}
}

func calculateZtest(crossMapPropotion map[string]float64, crossMapZScores map[string]float64, columnTotal []float64, dimensions []int) {

	for i := 0; i < dimensions[0]; i++ {
		for j := 0; j < dimensions[0]-1; j++ {
			for k := j; k < dimensions[1]-1; k++ {
				fmt.Printf("\nComparing %d-%d and %d-%d    ", i, j, i, k+1)
				z := calculateZscore(crossMapPropotion[strconv.Itoa(i)+"-"+strconv.Itoa(j)], crossMapPropotion[strconv.Itoa(i)+"-"+strconv.Itoa(k+1)], columnTotal[j], columnTotal[k+1])
				crossMapZScores[fmt.Sprintf("%d-%d,%d-%d", i, j, i, k+1)] = z
				fmt.Println(z)
			}
		}
	}
}

func calculateZscore(p1Cap float64, p2Cap float64, n1 float64, n2 float64) float64 {
	pCap := (p1Cap*n1 + p2Cap*n2) / (n1 + n2)
	qCap := 1 - pCap
	z := (p1Cap - p2Cap) / math.Sqrt((pCap * qCap * (1/n1 + 1/n2)))
	return z
}

func zScoreToPValue(z float64) float64 {
	// Create a standard normal distribution (mean=0, stddev=1)
	dist := distuv.Normal{Mu: 0, Sigma: 1}

	// Calculate the cumulative probability up to the z-score
	cdf := math.Min(dist.CDF(z), 1-dist.CDF(z))

	// For two-tailed test, p-value is 2 times the min of CDF and 1-CDF
	// pValue := 2 * (1 - cdf)
	// if cdf < 0.5 {
	// 	pValue = 2 * cdf
	// }
	pValue := 2 * cdf

	fmt.Println("pValue: ", pValue)

	return pValue
}

func calculateSignificance(crossMapZScores map[string]float64, crossMapIsSignificant map[string]bool) {
	for key, value := range crossMapZScores {
		pValue := zScoreToPValue(value)
		fmt.Println(key)
		fmt.Println(pValue)
		if pValue <= (1 - 0.95) {
			crossMapIsSignificant[key] = true
		} else {
			crossMapIsSignificant[key] = false
		}
	}

}
