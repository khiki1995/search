package search

import (
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

func All(ctx context.Context, phrase string, files []string) <-chan []Result {
	ch := make(chan []Result, len(files))
	chanFiles := make(chan []Result)
	defer close(chanFiles)
	for _, path := range files {	
		path, err := filepath.Abs(path)
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
		if len(arrTxt) > 0 {
			go func(chanFiles chan []Result, arrTxt []string, phrase string){
				var fileResult []Result
				for line, str := range arrTxt {
					if strings.Contains(str, phrase){	
						result := Result{phrase, str , int64(line)+1, int64(strings.Index(str, phrase))+1,}
						fileResult = append(fileResult, result)
					}
				}		
				chanFiles <- fileResult
			}(chanFiles, arrTxt, phrase)
		}
	}
	for  i := 0; i < len(files); i++ {
		select {
			case val := <- chanFiles:
				ch <- val
		}
		if i+1 == len(files){
			close(ch)
		}
	}
	return ch
}