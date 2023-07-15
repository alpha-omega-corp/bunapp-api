package resources

func sliceToChannel(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func Pipeline(nums []int) {
	// Stage 1
	dataChan := sliceToChannel(nums)

	// stage 2
	finalChan := square(dataChan)
	for n := range finalChan {
		println(n)
	}

}
