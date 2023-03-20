package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Billing struct {
	Object         string  `json:"object"`
	TotalGranted   float64 `json:"total_granted"`
	TotalUsed      float64 `json:"total_used"`
	TotalAvailable float64 `json:"total_available"`
	Grants         struct {
		Object string `json:"object"`
		Data   []struct {
			Object      string  `json:"object"`
			ID          string  `json:"id"`
			GrantAmount float64 `json:"grant_amount"`
			UsedAmount  float64 `json:"used_amount"`
			EffectiveAt float64 `json:"effective_at"`
			ExpiresAt   float64 `json:"expires_at"`
		} `json:"data"`
	} `json:"grants"`
}

func GetBalance() (Billing, error) {
	var data Billing
	url := "https://api.openai.com/dashboard/billing/credit_grants"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", viperConfig.GetString("openai_api_key")))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return data, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
