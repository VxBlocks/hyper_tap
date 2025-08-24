package handler

import (
	"fmt"
	"io"
	"net/http"
)

func main() {

	url := "https://api.coingecko.com/api/v3/global"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", "CG-SZjsFBHvHbyyf4v2Pu6HwoUD")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

}
