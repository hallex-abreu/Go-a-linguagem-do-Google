package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const MONITORAMENTO = 3
const DELAY = 5

func main() {
	showIntroduction()
	for {
		showMenu()

		command := readCommand()

		switch command {
		case 1:
			initWatch()
		case 2:
			fmt.Println("Exibindo logs")
			printLog()
		case 3:
			fmt.Println("Saindo do programa")
			os.Exit(0)
		default:
			fmt.Println("Não conhecemos esse comando")
			os.Exit(-1)
		}
	}
}

func showIntroduction() {
	name := "Hallex"
	var version float32 = 1.2

	fmt.Println("Olá, Sr.", name)
	fmt.Println("Este programa está na versão", version)
}

func showMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir logs")
	fmt.Println("0 - Sair do programa")
}

func readCommand() int {
	var readCommand int
	fmt.Scan(&readCommand)

	return readCommand
}

func initWatch() {
	fmt.Println("Monitorando...")
	sites := readySiteFile()

	for i := 0; i < MONITORAMENTO; i++ {
		for i, site := range sites {
			testingSite(i, site)
		}
		time.Sleep(DELAY * time.Second)
	}
}

func testingSite(indice int, site string) {
	result, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	if result.StatusCode == 200 {
		fmt.Println("Testando site:", indice, ", Url:", site, "foi carregado com sucesso!")
		registerLog(site, true)
	} else {
		fmt.Println("Testando site:", indice, ", Url:", site, "está com problemas, status code:", result.StatusCode)
		registerLog(site, false)
	}
}

func readySiteFile() []string {
	var sites []string
	file, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	read := bufio.NewReader(file)

	for {
		row, err := read.ReadString('\n')
		row = strings.TrimSpace(row)
		sites = append(sites, row)

		if err == io.EOF {
			break
		}
	}

	file.Close()

	return sites
}

func registerLog(site string, status bool) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + site + "-online: " + strconv.FormatBool(status) + "\n")

	file.Close()
}

func printLog() {
	file, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	fmt.Println(string(file))
}
