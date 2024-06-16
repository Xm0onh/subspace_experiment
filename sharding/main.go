package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"
)

type Bundle struct {
	Data string
}

type BundleHeader struct {
	Data string
}

type Block struct {
	ID      int
	Bundles []Bundle
}

type Farmer struct {
	ID      int
	ShardID int
	Bundles []Bundle
	Headers []BundleHeader
	Chain   []Block
	mu      sync.Mutex
}

type Shard struct {
	ID      int
	Farmers []*Farmer
	Leader  *Farmer
	mu      sync.Mutex
}

type Operator struct {
	ID        int
	Shard     *Shard
	Bandwidth int
}

type Metric struct {
	ShardID   int
	BlockID   int
	Stage     string
	TimeTaken float64
}

// Global variables for bundle sizes
var (
	bundleBodySize   = 1 * 1024 * 1024 * 8 // 1 MB in bits
	bundleHeaderSize = 10 * 1024 * 8       // 10 KB in bits
)

func chooseLeader(farmers []*Farmer, rng *rand.Rand) *Farmer {
	return farmers[rng.Intn(len(farmers))]
}

func chooseFarmer(farmers []*Farmer, round int) *Farmer {
	return farmers[round%len(farmers)]
}

func randomNetworkDelay() time.Duration {
	return time.Duration(rand.Intn(300)) * time.Millisecond
}

func timeToSendData(sizeBits, bandwidthMbps int) time.Duration {
	bandwidthBps := bandwidthMbps * 1024 * 1024
	return time.Duration(float64(sizeBits) / float64(bandwidthBps) * float64(time.Second))
}

