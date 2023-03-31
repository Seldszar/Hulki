package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/toml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/tidwall/gjson"
)

type Achievement struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Icon        string `json:"icon"`

	Achieved   bool  `json:"achieved"`
	UnlockedAt int64 `json:"unlockedAt"`
}

type State struct {
	Achievements []Achievement `json:"achievements"`
}

var (
	state State
)

func getSchemaForGame(key, appID, locale string) ([]byte, error) {
	resp, err := http.Get(
		fmt.Sprintf("https://api.steampowered.com/ISteamUserStats/GetSchemaForGame/v2/?key=%s&appid=%s&l=%s", key, appID, locale),
	)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return io.ReadAll(resp.Body)
	}

	return nil, nil
}

func getPlayerSummaries(key, steamID string) ([]byte, error) {
	resp, err := http.Get(
		fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=%s&steamids=%s", key, steamID),
	)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return io.ReadAll(resp.Body)
	}

	return nil, nil
}

func getPlayerAchievements(key, steamID, appID string) ([]byte, error) {
	resp, err := http.Get(
		fmt.Sprintf("https://api.steampowered.com/ISteamUserStats/GetPlayerAchievements/v1/?key=%s&steamid=%s&appid=%s", key, steamID, appID),
	)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return io.ReadAll(resp.Body)
	}

	return nil, nil
}

func startWebServer() error {
	http.Handle("/", http.FileServer(http.Dir("./web/dist")))

	http.HandleFunc("/api/state", func(w http.ResponseWriter, r *http.Request) {
		w.Header().
			Set("Access-Control-Allow-Origin", "*")

		w.Header().
			Set("Content-Type", "application/json")

		json.NewEncoder(w).
			Encode(state)
	})

	port := config.Int("port", 3000)

	log.Info().
		Msgf("Server is ready: http://localhost:%d", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func refreshAchievements() ([]Achievement, error) {
	key := config.String("key")
	steamID := config.String("steam_id")
	locale := config.String("locale")

	log.Debug().
		Str("steamid", steamID).
		Msg("Fetching player summaries...")

	playersBytes, err := getPlayerSummaries(key, steamID)

	if err != nil {
		log.Err(err).
			Str("steamid", steamID).
			Msg("An error occured while fetching player summaries")

		return nil, err
	}

	appID := gjson.GetBytes(playersBytes, "response.players.0.gameid").
		String()

	if appID == "" {
		return []Achievement{}, nil
	}

	log.Debug().
		Str("appid", appID).
		Msg("Fetching schema for game...")

	gameBytes, err := getSchemaForGame(key, appID, locale)

	if err != nil {
		log.Err(err).
			Str("appid", appID).
			Msg("An error occured while fetching schema for game")

		return nil, err
	}

	log.Debug().
		Str("appid", appID).
		Str("steamid", steamID).
		Msg("Fetching player achievements...")

	achievementsBytes, err := getPlayerAchievements(key, steamID, appID)

	if err != nil {
		log.Err(err).
			Str("appid", appID).
			Str("steamid", steamID).
			Msg("An error occured while fetching player achievements")

		return nil, err
	}

	res := make([]Achievement, 0)

	gjson.GetBytes(gameBytes, "game.availableGameStats.achievements").
		ForEach(func(key, value gjson.Result) bool {
			achievement := gjson.GetBytes(achievementsBytes, fmt.Sprintf(`playerstats.achievements.#(apiname=="%s")`, value.Get("name").String()))

			res = append(res, Achievement{
				Name:        value.Get("name").String(),
				DisplayName: value.Get("displayName").String(),
				Description: value.Get("description").String(),
				Icon:        value.Get("icon").String(),

				Achieved:   achievement.Get("achieved").Bool(),
				UnlockedAt: achievement.Get("unlocktime").Int(),
			})

			return true
		})

	return res, nil
}

func main() {
	config.AddDriver(toml.Driver)

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.SetGlobalLevel(zerolog.Level(config.Int("level", int(zerolog.InfoLevel))))

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stdout,
	})

	if err := config.LoadFiles("config.toml"); err != nil {
		log.Err(err).
			Msg("Unable to load configuration file")

		return
	}

	go startWebServer()

	for {
		if res, err := refreshAchievements(); err == nil {
			log.Debug().
				Interface("achievements", res).
				Msg("Achievements refreshed")

			state.Achievements = res
		}

		time.Sleep(time.Minute)
	}
}
