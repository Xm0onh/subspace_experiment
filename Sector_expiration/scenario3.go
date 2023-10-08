package main

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	totalPieces3 = 3000000
)

func Scenario3() {
	minGrowth := 100
	multiplier := 2
	pieces := make([]int, totalPieces3)
	farmers := make([]int, 3*farmerCount)

	for i := 0; i < 3*farmerCount; i++ {
		switch {
		case i < farmerCount:
			currentHeight := rand.Intn(1000000)
			latestHeight := currentHeight
			for currentHeight < totalPieces {
				maxGrowth := multiplier * currentHeight
				randomGrowth := rand.Intn(int(math.Abs(float64(maxGrowth-minGrowth+1)))) + minGrowth
				currentHeight += randomGrowth
				if currentHeight < totalPieces3/3 {
					latestHeight = currentHeight
				}
			}
			farmers[i] = latestHeight
			for j := 0; j < piecesPerFarmer; j++ {
				for {
					selectedPiece := hashValue(i, j) % uint32(latestHeight)
					if selectedPiece > totalPieces3 {
						continue
					}
					pieces[selectedPiece]++
					break
				}
			}
		case farmerCount <= i && i < 2*farmerCount:
			currentHeight := rand.Intn(2000000) - 1000000
			latestHeight := currentHeight
			for currentHeight < totalPieces {
				maxGrowth := multiplier * currentHeight
				randomGrowth := rand.Intn(int(math.Abs(float64(maxGrowth-minGrowth+1)))) + minGrowth
				currentHeight += randomGrowth
				if currentHeight < totalPieces3/3 {
					latestHeight = currentHeight
				}
			}
			farmers[i] = latestHeight + 1000000
			for j := 0; j < piecesPerFarmer; j++ {
				for {
					selectedPiece := hashValue(i, j) % uint32(latestHeight)
					if selectedPiece > totalPieces3 {
						continue
					}
					pieces[selectedPiece]++
					break
				}
			}
		case 2*farmerCount <= i && i < 3*farmerCount:

			currentHeight := rand.Intn(3000000) - 2000000
			latestHeight := currentHeight
			for currentHeight < totalPieces {
				maxGrowth := multiplier * currentHeight
				randomGrowth := rand.Intn(int(math.Abs(float64(maxGrowth-minGrowth+1)))) + minGrowth
				currentHeight += randomGrowth
				if currentHeight < totalPieces3/3 {
					latestHeight = currentHeight
				}
			}
			farmers[i] = latestHeight + 2000000
			for j := 0; j < piecesPerFarmer; j++ {
				for {
					selectedPiece := hashValue(i, j) % uint32(latestHeight)
					if selectedPiece > totalPieces3 {
						continue
					}
					pieces[selectedPiece]++
					break
				}
			}
		}
	}
	// fmt.Println(farmers[5000])
	// fmt.Println(farmers[15000])
	// fmt.Println(farmers[25000])

	// 	// Farmer i has to select piecesPerFarmer pieces out of the first latestHeight pieces

	// Check for missing pieces
	missingPieces := 0
	for _, count := range pieces {
		if count == 0 {
			missingPieces++

		}
	}

	// Calculate and display the missing piece information
	fractionMissing := float64(missingPieces) / float64(totalPieces)
	fmt.Printf("Number of Missing Pieces: %d\n", missingPieces)
	fmt.Printf("Fraction of Missing Pieces: %.4f\n", fractionMissing)
}
