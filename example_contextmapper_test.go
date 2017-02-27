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
		mapper.Add(pos, append(val.([]int), rand.Perm(10)...))

	}

	wg.Add(numRoutines)
	for x := 0; x < numRoutines; x++ {
		go updateMap()
	}

	wg.Wait()

	counter := 0
	for mb := range mapper.Iterate() {
		counter += mb.value.(int)
		//fmt.Println(mb.mapRef.([]int))
	}
	fmt.Println("Sum of all keys", counter)

	// Output:
	// Sum of all keys 5050

}