// Function for the operator to send bundle headers to all farmers
func (op *Operator) sendBundleHeaders(header BundleHeader, allFarmers []*Farmer, wg *sync.WaitGroup, metrics *[]Metric, shardID int, blockID int) {
	defer wg.Done()

	startTime := time.Now()
	for _, farmer := range allFarmers {
		time.Sleep(timeToSendData(bundleHeaderSize, op.Bandwidth)) // Simulate sending header time based on bandwidth
		farmer.mu.Lock()
		farmer.Headers = append(farmer.Headers, header)
		farmer.mu.Unlock()
		fmt.Printf("Operator %d sent bundle header to farmer %d\n", op.ID, farmer.ID)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Operator %d finished sending bundle headers in %v\n", op.ID, elapsed)
	*metrics = append(*metrics, Metric{shardID, blockID, "Header", elapsed.Seconds()})
}

// Function for the operator to send bundles to farmers within the shard
func (op *Operator) sendBundles(bundle Bundle, wg *sync.WaitGroup, metrics *[]Metric, shardID int, blockID int, bundleTimes *sync.Map) {
	defer wg.Done()

	startTime := time.Now()
	for _, farmer := range op.Shard.Farmers {
		time.Sleep(timeToSendData(bundleBodySize, op.Bandwidth)) // Simulate sending bundle time based on bandwidth
		farmer.mu.Lock()
		farmer.Bundles = append(farmer.Bundles, bundle)
		farmer.mu.Unlock()
		fmt.Printf("Operator %d sent bundle to farmer %d in shard %d\n", op.ID, farmer.ID, op.Shard.ID)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Operator %d finished sending bundles in %v\n", op.ID, elapsed)
	*metrics = append(*metrics, Metric{shardID, blockID, "Bundle", elapsed.Seconds()})
	bundleTimes.Store(shardID, elapsed.Seconds())
}

func (shard *Shard) proposeBlock(farmer *Farmer, blockID int, numBundles int, wg *sync.WaitGroup, metrics *[]Metric) {
	defer wg.Done()

	startTime := time.Now()
	if farmer != nil {
		time.Sleep(randomNetworkDelay())
		if len(farmer.Bundles) < numBundles {
			numBundles = len(farmer.Bundles)
		}
		newBlock := Block{ID: blockID, Bundles: farmer.Bundles[:numBundles]}

		// Longest chain simulation
		shard.mu.Lock()
		isConsensus := true
		for _, f := range shard.Farmers {
			time.Sleep(randomNetworkDelay())
			if len(f.Chain) < len(farmer.Chain) {
				isConsensus = false
				break
			}
		}

		if isConsensus {
			for _, f := range shard.Farmers {
				f.mu.Lock()
				f.Chain = append(f.Chain, newBlock)
				f.mu.Unlock()
			}
			fmt.Printf("Farmer %d of shard %d proposed block %d and consensus achieved\n", farmer.ID, shard.ID, blockID)
		} else {
			fmt.Printf("Farmer %d of shard %d proposed block %d but consensus not achieved\n", farmer.ID, shard.ID, blockID)
		}
		shard.mu.Unlock()
	} else {
		fmt.Printf("Shard %d has no farmer to propose a block\n", shard.ID)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Farmer %d finished proposing block %d in %v\n", farmer.ID, blockID, elapsed)
	*metrics = append(*metrics, Metric{shard.ID, blockID, "Proposal", elapsed.Seconds()})
}

func simulateNetwork(numShards, numFarmers, numOperators, numBlocks int, bundlesPerShard []int, bandwidth int) {
	var metrics []Metric
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	allFarmers := make([]*Farmer, numFarmers)
	for i := 0; i < numFarmers; i++ {
		allFarmers[i] = &Farmer{ID: i + 1, ShardID: (i % numShards) + 1}
	}

	shards := make([]*Shard, numShards)
	for i := 0; i < numShards; i++ {
		farmersInShard := []*Farmer{}
		for _, farmer := range allFarmers {
			if farmer.ShardID == i+1 {
				farmersInShard = append(farmersInShard, farmer)
			}
		}
		shards[i] = &Shard{
			ID:      i + 1,
			Farmers: farmersInShard,
			Leader:  chooseLeader(farmersInShard, rng),
		}

		fmt.Printf("Shard %d: Farmers ", shards[i].ID)
		for _, farmer := range farmersInShard {
			fmt.Printf("%d ", farmer.ID)
		}
		fmt.Println()
	}

	operators := make([]*Operator, numOperators)
	for i := 0; i < numOperators; i++ {
		if i < len(shards) {
			operators[i] = &Operator{
				ID:        i + 1,
				Shard:     shards[i],
				Bandwidth: bandwidth,
			}
			fmt.Printf("Operator %d assigned to Shard %d\n", operators[i].ID, operators[i].Shard.ID)
		}
	}

	for blockIndex := 0; blockIndex < numBlocks; blockIndex++ {
		fmt.Printf("---- Building Block %d ----\n", blockIndex+1)

		bundleTimes := &sync.Map{}

		var wg sync.WaitGroup
		for _, operator := range operators {
			numBundles := bundlesPerShard[operator.Shard.ID-1]
			for bundleIndex := 0; bundleIndex < numBundles; bundleIndex++ {
				header := BundleHeader{Data: fmt.Sprintf("Header from Operator %d, Bundle %d", operator.ID, bundleIndex+1)}
				bundle := Bundle{Data: fmt.Sprintf("Bundle from Operator %d, Bundle %d", operator.ID, bundleIndex+1)}

				wg.Add(1)
				go operator.sendBundleHeaders(header, allFarmers, &wg, &metrics, operator.Shard.ID, blockIndex+1)

				wg.Add(1)
				go operator.sendBundles(bundle, &wg, &metrics, operator.Shard.ID, blockIndex+1, bundleTimes)
			}
		}

		wg.Wait()

		// Simulate each shard's leader proposing a block
		for _, shard := range shards {
			farmer := chooseFarmer(shard.Farmers, blockIndex)
			numBundles := bundlesPerShard[shard.ID-1]
			wg.Add(1)
			go shard.proposeBlock(farmer, blockIndex+1, numBundles, &wg, &metrics)
		}

		wg.Wait()
	}

	saveMetricsToCSV(&metrics, "benchmark.csv")
}

func saveMetricsToCSV(metrics *[]Metric, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ShardID", "BlockID", "Stage", "TimeTaken"})

	// Sort metrics by ShardID, BlockID, and Stage
	sort.SliceStable(*metrics, func(i, j int) bool {
		if (*metrics)[i].ShardID != (*metrics)[j].ShardID {
			return (*metrics)[i].ShardID < (*metrics)[j].ShardID
		}
		if (*metrics)[i].BlockID != (*metrics)[j].BlockID {
			return (*metrics)[i].BlockID < (*metrics)[j].BlockID
		}
		return (*metrics)[i].Stage < (*metrics)[j].Stage
	})

	// Calculate total time by summing up bundle and proposal times for each shard and block
	totalTimes := make(map[string]float64)
	for _, metric := range *metrics {
		key := fmt.Sprintf("%d-%d", metric.ShardID, metric.BlockID)
		if metric.Stage != "Header" {
			totalTimes[key] += metric.TimeTaken
		}
	}

	for _, metric := range *metrics {
		if metric.Stage != "Total" {
			writer.Write([]string{
				fmt.Sprintf("%d", metric.ShardID),
				fmt.Sprintf("%d", metric.BlockID),
				metric.Stage,
				fmt.Sprintf("%f", metric.TimeTaken),
			})
		}
	}

	for key, totalTime := range totalTimes {
		var shardID, blockID int
		fmt.Sscanf(key, "%d-%d", &shardID, &blockID)
		writer.Write([]string{
			fmt.Sprintf("%d", shardID),
			fmt.Sprintf("%d", blockID),
			"Total",
			fmt.Sprintf("%f", totalTime),
		})
	}

	fmt.Println("Metrics saved to", filename)
}

func main() {
	numShards := 2
	numFarmers := 5
	numOperators := 2
	numBlocks := 1                 // Number of blocks to build
	bundlesPerShard := []int{1, 1} // Number of bundles in each shard block
	bandwidth := 10                // Bandwidth in Mbps

	simulateNetwork(numShards, numFarmers, numOperators, numBlocks, bundlesPerShard, bandwidth)
}
