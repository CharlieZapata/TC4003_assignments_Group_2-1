package cos418_hw1_1

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

// Sum numbers from channel `nums` and output sum to `out`.
// You should only output to `out` once.
// Do NOT modify function signature.
func sumWorker(nums chan int, out chan int) {
	// TODO: implement me
	// HINT: use for loop over `nums`
	sum := 0 //Local var for sum of each gorutine
	for i := 0; i < cap(nums); i++ {
		sum += <-nums
	}
	out <- sum //Write the sum at the end of each gorutine
}

// Read integers from the file `fileName` and return sum of all values.
// This function must launch `num` go routines running
// `sumWorker` to find the sum of the values concurrently.
// You should use `checkError` to handle potential errors.
// Do NOT modify function signature.
func sum(num int, fileName string) int {
	// TODO: implement me
	// HINT: use `readInts` and `sumWorkers`
	// HINT: used buffered channels for splitting numbers between workers
	// return 0
	file, err := os.Open(fileName) //Open the file for reading
	checkError(err)
	defer file.Close()

	numSlice, err := readInts(file)
	checkError(err)

	out := make(chan int, num) //Channel to save all the partial sums from gorutines
	n := len(numSlice) / num   //Size of each slice for each gorutine
	for i := 0; i < num; i++ {
		in := make(chan int, n) //One channel per gorutine for saving the values to sum
		for j := 0; j < n; j++ {
			in <- numSlice[(i*n)+j] //Passing to the channel the input values
		}
		go sumWorker(in, out) //Calling num gorutines
	}
	sum := 0
	for i := 0; i < num; i++ {
		sum += <-out //Reading the num partial results from the output channel
	}
	return sum
}

// Read a list of integers separated by whitespace from `r`.
// Return the integers successfully read with no error, or
// an empty slice of integers and the error that occurred.
// Do NOT modify this function.
func readInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var elems []int
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return elems, err
		}
		elems = append(elems, val)
	}
	return elems, nil
}
