package main

import (
	"encoding/json"
	"flag"
	"github.com/cespare/xxhash/v2"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

var (
	config *Config
)

const (
	OsFilePermissions os.FileMode = 0644
	DefaultConfigPath string      = "./config.json"
	DefaultOutputPath string      = "./output.json"
)

type (
	Config struct {
		ConfigPath string
		OutputPath string

		// Paths to directories, which has duplicates
		Directories []string `json:"dirs"`
	}

	File     string

	Executor struct {
		mutex *sync.Mutex
		wg    *sync.WaitGroup

		// Results, which will be saved in JSON
		// map[xxHash][]Path
		Results map[uint64][]File
	}
)

func isErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Check directories, which specified as paths with duplicates, to exists
// Returns non-existing directory
func (c *Config) IsDirsExists() (dir string) {
	for _, dir := range config.Directories {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return dir
		}
	}
	return
}

// Walkthrough specified directories and save paths to all found files
// Returns slice of found paths
func (c *Config) GetFiles() []File {
	output := make([]File, 0)
	for _, dir := range config.Directories {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			output = append(output, File(path))
			return nil
		})
	}
	return output
}

// Get uint64 hash from a file, using xxHash algorithm
// https://github.com/Cyan4973/xxHash#benchmarks
func (f File) GetHash() uint64 {
	raw, err := ioutil.ReadFile(string(f))
	isErr(err)
	return xxhash.Sum64(raw)
}

// Collect all files and save their hashes to the mapping
func (e *Executor) SaveFileHash(file File) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	defer e.wg.Done()

	hash := file.GetHash()
	e.Results[hash] = append(e.Results[hash], file)
}

func init() {
	config = new(Config)

	flag.StringVar(&config.ConfigPath, "config", DefaultConfigPath, "Use this flag to specify the path to your config file, which has paths to directories with duplicates.")
	flag.StringVar(&config.OutputPath, "output", DefaultOutputPath, "Use this flag to specify the path to the output file with results.")
	flag.Parse()

	raw, err := ioutil.ReadFile(config.ConfigPath)
	isErr(err)
	isErr(json.Unmarshal(raw, &config))
}

func main() {
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGKILL, syscall.SIGHUP)
		defer close(signalChan)

		<-signalChan
		log.Println("Shutting down the program...")
		os.Exit(0)
	}()

	if dir := config.IsDirsExists(); dir != "" {
		log.Fatalf("Directory %s can not be found.", dir)
	}

	files := config.GetFiles()
	log.Printf("Found %d files.", len(files))

	exec := &Executor{
		mutex:   new(sync.Mutex),
		wg:      new(sync.WaitGroup),
		Results: make(map[uint64][]File),
	}
	for _, file := range files {
		exec.wg.Add(1)
		go exec.SaveFileHash(file)
	}
	exec.wg.Wait()

	for hash, files := range exec.Results {
		if len(files) <= 1 {
			delete(exec.Results, hash)
		}
	}
	log.Printf("Found %d hashes of duplicates.", len(exec.Results))

	raw, err := json.Marshal(exec.Results)
	isErr(err)
	isErr(ioutil.WriteFile(config.OutputPath, raw, OsFilePermissions))
	log.Printf("All results saved to %s", config.OutputPath)
}
