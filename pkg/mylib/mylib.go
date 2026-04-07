package mylib

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

type MyLib struct {
	noder       Noder
	chunkes     [][]string
	chFiles     chan []string
	chResponses chan Response
}

func (s *MyLib) Run(ctx context.Context) <-chan Response {
	go func() {
		defer close(s.chFiles)
		mtx := sync.Mutex{}
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if len(s.chunkes) == 0 { // если закончились данные, то выйдем
					return
				}
				sl := s.chunkes[0]

				mtx.Lock()
				s.chunkes = s.chunkes[1:]
				mtx.Unlock()

				s.chFiles <- sl
			}
		}
	}()
	go func() {
		defer close(s.chResponses)
		for {
			select {
			case <-ctx.Done():
				return
			case sl, isOpened := <-s.chFiles:
				if !isOpened {
					return
				}
				s.handler(sl)
			}
		}
	}()
	return s.chResponses
}

func (s *MyLib) handler(files []string) {
	wg := sync.WaitGroup{}
	for _, filePath := range files {
		wg.Add(1)
		go func() { // каждый pipeline в своем потоке, но в pipeline соблюдается очередь
			defer wg.Done()

			node := s.noder
			for node != nil {
				resp := Response{
					Name: node.GetName(),
				}

				tmpResult, err := node.Execute(filePath)
				if err != nil {
					resp.Err = fmt.Errorf("failed to execute node: %w", err)
				} else {
					resp.Result = tmpResult
				}

				s.chResponses <- resp
				node = node.GetNextNode()
			}

			// time.Sleep(1 * time.Second)
		}()
	}
	wg.Wait()
}

func NewMyLib(searchDir string, amountParallels int, noder Noder) (*MyLib, error) {
	/*
		Протестировал получение мелких файлов с директорий на 1Gb.
		В среднем выполняется за 800ms, т.е. ни чего страшного.
		Поэтому сразу получим все расположения файлов.
		Использую готовую уже для этого ф-ию, т.к. она короче в написании.
	*/
	var files []string
	err := filepath.WalkDir(searchDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk dir %s: %v", path, err)
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	chunkes := make([][]string, 0, len(files)/amountParallels)
	for chunk := range slices.Chunk(files, amountParallels) {
		chunkes = append(chunkes, chunk)
	}

	return &MyLib{
		noder:       noder,
		chunkes:     chunkes,
		chFiles:     make(chan []string),
		chResponses: make(chan Response),
	}, nil
}
