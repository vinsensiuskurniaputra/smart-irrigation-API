package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type PlantQuery struct {
	Type           string  // jenis tanaman
	Temperature    float64 // suhu (°C)
	HumidityAir    float64 // kelembapan udara (%)
	HumiditySoil   float64 // kelembapan tanah (%)
	Question       string  // pertanyaan ke tanaman
}

func AskPlant(query PlantQuery) (string, error) {
	apiKey := os.Getenv("CEREBRAS_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("CEREBRAS_API_KEY belum diset")
	}

	url := "https://api.cerebras.ai/v1/chat/completions"

	// Buat prompt dinamis
	prompt := fmt.Sprintf(
		"Saya adalah tanaman %s. Kondisi saya sekarang: suhu %.1f°C, kelembapan udara %.1f%%, kelembapan tanah %.1f%%. Pertanyaan: %s",
		query.Type, query.Temperature, query.HumidityAir, query.HumiditySoil, query.Question,
	)

	payload := map[string]interface{}{
		"model": "llama-3.3-70b",
		"messages": []map[string]string{
			{"role": "system", "content": "Kamu adalah tanaman yang bisa berbicara dengan manusia. Jawablah seolah-olah kamu adalah tanaman tersebut."},
			{"role": "user", "content": prompt},
		},
	}

	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
