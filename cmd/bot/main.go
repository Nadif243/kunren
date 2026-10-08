package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"rennbunengine/internal/db"
	"rennbunengine/internal/worker"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load the secret token from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: No .env file found.")
	}

	// 2. Ignite the Crucible
	db.InitDB("./renbun.db")
	defer db.DB.Close()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("CRITICAL: DISCORD_TOKEN is missing!")
	}

	// 3. Initialize the Discord Session
	// The Discord API requires the exact prefix "Bot " before the token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}

	// 4. Register the Event Handler
	// This tells the Go runtime: "Every time a message is sent in the Discord server,
	// execute the 'messageCreate' function."
	dg.AddHandler(messageCreate)

	// 5. Declare Gateway Intents
	// This matches the "Message Content Intent" toggled in the Developer Portal.
	// Without this, Discord will censor the text of the messages sent to the bot.
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent

	// 6. Open the WebSocket connection
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection: ", err)
	}

	// Set the Alchemy Custom Status
	err = dg.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Type:  discordgo.ActivityTypeCustom,
				Name:  "Custom Status",        // Discord API requires this field to be populated
				State: "錬金術駆動の没入型エンジンで言葉を錬成中", //Refining text through an alchemy-driven immersion engine
			},
		},
	})
	if err != nil {
		fmt.Println("Warning: Could not set custom status:", err)
	}

	// 7. Keep the Server Alive
	fmt.Println("Renbun Engine (錬文) is online. Press CTRL-C to exit.")

	// Create a channel to listen for OS-level interrupt signals (like Ctrl+C).
	// The main Goroutine completely freezes at `<-sc`, keeping the bot alive indefinitely.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// 8. Clean Shutdown
	fmt.Println("\nShutting down Renbun Engine safely...")
	dg.Close()
}

// Helper function to read URLs from our text file
func loadTargets(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && strings.HasPrefix(line, "http") {
			urls = append(urls, line)
		}
	}
	return urls
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

	// !mine [keyword] command
	if strings.HasPrefix(m.Content, "!mine ") {
		// Extract the URL from the message
		keyword := strings.TrimSpace(strings.TrimPrefix(m.Content, "!mine "))

		// 1. Load the target list
		urls := loadTargets("target.txt")
		if len(urls) == 0 {
			s.ChannelMessageSend(m.ChannelID, "⚠️ `target.txt` is empty or missing.")
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("🔍 **Mining %d sources** for the word: `%s`...", len(urls), keyword))

		// 2. Ignite the Worker Pool (using 5 concurrent workers)
		results := worker.RunPool(urls, 5)

		// 3. The Goldilocks Filter
		var validSentences []string
		for _, res := range results {
			if res.Err != nil {
				continue // Silently ignore dead websites
			}
			for _, sentence := range res.Sentences {
				// Must contain the exact keyword
				if strings.Contains(sentence, keyword) {
					runeCount := len([]rune(sentence))
					// Must be longer than 15 chars (destroys headers/dates)
					// Must be shorter than 60 chars (destroys run-on paragraphs)
					if runeCount >= 15 && runeCount <= 60 {
						validSentences = append(validSentences, sentence)
					}
				}
			}
		}

		// 4. Handle No Results
		if len(validSentences) == 0 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⚠️ No high-quality sentences found for `%s` matching the length criteria.", keyword))
			return
		}

		// 5. Shuffle the array to ensure fresh examples every search attempt
		rand.Shuffle(len(validSentences), func(i, j int) {
			validSentences[i], validSentences[j] = validSentences[j], validSentences[i]
		})

		// Limit output to the top 3
		limit := 3
		if len(validSentences) < limit {
			limit = len(validSentences)
		}

		// 6. Format and Ship the Output as a Rich Embed
		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("錬文 Engine: %s", keyword),
			Description: fmt.Sprintf("Scanned **%d** sources.\nFound **%d** high-quality sentences.", len(urls), len(validSentences)),
			// Set the color to a sleek gold/amber (Hex: #F5A623)
			Color:  0xF5A623,
			Fields: []*discordgo.MessageEmbedField{},
		}

		for i := 0; i < limit; i++ {
			// Replace the target word with an inline code block so it has a distinct background
			highlighted := strings.ReplaceAll(validSentences[i], keyword, fmt.Sprintf("`%s`", keyword))

			// Append each sentence as its own distinct block in the UI
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("Result %d", i+1),
				Value:  highlighted,
				Inline: false,
			})
		}

		// Using SendEmbed instead of Send
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}

	// The Binding Command
	if m.Content == "!bind" {
		response, err := db.BindUser(m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⚠️ Alchemical failure during binding: "+err.Error())
			return
		}
		s.ChannelMessageSend(m.ChannelID, "🩸 "+response)
		return
	}

	// The Vault Creation Command (!vault create <name>)
	if strings.HasPrefix(m.Content, "!vault create ") {
		vaultName := strings.TrimSpace(strings.TrimPrefix(m.Content, "!vault create "))
		if vaultName == "" {
			s.ChannelMessageSend(m.ChannelID, "⚠️ You must provide a name for the vault. (e.g., `!vault create Technical Vault`)")
			return
		}

		response, err := db.CreateVault(m.Author.ID, vaultName)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⚠️ Failed to structure the vault: "+err.Error())
			return
		}
		s.ChannelMessageSend(m.ChannelID, "📖 "+response)
		return
	}
}
