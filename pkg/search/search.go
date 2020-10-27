package search

import (
	"sync"
	"path/filepath"
	"log"
	"strings"
	"context"
	"os"
)

type Result struct {
	Phrase string
	Line string
	LineNum int64
	ColNum int64
}

func FindTextInFile(phrase string, filename string) []Result {
	path, err := filepath.Abs(filename)
	if err != nil {
		log.Printf("Error with dirict, dir = %v", path)
	}
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error with open file! error = %v", err)
	}
	buf := make([]byte,4096)
	read, err := file.Read(buf) 
	if err != nil {
		log.Printf("Error in reading file! error = %v", err)
	}
	data := string(buf[:read])
	arrTxt := strings.Split(data,"\n")
	var fileResult []Result
	for line, str := range arrTxt {
		if strings.Contains(str, phrase){	
			result := Result{phrase, str , int64(line)+1, int64(strings.Index(str, phrase))+1,}
			fileResult = append(fileResult, result)
		}
	}		
	return fileResult
}

func FindAnyTextInFile(phrase string, filename string) Result {
	path, err := filepath.Abs(filename)
	if err != nil {
		log.Printf("Error with dirict, dir = %v", path)
	}
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error with open file! error = %v", err)
	}
	buf := make([]byte,4096)
	read, err := file.Read(buf) 
	if err != nil {
		log.Printf("Error in reading file! error = %v", err)
	}
	data := string(buf[:read])
	arrTxt := strings.Split(data,"\n")
	var result Result
	for line, str := range arrTxt {
		if strings.Contains(str, phrase){	
			return Result{phrase, str , int64(line)+1, int64(strings.Index(str, phrase))+1,}
		}
	}
	return result	
}

func All(ctx context.Context, phrase string, files []string) <-chan []Result {
	ch := make(chan []Result)
	wg := sync.WaitGroup{}
	
	for i := 0; i < len(files); i++ {
		wg.Add(1)
		go func(file string, ch chan<- []Result){
			defer wg.Done()
			res := FindTextInFile(phrase, file)

			if len(res) > 0 {
				ch <- res
			}
		}(files[i], ch)
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	return ch
}

func Any(ctx context.Context, phrase string, files []string) <-chan Result {
	ctx, cancel :=  context.WithCancel(ctx)
	ch := make(chan Result, 1)
	chFirst := make(chan Result)

	for i := 0; i < len(files); i++ {
		go func(ctx context.Context, chFirst chan Result, phrase string, file string, lastGor bool) {
			select {
				case <-ctx.Done():
					close(chFirst)
				default:
					res := FindAnyMatchTextInFile(phrase, file)
					if res != (Result{}){
						chFirst <- res
					}else {
						if lastGor {
							chFirst <- Result{}
						}
					}
			}
		}(ctx, chFirst, phrase, files[i], i+1 == len(files))
		
	}
	val := <-chFirst
	cancel()
	ch <- val
	close(ch)
	return ch
}

func FindAnyMatchTextInFile(phrase, filetext string) (res Result) {

	//ch := make(chan Result)

	temp := strings.Split(filetext, "\n")

	for i, line := range temp {
		//fmt.Println("[", i+1, "]\t", line)
		if strings.Contains(line, phrase) {

			return Result{
				Phrase:  phrase,
				Line:    line,
				LineNum: int64(i + 1),
				ColNum:  int64(strings.Index(line, phrase)) + 1,
			}

		}
	}

	return res
}