package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil" 
	"log"
	"os"
	"github.com/bwmarrin/discordgo"
)

type Configuration struct {
	BotToken      string
	ClienId 	  string
	ClientSecret  string

}

var botID string

func main() {
	config := loadConfig("conf.json")
	// Initialise Discord client
	fmt.Println("Initialising...")
	discord, err := discordgo.New("Bot " + string(config.BotToken))
	if err != nil {
		log.Panicln("Failed to create a discord session", err)
	}

	bot, err := discord.User("@me")
	if err != nil {
		log.Panicln("Could not access accout", err)
	}

	botID = bot.ID

	discord.AddHandler(scanAttachments)
	err = discord.Open()
	if err != nil {
		log.Println("Unable to establish connection", err)
	}

	defer discord.Close()

	//Hacky way to make program wait forever without loop
	<-make(chan struct{})
}

func scanAttachments(s*discordgo.Session, msg*discordgo.MessageCreate) {
	user := msg.Author
	if user.ID == botID || user.Bot {
		// Don't check the bots own messages
		return
	}

	attatchments := msg.Attachments
	if attatchments != nil {
		for _, atch := range attatchments {
			sha256, err := getSHA(atch.URL)
			if err != nil {
				log.Println("Problem determining SHA256 skipping...")
			}
			fmt.Println(sha256)
		}
	}
	

	content := msg.Content
	s.ChannelMessageSend(msg.ChannelID, content)
}

func getSHA(FileUrl string) (string, error) {
	resp, err := http.Get(FileUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write([]byte(body))
	hash := hex.EncodeToString(h.Sum(nil))

	return hash, err
}

func loadConfig(file string) Configuration {
	var config Configuration
	configFile, err := os.Open(file)
	defer configFile.Close()
	jsonDecoder := json.NewDecoder(configFile)
	jsonDecoder.Decode(&config)
	if err != nil {
		log.Panic("Problem opening config file", err)
	}

	return config
}
