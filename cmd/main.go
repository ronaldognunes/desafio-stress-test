package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
)

var (
	url        string
	totalReqs  int
	concurrent int
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "stress-test",
		Short: "Teste de carga para um serviço web",
		Run:   runLoadTest,
	}

	rootCmd.Flags().StringVar(&url, "url", "", "URL do serviço a ser testado (obrigatório)")
	rootCmd.Flags().IntVar(&totalReqs, "requests", 100, "Número total de requests")
	rootCmd.Flags().IntVar(&concurrent, "concurrency", 10, "Número de chamadas simultâneas")
	rootCmd.MarkFlagRequired("url")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Erro ao executar comando:", err)
	}
}

func runLoadTest(cmd *cobra.Command, args []string) {
	fmt.Printf("Iniciando teste de carga em %s com %d requests e %d concorrentes\n", url, totalReqs, concurrent)

	start := time.Now()

	var successCount int32
	var statusCodes sync.Map

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrent)

	for i := 0; i < totalReqs; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			resp, err := http.Get(url)
			if err == nil {

				atomic.AddInt32(&successCount, 1)
				statusCodes.Store(resp.StatusCode, struct{}{})
				resp.Body.Close()
				fmt.Printf("Request %d: %d\n", i, resp.StatusCode)
			} else {
				fmt.Printf("Request %d: %s\n", i, err)
			}
		}()
	}

	wg.Wait()
	totalTime := time.Since(start)

	fmt.Println("\n--- Relatório ---")
	fmt.Printf("Tempo total: %v\n", totalTime)
	fmt.Printf("Requests bem-sucedidos: %d/%d\n", successCount, totalReqs)
	fmt.Println("Distribuição de Status HTTP:")
	statusCodes.Range(func(key, value interface{}) bool {
		fmt.Printf("Status %d\n", key.(int))
		return true
	})
}
