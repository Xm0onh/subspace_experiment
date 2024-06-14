## Parameters

The simulation allows you to configure the following parameters:

1. **Number of Shards (`numShards`)**: The number of shards in the network.
2. **Number of Farmers (`numFarmers`)**: The total number of farmers distributed across the shards.
3. **Number of Operators (`numOperators`)**: The total number of operators managing the shards.
4. **Number of Blocks (`numBlocks`)**: The total number of blocks to be produced in the simulation in each shard (shard blocks).
5. **Bundles per Shard (`bundlesPerShard`)**: An array specifying the number of bundles to be included in each shard's block.
6. **Bandwidth (`bandwidth`)**: The bandwidth limitation, defining how many bundles can be sent per second.

## How to Change Parameters

To change the parameters, you need to modify the values in the `main` function of the `main.go` file. 
```go
func main() {
    numShards :=                // Number of shards in the network
    numFarmers :=               // Total number of farmers
    numOperators :=             // Total number of operators
    numBlocks :=                // Number of blocks to build
    bundlesPerShard := []int{}  // Number of bundles in each shard block
    bandwidth :=                // Bundles per second

    simulateNetwork(numShards, numFarmers, numOperators, numBlocks, bundlesPerShard, bandwidth)
}
```

### Parameters Explanation

1. **`numShards`**: Set this to the number of shards you want in your network.
2. **`numFarmers`**: Define the total number of farmers.
3. **`numOperators`**: Specify the number of operators.
4. **`numBlocks`**: Set this to the number of blocks to be produced. 
5. **`bundlesPerShard`**: This is an array where each element specifies the number of bundles for each shard's block. For example, `bundlesPerShard := []int{3, 4, 5}` means the first shard will have 3 bundles per block, the second shard will have 4, and the third will have 5.
6. **`bandwidth`**: Define the bandwidth limitation. For example, `bandwidth := 5` means five bundles can be sent per second.

## Running the Simulation

After setting the desired parameters, you can run the simulation by executing the `main.go` file. The simulation will output the metrics for each stage and save them to a CSV file named `metrics.csv`.

```sh
go run .
```

The `metrics.csv` file will contain the following columns:

- **ShardID**: The ID of the shard.
- **BlockID**: The ID of the block.
- **Stage**: The stage of the process (e.g., Header, Bundle, Proposal, Total).
- **TimeTaken**: The time taken for the stage in seconds.


## More Details

Farmers are evenly distributed across shards in a round-robin fashion. This means that farmers are assigned to shards based on their index, cycling through the available shards. For example, if there are 5 farmers and 2 shards, farmers 1, 3, and 5 will be assigned to shard 1, and farmers 2 and 4 will be assigned to shard 2. Operators are similarly assigned to shards, where each operator manages a specific shard.

For leadership within shards, the function chooseLeader randomly selects a leader from the farmers within the shard using a pseudo-random number generator (rng).

For block proposals, the function `chooseFarmer` deterministically selects a farmer based on the current round number.