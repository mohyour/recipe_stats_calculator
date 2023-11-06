package app

import (
	"context"
	"encoding/json"
	"log"
	"sync"
)

// Process in batches for faster time, but seems inaccurate as some rows are lost runtime
// TODO: revisit
type processedRecipes struct {
	numOfRows int
}

func ProcessFixtureFileInBatches(file string, numOfWorkers, batchSize int) (Result, error) {
	res := Result{}
	loadedFile, err := loadFile(file)
	if err != nil {
		log.Println("Unable to load file")
		return Result{}, err
	}
	var serviceMu sync.Mutex
	reader := func(ctx context.Context, rowsBatch *[]RecipeData) (<-chan []RecipeData, error) {
		out := make(chan []RecipeData)
		dec := json.NewDecoder(loadedFile)
		t, err := dec.Token()
		if err != nil {
			log.Println(err, "==")
			return nil, err
		}
		if t != json.Delim('[') {
			log.Println(err, "delim")
			return nil, err
		}

		go func() {
			defer close(out)
			for dec.More() {
				var recipe RecipeData
				select {
				case <-ctx.Done():
					return
				default:
					if err := dec.Decode(&recipe); err != nil {
						log.Println(err)
						return
					}
					if len(*rowsBatch) == batchSize {
						out <- *rowsBatch
						*rowsBatch = []RecipeData{} // clear batch
					}
					*rowsBatch = append(*rowsBatch, recipe) // add row to current batch
				}
			}
		}()

		// While there are more tokens in the JSON stream...
		return out, nil
	}
	worker := func(ctx context.Context, rowsBatch <-chan []RecipeData) <-chan processedRecipes {
		out := make(chan processedRecipes)
		go func() {
			defer close(out)

			p := processedRecipes{}

			for rowBatch := range rowsBatch {

				for _, recipe := range rowBatch {
					serviceMu.Lock()
					calculateUniqueRecipes(recipe.Recipe)
					serviceMu.Unlock()

					serviceMu.Lock()
					calculateUniquePostcodes(recipe.Postcode)
					serviceMu.Unlock()

					serviceMu.Lock()
					calculateDeliveryCount(recipe.Delivery, recipe.Postcode)
					serviceMu.Unlock()

					serviceMu.Lock()
					matchRecipeName(recipe.Recipe)
					serviceMu.Unlock()
					serviceMu.Lock()
					p.numOfRows++
					serviceMu.Unlock()
				}
			}
			out <- p
		}()
		return out
	}

	combiner := func(ctx context.Context, inputs ...<-chan processedRecipes) <-chan processedRecipes {
		out := make(chan processedRecipes)

		var wg sync.WaitGroup
		multiplexer := func(p <-chan processedRecipes) {
			defer wg.Done()
			for in := range p {
				select {
				case <-ctx.Done():
				case out <- in:
				}
			}
		}

		// add length of input channels to be consumed by mutiplexer
		wg.Add(len(inputs))
		for _, in := range inputs {
			go multiplexer(in)
		}

		// close channel after all inputs channels are closed
		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rowsBatch := []RecipeData{}
	rowsCh, err := reader(ctx, &rowsBatch)
	if err != nil {
		return Result{}, err
	}

	workersCh := make([]<-chan processedRecipes, numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		workersCh[i] = worker(ctx, rowsCh)
	}

	for processed := range combiner(ctx, workersCh...) {
		// add number of rows processed by worker
		res.NumOfRows += processed.numOfRows
	}

	formatOutputData(res)

	return res, nil
}
