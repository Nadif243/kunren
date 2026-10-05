package main

import (
	"fmt"
	"kunnrenengine/internal/scraper"
	"log"
	"os"
	"os/signal"
	"strings"
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
	// This matches the "Message Content Intent" toggled in the Developer Portal.
	// Without this, Discord will censor the text of the messages sent to the bot.
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent

	// 5. Open the WebSocket connection
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection: ", err)
	}

	// 6. Keep the Server Alive
	fmt.Println("Renbun Engine (錬文) is online. Press CTRL-C to exit.")

	// Create a channel to listen for OS-level interrupt signals (like Ctrl+C).
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
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong! 錬文 engine is listening.")
		return
	}

	// !scrape [URL] command test
	if strings.HasPrefix(m.Content, "!scrape ") {
		// Extract the URL from the message
		url := strings.TrimSpace(strings.TrimPrefix(m.Content, "!scrape "))

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⚙️ Fetching and parsing HTML from: <%s>...", url))

		// Call our new internal package
		sentences, err := scraper.ExtractSentences(url)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ **Scrape Failed:**\n`%v`", err))
			return
		}

		if len(sentences) == 0 {
			s.ChannelMessageSend(m.ChannelID, "⚠️ No valid Japanese sentences found in <p> tags on this page.")
			return
		}

		// Preview the first 3 sentences found
		preview := fmt.Sprintf("✅ **Scrape Successful!** Found %d sentences.\n\n**Preview:**\n", len(sentences))
		for i := 0; i < 3 && i < len(sentences); i++ {
			preview += fmt.Sprintf("%d. %s\n", i+1, sentences[i])
		}

		s.ChannelMessageSend(m.ChannelID, preview)
	}
}
