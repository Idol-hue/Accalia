package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var filePath string
	var websitePath string

	flag.StringVar(&websitePath, "w", "", "Website URL for bruteforce. Default-*Empty String*")
	flag.StringVar(&filePath, "f", "res/default.txt", "File location for bruteforce. Default-res/default.txt")
	mode := flag.Int("m", 0, "Set mode of Accalia. Default-0. 0-Crawl")
	workerGoNumber := flag.Int("g", 200, "Set the number of worker goroutines. Default-200")
	silence := flag.Bool("s", false, "Show cool ASCII art? Default-true")
	help := flag.Bool("h", false, "Do you want help? Default-false")
	flag.Parse()
	if *help == true {
		helpMSG()
	}
	if websitePath == "" {
		fmt.Println("Website URL is empty")
		os.Exit(1)
	}
	if *silence == false {
		header()
	}

	if *mode == 0 {
		words := make(chan string)
		for goroutines := 0; goroutines < *workerGoNumber; goroutines++ {
			go crawlWorker(&wg, words, &websitePath)
		}
		readLineAndSendToChan(&filePath, words, &wg)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func checkError(err error) {
	if err != nil {
		fmt.Print(err)
	}
}

func loadFileToRAM(fileName string) (RAMAdress *string) {
	fileContent, err := ioutil.ReadFile(fileName)
	checkError(err)
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

func helpMSG() {
	fmt.Println("Accalia- Written by Idol\n")
	fmt.Println("Usage:\n")
	fmt.Println("        Accalia [arguments]\n")
	fmt.Println("The arguments are:")
	fmt.Println("        m       Set mode of Accalia. Default-0. 0-Crawl")
	fmt.Println("        g       Set the number of worker goroutines. Default-200")
	fmt.Println("        s       Show cool ASCII art? Default-true")
	fmt.Println("        h       Do you want help? Default-false")
	fmt.Println("        f       File location for bruteforce. Default-res/default.txt")
	os.Exit(0)
}

func readLineAndSendToChan(path *string, channel chan string, wg *sync.WaitGroup) {
	inFile, err := os.Open(*path)
	checkError(err)
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
		checkError(err)
		if response.StatusCode == 200 {
			fmt.Println(*rootPath + word)
		}
		wg.Done()
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
