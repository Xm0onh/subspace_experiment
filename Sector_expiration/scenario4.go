package main

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	totalPieces4 = 3000000
)

func Scenario4() {
	minGrowth := 100
	multiplier := 2
	pieces := make([]int, totalPieces4)
	farmers := make([]int, 3*farmerCount)

	for i := 0; i < 3*farmerCount; i++ {
		currentHeight := rand.Intn(3000000)
		latestHeight := currentHeight
		for currentHeight < totalPieces4 {
			maxGrowth := multiplier * currentHeight
			randomGrowth := rand.Intn(int(math.Abs(float64(maxGrowth-minGrowth+1)))) + minGrowth
			currentHeight += randomGrowth
			if currentHeight < totalPieces4 {
				latestHeight = currentHeight
			}
		}
		farmers[i] = latestHeight
		for j := 0; j < piecesPerFarmer; j++ {
			for {
				selectedPiece := hashValue(i, j) % uint32(farmers[i])

				if selectedPiece > totalPieces4 {
					continue
				}
				pieces[selectedPiece]++
				break
			}
		}
	}

	// 	// Farmer i has to select piecesPerFarmer pieces out of the first latestHeight pieces

	// Check for missing pieces
	missingPieces := 0
	for _, count := range pieces {
		if count == 0 {
			missingPieces++

		}
	}
	// Calculate and display the missing piece information
	fractionMissing := float64(missingPieces) / float64(totalPieces4)
	fmt.Println("Scenario 4")
	fmt.Printf("Number of Missing Pieces: %d\n", missingPieces)
	fmt.Printf("Fraction of Missing Pieces: %.4f\n", fractionMissing)

}
