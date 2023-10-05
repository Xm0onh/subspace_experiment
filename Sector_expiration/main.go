package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

const (
	totalPieces     = 1000000
	farmerCount     = 10000
	piecesPerFarmer = 1000
)

func hashValue(farmerID int, piecePosition int) uint32 {
	data := fmt.Sprintf("%d%d", farmerID, piecePosition)
	hasher := sha256.New()
	hasher.Write([]byte(data))
	hashBytes := hasher.Sum(nil)
	return binary.BigEndian.Uint32(hashBytes[4:8])
}

func main() {

	// expOne()
	expTwo()
}

// func expTwo() {
// 	n0 := 10000
// 	minGrowth := 1000
// 	var multiplier = 3

// 	counter2 := 0

// 	for n0 <= 2*totalPieces {
// 		maxGrowth := multiplier * n0
// 		// Randomly select a value between minGrowth and maxGrowth
// 		randomGrowth := rand.Intn(maxGrowth-minGrowth+1) + minGrowth
// 		fracCal(n0, counter2)
// 		n0 += randomGrowth // Increase n0 by the random growth value
// 		counter2++
// 		fmt.Println("n0 ->", n0)
// 	}
// }

func fracCal(n int, i int) {
	// size := make([]bool, n)
	pieces := make([]int, n)
	count := 0
	replFac := piecesPerFarmer * farmerCount / n
	fmt.Println("Replication Factor ->", replFac)
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < farmerCount; i++ {
		for j := 0; j < piecesPerFarmer; j++ {
			selectedPiece := hashValue(i, j) % uint32(n)
			pieces[selectedPiece]++
		}
	}

	missingPieces := 0
	for _, count := range pieces {
		if count == 0 {
			missingPieces++
		}
	}

	fractionMissing := float64(missingPieces) / float64(n)
	fmt.Printf("Number of Missing Pieces: %d\n", count)
	fmt.Printf("i: %d\n", i)
	fmt.Printf("Fraction of Missing Pieces: %.4f\n", fractionMissing)
	fmt.Printf("====================================\n")
}

func expOne() {
	// Step 1: Initialization
	pieces := make([]int, totalPieces)

	// Step 2: Distribution Simulation
	for farmerID := 0; farmerID < farmerCount; farmerID++ {
		for piecePos := 0; piecePos < piecesPerFarmer; piecePos++ {
			selectedPiece := hashValue(farmerID, piecePos) % totalPieces
			pieces[selectedPiece]++
		}
	}

	// Step 3: Check for Missing Pieces
	missingPieces := 0
	for _, count := range pieces {
		if count == 0 {
			missingPieces++
		}
	}

	// Step 4: Calculate Fraction of Missing Pieces
	fractionMissing := float64(missingPieces) / float64(totalPieces)
	fmt.Printf("Number of Missing Pieces: %d\n", missingPieces)
	fmt.Printf("Fraction of Missing Pieces: %.4f\n", fractionMissing)
}

func expTwo() {
	initialHeight := 1000
	minGrowth := 100
	multiplier := 4

	// Create a slice to keep track of the pieces across all farmers
	pieces := make([]int, totalPieces)

	// Open file for writing. If file doesn't exist, create it, or append to the file
	file, err := os.OpenFile("latestHeight10000.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// For each farmer, calculate the latest height before 1,000,000 at which they perform a piece selection
	for i := 0; i < farmerCount; i++ {
		currentHeight := initialHeight
		latestHeight := initialHeight
		for currentHeight < totalPieces {
			maxGrowth := multiplier * currentHeight
			randomGrowth := rand.Intn(maxGrowth-minGrowth+1) + minGrowth
			currentHeight += randomGrowth
			if currentHeight < totalPieces {
				latestHeight = currentHeight
			}
		}

		// Write the latestHeight for the current farmer to the file
		_, err := file.WriteString("Farmer#" + strconv.Itoa(i) + "--" + strconv.Itoa(latestHeight) + "\n")
		if err != nil {
			panic(err)
		}

		// Farmer i has to select piecesPerFarmer pieces out of the first latestHeight pieces
		for j := 0; j < piecesPerFarmer; j++ {
			selectedPiece := hashValue(i, j) % uint32(latestHeight)
			pieces[selectedPiece]++
		}
	}

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
