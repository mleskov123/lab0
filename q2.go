package cos418_hw1_1

import (
	"bufio"
	"errors"
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

	if num < 1{
		checkError(errors.New("number of threads must be at least one"))
	}

	file, err := os.Open(fileName)
	checkError(err)
	vals, err := readInts(file)
	checkError(err)

	// no work to do
	if len(vals) == 0 {
		return 0
	}

	// do not need all the threads asked for. we could simply
	// give the extra threads a zero value, but this is extra work than needed
	if len(vals) < num{
		num = len(vals)
	}

	shareOfWork := int(len(vals)/num)
	// each go routine needs individual out channel
	var sumCh = make([]chan int, num)
	for i:=0; i < num; i++{
		sumCh[i] = make(chan int)
	}

	// if there is extra work to be done, some threads will
	// have to put in more work than others
	var extraWork int = len(vals) - (shareOfWork*num)

	var idx int = 0
	for i := 0; i < num; i++ {
		workToDispense := shareOfWork
		if extraWork > 0{
			workToDispense++
			extraWork--
		}
		numsCh := make(chan int, workToDispense)
		for j := 0; j < workToDispense; j++ {
			numsCh <- vals[idx]
			idx++
		}
		close(numsCh)
		go sumWorker(numsCh, sumCh[i])


	}
	var sum int = 0
	for i:=0; i < num; i++ {
		sum = sum + <-sumCh[i]
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
