package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

const (
	totalPieces3 = 3000000
)

func Scenario3() {
	minGrowth := 100
	multiplier := 2
	pieces := make([]int, totalPieces3)
	farmers := make([]int, 3*farmerCount)
	checkPoints := make([]int, 11)
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
				if currentHeight < totalPieces3/3 {
					latestHeight = currentHeight
				}
			}
			farmers[i] = latestHeight + 1000000
			for j := 0; j < piecesPerFarmer; j++ {
				for {
					selectedPiece := hashValue(i, j) % uint32(farmers[i])

					if selectedPiece > totalPieces3 {
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
		for i := 2 * farmerCount; i < 3*farmerCount; i++ {

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
					selectedPiece := hashValue(i, j) % uint32(farmers[i])
					if selectedPiece > totalPieces3 {
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

	fmt.Println(farmers[5000])
	fmt.Println(farmers[15000])
	fmt.Println(farmers[25000])
	test := 0
	for i := 20000; i < 30000; i += 1 {
		if (farmers[i] < 2000000) || (farmers[i] > 3000000) {
			test++
		}
	}

	fmt.Println(test)

	// 	// Farmer i has to select piecesPerFarmer pieces out of the first latestHeight pieces

	// Check for missing pieces
	missingPieces := 0

	for index, count := range pieces {

		if count == 0 {
			missingPieces++

		}
		switch {
		case index == len(pieces)*20/30:
			checkPoints[0] = missingPieces
		case index == len(pieces)*21/30:
			checkPoints[1] = missingPieces
		case index == len(pieces)*22/30:
			checkPoints[2] = missingPieces
		case index == len(pieces)*23/30:
			checkPoints[3] = missingPieces
		case index == len(pieces)*24/30:
			checkPoints[4] = missingPieces
		case index == len(pieces)*25/30:
			checkPoints[5] = missingPieces
		case index == len(pieces)*26/30:
			checkPoints[6] = missingPieces
		case index == len(pieces)*27/30:
			checkPoints[7] = missingPieces
		case index == len(pieces)*28/30:
			checkPoints[8] = missingPieces
		case index == len(pieces)*29/30:
			checkPoints[9] = missingPieces
		case index == len(pieces)-1:
			checkPoints[10] = missingPieces
		}

	}

	// Calculate and display the missing piece information
	fractionMissing := float64(missingPieces) / float64(totalPieces3)
	fmt.Println("Scenario 3")
	fmt.Printf("Number of Missing Pieces: %d\n", missingPieces)
	fmt.Printf("Fraction of Missing Pieces: %.4f\n", fractionMissing)
	for i := 0; i < len(checkPoints); i++ {
		fmt.Printf("CheckPoint %d: %d\n", i, checkPoints[i])
	}
}
