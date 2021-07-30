package main

import (
	"log"
	"os"
	"regexp"
	"time"

	"github.com/DisgoOrg/disgohook"
	"github.com/joho/godotenv"
	"github.com/shirou/gopsutil/v3/process"
)

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

var (
	hook, err = disgohook.NewWebhookClientByToken(nil, nil, goDotEnvVariable("HOOK"))
	placeId   string
	reset     = false
)

type MarketPlaceInfo struct { // https://mholt.github.io/json-to-go/
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Creator     struct {
		ID              int    `json:"Id"`
		Name            string `json:"Name"`
		CreatorType     string `json:"CreatorType"`
		CreatorTargetID int    `json:"CreatorTargetId"`
	} `json:"Creator"`
	IconImageAssetID int64 `json:"IconImageAssetId"`
}

func GetProcessByName(targetProcessName string) *process.Process {
	processes, _ := process.Processes()
	for _, proc := range processes {
		name, _ := proc.Name()
		if name == targetProcessName {
			return proc
		}
	}
	return nil
}

func UpdateRobloxPresence() {
	roblox := GetProcessByName("RobloxPlayerBeta.exe")
	for roblox == nil {
		roblox = GetProcessByName("RobloxPlayerBeta.exe")
		if reset == false {
			reset = true
		}
	}

	reset = false
	args, _ := roblox.Cmdline()
	placePattern := regexp.MustCompile(`placeId=(\d+)`)
	placeMatch := placePattern.FindStringSubmatch(args)[1]
	if placeMatch != placeId {
		placeId = placeMatch
		_, _ = hook.SendContent("`Started playing:`\nhttps://www.roblox.com/games/" + placeId + "/-")
	}
}

func main() {
	for {
		UpdateRobloxPresence()
		time.Sleep(time.Second * 3)
	}
}
