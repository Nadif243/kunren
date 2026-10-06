package worker

import (
	"fmt"
	"kunnrenengine/internal/scraper"
	"sync"
)

// JobResult bundles the URL with the sentences it found (or the error it hit)
type JobResult struct {
	URL       string
	Sentences []string
	Err       error
}

// worker is the infinite loop that runs inside each Goroutine.
// It listens to the jobs channel, scrapes the URL, and sends data to the results channel.
func worker(id int, jobs <-chan string, results chan<- JobResult, wg *sync.WaitGroup) {
	// When this worker eventually dies, tell the WaitGroup it's done
	defer wg.Done()

	// Constantly pull URLs from the channel until the channel is closed
	for url := range jobs {
		fmt.Printf("[Worker %d] Engaging target: %s\n", id, url)

		sentences, err := scraper.ExtractSentences(url)

		// Push the results back to the main thread
		results <- JobResult{
			URL:       url,
			Sentences: sentences,
			Err:       err,
		}
	}
}

// RunPool is the main orchestrator. It takes a list of URLs, spins up the workers,
// and returns all the gathered sentences when everything is finished.
func RunPool(urls []string, numWorkers int) []JobResult {
	// 1. Setup the Channels
	jobs := make(chan string, len(urls))
	results := make(chan JobResult, len(urls))
	var wg sync.WaitGroup

	// 2. Boot up the Worker Goroutines
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// 3. Dump all the URLs into the jobs queue
	for _, url := range urls {
		jobs <- url
	}
	// Close the jobs channel so workers know no more URLs are coming
	close(jobs)

	// 4. Wait for all workers to finish in the background
	go func() {
		wg.Wait()
		close(results) // Close results only AFTER all workers finish
	}()

	// 5. Collect all the data
	var allData []JobResult
	for result := range results {
		allData = append(allData, result)
	}

	return allData
}
