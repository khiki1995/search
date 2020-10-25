package search

import (
	"context"
	"testing"
	"log"
)

func TestSearch(t *testing.T){
	root := context.Background()
	ctx, cancel := context.WithCancel(root)
	files := []string{"./file1.txt","./file2.txt"}
	// var answerResult []Result
	// answerResult = append(answerResult, Result{"khiki", "khiki;1;12;test;khiki", int64(2), int64(1),})
	// answerResult = append(answerResult, Result{"khiki", "eb7af339-e46a-417e-84b9-d3c211d028f0;1;12;test;khiki", int64(4), int64(48),})
	
	ch := All(ctx, "khiki", files)
	for val := range ch {
		// if !reflect.DeepEqual(answerResult, val) {
		// 	t.Errorf("want = %v , got = %v",answerResult,  val)
		// }
		log.Println(val)
	}
	// log.Println(<-ch)
	cancel()

}