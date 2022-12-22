package kuebiko

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Bot paremeters
var (
	BotToken string
	GuildID  = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	///GuildID  = ""
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)
var s *discordgo.Session

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "basic-command",
		Description: "Kuebiko Responds!",
	},
	{
		Name:        "animelist",
		Description: "get AnimeList",
		// String option here ->
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "username",
				Description: "Type the username of the anime list",
				Required:    true,
			},
		},
	},
	{
		Name:        "コード",
		Description: "Bla",
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hey there! Congratulations, you just executed your first slash command",
			},
		})
	},
	"animelist": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}
		output := optionMap["username"]

		//margs := make([]interface{}, 0, len(options))
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: output.StringValue(),
			},
		})

	},
	"コード": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "たわけ者",
			},
		})

	},
}

func Run() {
	//fmt.Println("Got Tokens: ", BotToken)
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the handlers
	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	//dg.AddHandler(commandHandler)

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	//dg.AddHandler(messageCreate)

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

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	// Adds the commands
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
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
	s.UpdateGameStatus(0, "!golang")
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
