package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

const (
	totalPieces2 = 2000000
)

func Scenario1() {
	minGrowth := 100
	multiplier := 2
	pieces := make([]int, totalPieces2)
	farmers := make([]int, 2*farmerCount)
	// checkPoints := make([]int, 11)
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(1)

	go func() {
		defer wg.Done()
		for i := 0; i < farmerCount; i++ {
			currentHeight := rand.Intn(1000000)
			latestHeight := currentHeight
			for currentHeight < totalPieces {
				maxGrowth := multiplier * currentHeight
				randomGrowth := rand.Intn(int(math.Abs(float64(maxGrowth-minGrowth+1)))) + minGrowth
				currentHeight += randomGrowth
				if currentHeight < totalPieces2/2 {
					latestHeight = currentHeight
				}
			}
			farmers[i] = latestHeight
			for j := 0; j < piecesPerFarmer; j++ {
				for {
					selectedPiece := hashValue(i, j) % uint32(latestHeight)
					if selectedPiece > totalPieces2 {
						continue
					}
					mu.Lock()
					pieces[selectedPiece]++
					mu.Unlock()
					break
				}
			}

		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := farmerCount; i < 2*farmerCount; i++ {

			currentHeight := rand.Intn(2000000) - 1000000
			latestHeight := currentHeight
			for currentHeight < totalPieces {
				maxGrowth := multiplier * currentHeight
				randomGrowth := rand.Intn(int(math.Abs(float64(maxGrowth-minGrowth+1)))) + minGrowth
				currentHeight += randomGrowth
				if currentHeight < totalPieces2/2 {
					latestHeight = currentHeight
				}
			}
			farmers[i] = latestHeight + 1000000
			for j := 0; j < piecesPerFarmer; j++ {
				for {
					selectedPiece := hashValue(i, j) % uint32(farmers[i])

					if selectedPiece > totalPieces2 {
						continue
					}
					mu.Lock()
					pieces[selectedPiece]++
					mu.Unlock()
					break
				}
			}
		}
	}()

	wg.Wait()

	// 	// Farmer i has to select piecesPerFarmer pieces out of the first latestHeight pieces

	// Check for missing pieces
	missingPieces := 0

	for _, count := range pieces {

		if count == 0 {
			missingPieces++

		}

	}

	// Calculate and display the missing piece information
	fractionMissing := float64(missingPieces) / float64(totalPieces2)
	fmt.Println("Scenario 1")
	fmt.Printf("Number of Missing Pieces: %d\n", missingPieces)
	fmt.Printf("Fraction of Missing Pieces: %.4f\n", fractionMissing)

}
