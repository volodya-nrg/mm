package mylib

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type MyLib struct {
	searchDir       string
	amountParallels int
	noder           Noder
	chResponses     chan Response
	// эти ниже два св-ва все таки для того чтоб в итоге можно было выйти из программы
	queueFiles []string
	mtx        sync.Mutex
}

func (s *MyLib) Run(ctx context.Context) <-chan Response {
	go s.walkToFolders(s.searchDir)   // если файлов много, то на фоне будет копится очередь из файлов
	time.Sleep(10 * time.Millisecond) // все таки явно подождем, на всякий случай

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				s.workerPool()
			}
		}
	}()

	return s.chResponses
}

// walkToFolders пробегается по файловой системе и кладет в очередь пути к файлам
func (s *MyLib) walkToFolders(dirPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		s.chResponses <- Response{
			Err: fmt.Errorf("failed to read dir (%s): %w", dirPath, err),
		}
		return
	}

	var result []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirName := entry.Name()
			if dirName == "." || dirName == ".." { // на всякий случай
				continue
			}

			s.walkToFolders(filepath.Join(dirPath, dirName))
		} else {
			result = append(result, filepath.Join(dirPath, entry.Name()))
		}
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.queueFiles = append(s.queueFiles, result...)
}

// workerPool берет часть сначала и обрабатывает такое кол-во сколько нод имеется
func (s *MyLib) workerPool() {
	if len(s.queueFiles) == 0 { // предохранитель на всякий случай, защита от двойного close ch
		return
	}

	defer func() {
		if len(s.queueFiles) == 0 {
			close(s.chResponses) // закрываем канал чтоб в итоге программа завершилась
		}
	}()

	amountElementsExecuted := 0
	wg := sync.WaitGroup{}

	for k, incomingFile := range s.queueFiles {
		if k > s.amountParallels { // обрабатываем только до определенного лимита
			break
		}

		wg.Add(1)
		go func() { // каждый pipeline в своем потоке, но в pipeline соблюдается очередь
			defer wg.Done()

			node := s.noder
			for node != nil {
				resp := Response{
					Name: node.GetName(),
				}

				tmpResult, err := node.Execute(incomingFile)
				if err != nil {
					resp.Err = fmt.Errorf("failed to execute node: %w", err)
				} else {
					resp.Result = tmpResult
				}

				s.chResponses <- resp
				node = node.GetNextNode()
			}
		}()
		amountElementsExecuted++
	}
	wg.Wait()
	slog.Debug("-----new batch-----")

	// уберем те данные которые уже обработали
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.queueFiles = s.queueFiles[amountElementsExecuted:]
}

func NewMyLib(searchDir string, amountParallels int, noder Noder) *MyLib {
	return &MyLib{
		searchDir:       searchDir,
		amountParallels: amountParallels,
		noder:           noder,
		chResponses:     make(chan Response, amountParallels), // пусть буфер будет
	}
}
