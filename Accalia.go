package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var green = "\033[32m"
var red = "\033[31m"
var white = "\033[37m"

func main() {
	var wg sync.WaitGroup
	var filePath string
	var websitePath string

	flag.StringVar(&websitePath, "w", "", "Website URL for bruteforce (default *Empty String*)")
	flag.StringVar(&filePath, "f", "res/default.txt", "File location for bruteforce")
	workerGoNumber := flag.Int("g", 200, "Set the number of worker goroutines")
	silence := flag.Bool("s", false, "Show cool ASCII art? (default true)")
	flag.Parse()

	if websitePath == "" {
		fmt.Println("Website URL is empty")
		os.Exit(1)
	}
	if *silence == false {
		header()
	}

	fmt.Println("Creating " + strconv.Itoa(*workerGoNumber) + " crawl worker threads")
	words := make(chan string)
	for goroutines := 0; goroutines < *workerGoNumber; goroutines++ {
		go crawlWorker(&wg, words, &websitePath)
	}
	readLineAndSendToChan(&filePath, words, &wg)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func checkError(err *error) {
	if *err != nil {
		fmt.Println(string(red), "[-] ", string(white), *err)
	}
}

func loadFileToRAM(fileName string) (RAMAdress *string) {
	fileContent, err := ioutil.ReadFile(fileName)
	checkError(&err)
	fileContentString := string(fileContent)
	return &fileContentString
}

func header() {
	fmt.Println("\n			    Accalia")
	fmt.Println("			Written by - Idol")
	fmt.Print("\n")
	RAMAdress := loadFileToRAM("res/Accalia.txt")
	fmt.Println(*RAMAdress)
}

func readLineAndSendToChan(path *string, channel chan string, wg *sync.WaitGroup) {
	inFile, err := os.Open(*path)
	checkError(&err)
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		wg.Add(1)
		channel <- scanner.Text()
	}
	wg.Wait()
}

func crawlWorker(wg *sync.WaitGroup, words chan string, rootPath *string) {
	for word := range words {
		response, err := http.Get(*rootPath + word)
		checkError(&err)
		if response.StatusCode == 200 {
			fmt.Print(string(green), "[+] ", string(white))
			fmt.Println(*rootPath + word)
		}
		wg.Done()
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
