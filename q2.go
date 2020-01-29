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
	var sum int = 0
	for v := range nums{
		sum += v
	}
	out <- sum
	close(out)
	return
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
	file, err := os.Open(fileName)
	checkError(err)
	vals, err := readInts(file)
	checkError(err)

	shareOfWork := int(len(vals)/num)
	// each go routine needs individual out channel
	var sumCh = make([]chan int, num)
	for i:=0; i < num; i++{
		sumCh[i] = make(chan int)
	}
	var sum int = 0
	var idx int = 0
	for i := 0; i < num-1; i++ {
		numsCh := make(chan int, shareOfWork)
		for j := 0; j < shareOfWork; j++ {
			numsCh <- vals[idx]
			idx++
		}
		close(numsCh)
		go sumWorker(numsCh, sumCh[i])

	}
	{
		numsCh := make(chan int, len(vals) - idx)
		for ; idx < len(vals); idx++ {
			tmp := vals[idx]
			numsCh <- tmp
		}
		close(numsCh)
		go sumWorker(numsCh, sumCh[num-1])
	}
	var count int = 1
	for i:=0; i < num; i++ {

		sum = sum + <-sumCh[i]
		count++
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
