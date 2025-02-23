package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func main() {

	// Fetch and validate arguments
	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatal("Incorrect arguments. Usage: apriori [filename] [min support level]")
	}

	filename := args[0]
	minsumStr := args[1]

	// Open file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Parse dataset from file
	var dataset [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		intStrs := strings.Split(txt, " ")

		var transaction []int
		for _, intStr := range intStrs {
			item, err := strconv.ParseInt(intStr, 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			transaction = append(transaction, int(item))
		}
		dataset = append(dataset, transaction)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Parse min support level from args
	minsum, err := strconv.ParseFloat(minsumStr, 64)
	if err != nil {
		log.Fatal(err)
	}

	// Get initial candidates
	frequentItems := make([][]int, 0)
	initialCandMap := make(map[int]bool)
	for tranIdx := range dataset {
		tran := dataset[tranIdx]
		for itemIdx := range tran {
			item := tran[itemIdx]
			initialCandMap[item] = true
		}
	}
	var initialCands [][]int
	for key := range initialCandMap {
		initialCands = append(initialCands, []int{key})
	}
	slices.SortFunc(initialCands, sliceSort)

	// Get frequent items for k = 1
	tempFreq := generateFrequentItems(dataset, initialCands, minsum)

	for k := 2; len(tempFreq) > 0; k++ {
		// Append frequent items to overall list
		frequentItems = append(frequentItems, tempFreq...)
		// Generate list of candidate sets from previous frequent sets
		candidates := generateCandidates(tempFreq)
		// Prune candidates using previous frequent sets
		prunedCandidates := pruneCandidates(candidates, tempFreq)
		// Scan DB and generate list of frequent sets from pruned candidate list
		tempFreq = generateFrequentItems(dataset, prunedCandidates, minsum)
	}

	fmt.Printf("Minimum Support Level: %.0f%%\n", minsum)
	fmt.Println("Frequent Items:")
	for _, item := range frequentItems {
		fmt.Printf("%v\n", item)
	}
}

// Generates possible candidates of length k from frequent items of length k-1
func generateCandidates(freqSets [][]int) [][]int {
	candidates := make([][]int, 0)

	freqLength := len(freqSets)
	for setIdx, set := range freqSets {
		if setIdx == freqLength-1 {
			break
		}
		setLength := len(set)
		for _, set2 := range freqSets[setIdx+1:] {
			appendSets := true
			for itemIdx, item := range set {
				if (itemIdx < setLength-1) && item != set2[itemIdx] {
					appendSets = false
				}
			}
			if appendSets {
				candidate := set[:setLength]
				candidate = append(candidate, set2[setLength-1])
				candidates = append(candidates, candidate)
			}
		}
	}

	return candidates
}

// Prunes candidates of length k using frequent items of length k-1
func pruneCandidates(candidates [][]int, freqSets [][]int) [][]int {
	prunedCandidates := make([][]int, 0)

	for _, candidate := range candidates {
		candidateLength := len(candidate)
		if candidateLength == 1 {
			prunedCandidates = append(prunedCandidates, candidate)
		}

		p1 := candidate[:candidateLength-1]
		foundP1 := false
		p2 := candidate[1:candidateLength]
		foundP2 := false

		for _, freq := range freqSets {
			if reflect.DeepEqual(p1, freq) {
				foundP1 = true
			}
			if reflect.DeepEqual(p2, freq) {
				foundP2 = true
			}
		}

		if foundP1 && foundP2 {
			prunedCandidates = append(prunedCandidates, candidate)
		}
	}

	return prunedCandidates
}

// Finds frequent items
func generateFrequentItems(dataset [][]int, candidates [][]int, minsum float64) [][]int {
	// Get absolute support of all items
	// idx of candidate -> freq
	candFreqs := make(map[int]int)
	for tranIdx := range dataset {
		tran := dataset[tranIdx]
		for candIdx, cand := range candidates {
			found := true
			for _, c := range cand {
				if find(tran, c) == -1 {
					found = false
				}
			}
			if found {
				candFreqs[candIdx]++
			}
		}
	}

	// Get frequent items
	datasetLength := len(dataset)
	tempFreq := make([][]int, 0)
	for candIdx, freq := range candFreqs {
		relativeSupport := (float64(freq) / float64(datasetLength)) * 100
		if relativeSupport >= minsum {
			tempFreq = append(tempFreq, candidates[candIdx])
		}
	}
	slices.SortFunc(tempFreq, sliceSort)

	return tempFreq
}

// Sort slices of slices of integers
func sliceSort(a []int, b []int) int {
	sort.Ints(a)
	sort.Ints(b)

	if len(b) < len(a) {
		return 1
	}
	for idx, attr := range a {
		attr2 := b[idx]
		if attr < attr2 {
			return -1
		}
	}

	if len(a) == len(b) {
		return 0
	}

	return 1
}

// Finds integer in slice
func find(slice []int, val int) int {
	for idx, v := range slice {
		if v == val {
			return idx
		}
	}
	return -1
}
