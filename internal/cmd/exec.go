/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	tr = &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		ForceAttemptHTTP2: true,
	}
	clientHttp = &http.Client{
		Transport: tr,
		Timeout:   50 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {

			return nil
		},
	}
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Stress test",
	Long:  `Teste de performance que visa avaliar a estabilidade e confiabilidade de um sistema sob condições extremas`,
	Run: func(cmd *cobra.Command, args []string) {

		url, _ := cmd.Flags().GetString("url")
		requests, _ := cmd.Flags().GetString("requests")
		concurrency, _ := cmd.Flags().GetString("concurrency")

		reqCount, errReqConv := strconv.Atoi(requests)
		if errReqConv != nil {
			reqCount = 1
		}
		concurrencyCount, errConcurrencyConv := strconv.Atoi(concurrency)
		if errConcurrencyConv != nil {
			concurrencyCount = 1
		}

		fmt.Printf("Iniciando com os parametros: url: %v - requests: %v - concurrency: %v", url, requests, concurrency)

		result := exectRotine(url, reqCount, concurrencyCount)
		print(result)

	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringP("url", "u", "", "URL do serviço a ser testado")
	execCmd.Flags().StringP("requests", "r", "", "Número total de requests")
	execCmd.Flags().StringP("concurrency", "c", "", "Número de chamadas simultâneas")

}

func exectRotine(url string, requests int, concurrency int) string {

	startTime := time.Now()
	report := initReport(url, requests, concurrency, startTime)
	status := make(map[string]int)

	var wg sync.WaitGroup
	wg.Add(requests)
	paralellReqControl := make(chan struct{}, concurrency)
	erroInRequisition := make(chan struct{})

	fmt.Print("\033[s")
	go func() {
		for {
			select {
			case <-erroInRequisition:
				paralellReqControl <- struct{}{}
				wg.Add(1)
				go execGet(&wg, url, status, paralellReqControl, requests)
			}
		}
	}()

	for i := 0; i < requests; i++ {
		paralellReqControl <- struct{}{}
		go execGet(&wg, url, status, paralellReqControl, requests)
	}

	wg.Wait()
	fmt.Print("\033[u\033[K")
	return finishReport(report, status, startTime)

}

func execGet(wg *sync.WaitGroup, url string, statusCounter map[string]int, p <-chan struct{}, total int) {
	defer wg.Done()

	res, err := clientHttp.Get(url)

	if err != nil {
		fmt.Printf("\nErro ao acessar o site: %s\n", err)
		os.Exit(1)
	}

	index := strconv.Itoa(res.StatusCode)
	if _, exists := statusCounter[index]; exists {
		statusCounter[strconv.Itoa(res.StatusCode)] += 1
	} else {
		statusCounter[strconv.Itoa(res.StatusCode)] = 1
	}

	fmt.Print("\033[u\033[K")
	var currentV int
	for c := range statusCounter {
		currentV += statusCounter[c]
	}

	fmt.Printf("\nAguarde -> %v/%v", currentV, total)

	<-p
	fmt.Print("\033[u\033[K")
}

func initReport(url string, requests int, concurrency int, timeStart time.Time) string {
	var x string

	x = "\n================= Relatorio do teste de Stress ================="
	x += fmt.Sprintf("\nUrl testada: %v\nQuantidade de requisicoes: %v\nParalelismo: %v", url, requests, concurrency)
	x += fmt.Sprintf("\nIniciado em   %v", timeStart.Format(("02/01/2006 15:04:05 ")))

	return x
}

func finishReport(r string, statusCounter map[string]int, timeStart time.Time) string {

	end := time.Now()
	r += fmt.Sprintf("\nFinalizado em %v", end.Format(("02/01/2006 15:04:05 ")))
	r += fmt.Sprintf("\nTempo decorrido  %v", end.Sub(timeStart))
	r += "\nStatus code que foram recebidos: "
	var total int
	for i, v := range statusCounter {
		r += fmt.Sprintf("\n  => status %v = %v", i, v)
		total += v
	}
	r += fmt.Sprintf("\nTotal de requisicoes realizadas  %v", total)
	r += "\n================================================================"

	return r
}
