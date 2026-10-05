package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load the secret token from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: No .env file found or error reading it.")
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("CRITICAL: DISCORD_TOKEN is missing!")
	}

	// 2. Initialize the Discord Session
	// The Discord API requires the exact prefix "Bot " before the token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}

	// 3. Register the Event Handler
	// This tells the Go runtime: "Every time a message is sent in the Discord server,
	// execute the 'messageCreate' function."
	dg.AddHandler(messageCreate)

	// 4. Declare Gateway Intents
	// This matches the "Message Content Intent" you toggled in the Developer Portal.
	// Without this, Discord will censor the text of the messages sent to your bot.
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent

	// 5. Open the WebSocket connection
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection: ", err)
	}

	// 6. Keep the Server Alive
	fmt.Println("Renbun Engine (錬文) is online. Press CTRL-C to exit.")

	// We create a channel to listen for OS-level interrupt signals (like Ctrl+C).
	// The main Goroutine completely freezes at `<-sc`, keeping the bot alive indefinitely.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// 7. Clean Shutdown
	fmt.Println("\nShutting down Renbun Engine safely...")
	dg.Close()
}

// This function is the callback for the event listener.
// discordgo automatically spawns a NEW Goroutine for this function every single time a message is received.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself to prevent infinite echoing loops.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// A simple ping-pong test to verify the WebSocket is receiving and transmitting.
	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong! 錬文 engine is listening.")
	}
}
