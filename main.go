package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	var baseURL string = "https://la2.api.riotgames.com"
	var token string = os.Getenv("RIOT_TOKEN")

	client := &http.Client{}
	summonerName := os.Args[1]

	summonerReq, err := http.NewRequest("GET", baseURL+"/lol/summoner/v4/summoners/by-name/"+summonerName, nil)
	if err != nil {
		panic(err)
	}
	summonerReq.Header.Add("X-Riot-Token", token)
	summonerResp, err := client.Do(summonerReq)
	if err != nil {
		panic(err)
	}
	defer summonerResp.Body.Close()

	var summonerResult map[string]interface{}
	json.NewDecoder(summonerResp.Body).Decode(&summonerResult)

	var encryptedSummonerId string = summonerResult["id"].(string)
	matchReq, err := http.NewRequest("GET", baseURL+"/lol/spectator/v4/active-games/by-summoner/"+encryptedSummonerId, nil)
	if err != nil {
		panic(err)
	}
	matchReq.Header.Add("X-Riot-Token", token)
	matchResp, err := client.Do(matchReq)
	if err != nil {
		panic(err)
	}
	defer matchResp.Body.Close()

	if matchResp.StatusCode == 404 {
		fmt.Println("Match not found 😔")
	} else {
		var matchResult map[string]interface{}
		json.NewDecoder(matchResp.Body).Decode(&matchResult)
		fmt.Println(matchResult)
	}
}
