package search

import (
	"io/ioutil"
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
	ch := make(chan Result)
	wg := sync.WaitGroup{}
	result := Result{}
	for i := 0; i < len(files); i++ {

		data, err := ioutil.ReadFile(files[i])
		if err != nil {
			log.Println("ошибка при открытии файла : ", err)
		}
		filetext := string(data)

		if strings.Contains(filetext, phrase) {
			res := FindAnyTextInFile(phrase, filetext)
			if (Result{}) != res {
				result = res
				break
			}
		}
	}
	wg.Add(1)
	go func(ctx context.Context, ch chan<- Result) {
		defer wg.Done()
		if (Result{}) != result {
			ch <- result
		} 
	}(ctx, ch)
	go func() {
		wg.Wait()
		close(ch)
	}()	
	return ch
}