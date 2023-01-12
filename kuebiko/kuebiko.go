package kuebiko

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/jikan-go"
)

// Bot paremeters
var (
	BotToken string
	GuildID  = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	///GuildID  = ""
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

//var s *discordgo.Session

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "anime",
		Description: "Get the Anime from MyAnimeList",
		// String option here ->
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "Type the anime name",
				Required:    true,
			},
		},
	},
	{
		Name:        "kuebiko",
		Description: "Get senpai to notice you.",
	},
	{
		Name:        "random-anime",
		Description: "Get a random anime!",
	},
	//https://media.tenor.com/NL_qdSSs28EAAAAd/eevee-dancing.gif
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"anime": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		output := optionMap["name"]
		//fmt.Printf(output.Name)

		query := url.Values{}
		fmt.Println("Input: ", output.StringValue())
		query.Set("q", output.StringValue())
		//query.Set("type", "tv")

		// Search anime
		search, err := jikan.GetAnimeSearch(query)
		if err != nil {
			panic(err)
		}
		animeData := search.Data[0]
		animeYear := fmt.Sprintf("%s (%d)", strings.Title(animeData.Season), animeData.Year)
		animeTitle := fmt.Sprintf("%s (%s)", animeData.TitleEnglish, (animeData.TitleJapanese))
		animeGenre := formatGenre(animeData.Genres)
		animeRanking := fmt.Sprintf("#%s", strconv.Itoa(animeData.Rank))
		animeSym := word_limiter(animeData.Synopsis, 180)
		fmt.Printf(animeSym)
		//animeDesc := animeData.Background
		fmt.Printf(animeGenre)

		//animeYear := `${animeData.Season}, animeData.Year`

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Author: &discordgo.MessageEmbedAuthor{},
						Title:  animeTitle,
						Color:  0xb004b2, // Green
						/*
							https://github.com/cadecuddy/sauce/blob/main/utils/print.go
						*/
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "❓ Type⤵",
								Value:  animeData.Type,
								Inline: true,
							},
							{
								Name:   "🕐 Episodes⤵",
								Value:  strconv.Itoa(animeData.Episodes),
								Inline: true,
							},
							{
								Name:   "💡 Status⤵",
								Value:  animeData.Status,
								Inline: true,
							},
							{
								Name:   "🎥 Studios⤵",
								Value:  animeData.Studios[0].Name,
								Inline: true,
							},
							{
								Name:   "📅 Year⤵",
								Value:  animeYear,
								Inline: true,
							},
							{
								Name:   "📈 Score⤵",
								Value:  strconv.FormatFloat(animeData.Score, 'g', 5, 64),
								Inline: true,
							},
							{
								Name:   "🏆 Ranking⤵",
								Value:  animeRanking,
								Inline: true,
							},
							{
								Name:   "📕 Source⤵",
								Value:  animeData.Source,
								Inline: true,
							},
							{
								Name:   "📜 Genres⤵",
								Value:  animeGenre,
								Inline: true,
							},
							{

								Name:   "Description⤵",
								Value:  animeData.Synopsis,
								Inline: false,
							},
						},
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: animeData.Images.Jpg.LargeImageUrl,
						},
					},
				},
			},
		})

	},
	"kuebiko": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "https://media.tenor.com/NL_qdSSs28EAAAAd/eevee-dancing.gif",
			},
		})
	},
	"random-anime": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		// Search anime
		getRandom, err := jikan.GetRandomAnime()
		if err != nil {
			panic(err)
		}
		fmt.Println(getRandom.Data.Title)
		gotRandom := getRandom.Data

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Author:      &discordgo.MessageEmbedAuthor{},
						Title:       gotRandom.Title,
						Color:       0xb004b2, // Green
						Description: gotRandom.Synopsis,
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: gotRandom.Images.Jpg.LargeImageUrl,
						},
					},
				},
			},
		})
	},
}

func formatGenre(genres []jikan.MalItem) string {
	var result string

	for i, genre := range genres {
		if i != 0 {
			//result += genre.Name
			result += fmt.Sprintf(", %s", genre.Name)
			//result += (", %s", genre.Name)
			//result := (", ", genre.Name)
			//o:= (", ", genre.Name)
		} else {
			result += genre.Name
		}
	}

	//r := fmt.Sprintf("%s", o)
	return result
	//return out
}
func word_limiter(s string, limit int) string {

	if strings.TrimSpace(s) == "" {
		return s
	}

	// convert string to slice
	strSlice := strings.Fields(s)

	// count the number of words
	numWords := len(strSlice)

	var result string

	if numWords > limit {
		// convert slice/array back to string
		result = strings.Join(strSlice[0:limit], " ")
		result = result + "..."
	} else {

		result = s
	}

	return string(result)

}

func Run() {
	dg, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the handlers
	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	// Add the intents for the bot!
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Printf("Logged in as: %v#%v", dg.State.User.Username, dg.State.User.Discriminator)
	fmt.Println("Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
func ready(s *discordgo.Session, event *discordgo.Ready) {
	/*
		log.Println("Adding commands...")
		registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
		// Adds the commands
		for i, v := range commands {
			cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
			registeredCommands[i] = cmd
		}*/
	/*
		if *RemoveCommands {
			log.Println("Removing commands...")

			for _, v := range registeredCommands {
				err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
				if err != nil {
					log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
				}
			}
		}*/
	// Set the playing status.
	s.UpdateGameStatus(2, "!golang")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//fmt.Println(m.Content)
	// Don't respond to himself
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "!golang") {
		s.ChannelMessageSend(m.ChannelID, "Hey! I'm using golang!")
	}

}
