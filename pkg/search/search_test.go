package search

import (
	"log"
	//"reflect"
	"context"
	"testing"
)

func TestAll_positive(t *testing.T){
	root := context.Background()
	//ctx, cancel := context.WithCancel(root)
	files := []string{"./file1.txt","./file1.txt","./file1.txt"}
	var answerResult []Result
	answerResult = append(answerResult, Result{"khiki", "eb7af339-e46a-417e-84b9-d3c211d028f0;1;12;test;khiki", int64(3), int64(48),})
	answerResult = append(answerResult, Result{"khiki", "eb7af339-e46a-417e-84b9-d3c211d028f0;1;12;test;khiki", int64(5), int64(48),})
	
	ch := All(root, "khiki", files)
	for ans := range ch{
		// if !reflect.DeepEqual(answerResult, ans) {
		// 	t.Errorf("want = %v , got = %v",answerResult,  ans)
		// }
		log.Println(ans)		
	}
	log.Println("Хорошо")	

}