package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Estrutura para armazenar dados do CEP
type CEPInfo struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

// -------------------------------
// Consulta ViaCEP
// -------------------------------
func consultaViaCEP(cep string) (*CEPInfo, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição ViaCEP: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var resultado CEPInfo
	if err := json.Unmarshal(body, &resultado); err != nil {
		return nil, fmt.Errorf("erro ao ler JSON ViaCEP: %v", err)
	}

	// Se o CEP não existir, ViaCEP devolve campo "cep" vazio
	if resultado.Cep == "" {
		return nil, fmt.Errorf("CEP não encontrado no ViaCEP")
	}

	return &resultado, nil
}

// -------------------------------
// Consulta BrasilAPI
// -------------------------------
func consultaBrasilAPI(cep string) (*CEPInfo, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição BrasilAPI: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var resultado CEPInfo
	if err := json.Unmarshal(body, &resultado); err != nil {
		return nil, fmt.Errorf("erro ao ler JSON BrasilAPI: %v", err)
	}

	// Mesma validação
	if resultado.Cep == "" {
		return nil, fmt.Errorf("CEP não encontrado na BrasilAPI")
	}

	return &resultado, nil
}

// -------------------------------
// Função principal
// -------------------------------
func main() {
	var cep string

	fmt.Print("Digite um CEP (somente números): ")
	fmt.Scan(&cep)

	cep = strings.TrimSpace(cep)

	fmt.Println("\n🔎 Consultando CEP...\n")

	// consulta ViaCEP
	viaCEP, err1 := consultaViaCEP(cep)
	if err1 != nil {
		fmt.Println("ViaCEP: ❌", err1)
	} else {
		fmt.Println("ViaCEP: ✅ Encontrado!")
		fmt.Printf("CEP: %s\nRua: %s\nBairro: %s\nCidade: %s\nUF: %s\n\n",
			viaCEP.Cep, viaCEP.Logradouro, viaCEP.Bairro, viaCEP.Localidade, viaCEP.Uf)
	}

	// consulta BrasilAPI
	brAPI, err2 := consultaBrasilAPI(cep)
	if err2 != nil {
		fmt.Println("BrasilAPI: ❌", err2)
	} else {
		fmt.Println("BrasilAPI: ✅ Encontrado!")
		fmt.Printf("CEP: %s\nRua: %s\nBairro: %s\nCidade: %s\nUF: %s\n\n",
			brAPI.Cep, brAPI.Logradouro, brAPI.Bairro, brAPI.Localidade, brAPI.Uf)
	}

	fmt.Println("🔚 Fim da consulta.")
}
