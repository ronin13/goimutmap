package immap

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
)

func ExampleContextMapper() {

	var wg sync.WaitGroup
	mapper, cFunc := NewcontextMapper(context.Background())
	defer cFunc()

	numRoutines := 1000

	updateMap := func() {
		defer wg.Done()
		pos := rand.Intn(101)
		val, _ := mapper.Exists(pos)
		if val == nil {
			val = make([]int, 0)
		}

		cint := val.([]int)

		newslice := make([]int, len(cint))
		copy(newslice, cint)
		newslice = append(newslice, rand.Perm(10)...)

		mapper.Add(pos, newslice)

	}

	wg.Add(numRoutines)
	for x := 0; x < numRoutines; x++ {
		go updateMap()
	}

	wg.Wait()

	counter := 0
	for mb := range mapper.Iterate() {
		counter += mb.key.(int)
		//fmt.Println(mb.mapRef.([]int))
	}
	fmt.Println("Sum of all keys", counter)

	// Output:
	// Sum of all keys 5050

}
