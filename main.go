package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	clitable "gopkg.in/benweidig/cli-table.v2"
)

type summoner struct {
	AccountID     string `json:"accountId"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int    `json:"revisionDate"`
	Name          string `json:"name"`
	ID            string `json:"id"`
	PUUID         string `json:"puuid"`
	SummonerLevel int    `json:"summonerLevel"`
}

type observer struct {
	EncryptionKey string `json:"encriptionKey"`
}

type perks struct {
	PerksID      []int `json:"perkIds"`
	PerkStyle    int   `json:"perkStyle"`
	PerkSubStyle int   `json:"perkSubStyle"`
}

type gameCustomizationObject struct {
	Category string `json:"category"`
	Content  string `json:"content"`
}

type currentGameParticipant struct {
	ChampionID               int                       `json:"championId"`
	Perks                    perks                     `json:"perks"`
	ProfileIconID            int                       `json:"profileIconId"`
	Bot                      bool                      `json:"bot"`
	TeamID                   int                       `json:"teamId"`
	SummonerName             string                    `json:"summonerName"`
	SummonerID               string                    `json:"summonerId"`
	SpellID                  int                       `json:"spellId"`
	Spell2ID                 int                       `json:"spell2Id"`
	GameCustomizationObjects []gameCustomizationObject `json:"gameCustomizationObjects"`
}

type bannedChampion struct {
	PickTurn   int `json:"pickTurn"`
	ChampionID int `json:"championId"`
	TeamID     int `json:"teamId"`
}

type currentGameInfo struct {
	GameID            int                      `json:"gameId"`
	GameType          string                   `json:"gameType"`
	GameStartTime     int                      `json:"gameStartTime"`
	MapID             int                      `json:"mapId"`
	GameLength        int                      `json:"gameLength"`
	PlatformID        string                   `json:"platformId"`
	GameMode          string                   `json:"gameMode"`
	BannedChampions   []bannedChampion         `json:"bannedChampions"`
	GameQueueConfigID int                      `json:"gameQueueConfigId"`
	Observers         observer                 `json:"observers"`
	Participants      []currentGameParticipant `json:"participants"`
}

func main() {
	var baseURL string = "https://la2.api.riotgames.com"
	var token string = os.Getenv("RIOT_TOKEN")

	client := &http.Client{}
	summonerName := os.Args[1]

	table := clitable.New()
	table.AddRow("ID", "Name")

	summoner := summoner{}
	currentGame := currentGameInfo{}

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

	json.NewDecoder(summonerResp.Body).Decode(&summoner)

	encryptedSummonerID := summoner.ID
	matchReq, err := http.NewRequest("GET", baseURL+"/lol/spectator/v4/active-games/by-summoner/"+encryptedSummonerID, nil)
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
		err := json.NewDecoder(matchResp.Body).Decode(&currentGame)
		if err != nil {
			panic(err)
		}
		for _, participant := range currentGame.Participants {
			table.AddRow(participant.SummonerID, participant.SummonerName)
		}
		fmt.Println("Match found 🥳")
		fmt.Println(table.String())
	}
}
