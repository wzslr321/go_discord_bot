package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.AddHandler(ready)

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	closeErr := dg.Close()
	if closeErr != nil {
		fmt.Println("error closing connection", closeErr)
		return
	}

}

func ready(s *discordgo.Session, event *discordgo.Ready){
	 s.UpdateStatus(0, "!talk")
}


func newMessageEmbed(title string ,description string ) *discordgo.MessageEmbed {
	msg := discordgo.MessageEmbed{
		Title:title,
		Description: description,
	}
	return &msg
}


var commandListsEmbed = newMessageEmbed(
	"Here is list of my commands!",
	"1 : Embed message creator \n 2 : Check time \n 3:  Open calculator \n Type !talk <option> to choose one",
)

func sendCommands(s * discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSendEmbed(m.ChannelID,commandListsEmbed)
	if err != nil{
		fmt.Println("error occurred ::",err)
		return
	}
}

var isMenuNumberTrue bool
var selectedOption int

func getEnteredOption(option string,s *discordgo.Session, m *discordgo.MessageCreate ) {

	index := strings.Index(option,"")
	options := [3]string{"1","2","3"}
	opt := option[index + 6:]
	for _, opts := range options {
		if opt == opts{
			msg  := "You have selected option:" + opt
			optNum, _ := strconv.Atoi(opt)
			selectedOption = optNum
			_,err := s.ChannelMessageSend(m.ChannelID,msg)
			if err != nil{
				fmt.Println("error occurred ::",err)
				return
			}
			isMenuNumberTrue = true
			return
		}
	}
	isMenuNumberTrue = false
	return
}

func invokeMenuFuncs(opt int, s *discordgo.Session, m *discordgo.MessageCreate) {
	switch opt {
	case 1:
		embedInstruction(s,m)
	case 2:
		showTime(s,m)
	case 3:
		openCalculator(s,m)
	default:
		fmt.Println("Error occurred while invoking menu func")
		return
	}
}

func embedInstruction(s *discordgo.Session, m *discordgo.MessageCreate) {

	howToUseEmbed := newMessageEmbed(
		"How to create message embed",
		"Enter message starting with @ which contains title and description of message separated by '|' \n \n" +
			"Example: !talk @ cool title | cool description ",
	)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID,howToUseEmbed)
	if err != nil{
		fmt.Println("error occurred ::",err)
		return
	}
}
func showTime(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("show time")
}
func openCalculator(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("open calculator")
}

func createEmbed(s *discordgo.Session, m *discordgo.MessageCreate, title string, description string) {
	if len(title) >= 1  && len(description) >= 1 {
		createdMsgEmbed := newMessageEmbed(
			title,
			description,
		)

		_, err := s.ChannelMessageSendEmbed(m.ChannelID,createdMsgEmbed)
		if err != nil{
			fmt.Println("error occurred ::",err)
			return
		}
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {


	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!talk" {
		sendCommands(s,m)
	}

	if strings.HasPrefix(m.Content, "!talk") && len(m.Content) != 5{
		if len(m.Content) == 7{
			getEnteredOption(m.Content,s,m)
			if  isMenuNumberTrue == false {
				_, err := s.ChannelMessageSend(m.ChannelID,"Please enter correct number")
				if err != nil{
					fmt.Println("error occurred ::",err)
					return
				}
			} else {
				invokeMenuFuncs(selectedOption,s,m)
			}
		}

		if len(m.Content) >= 8 {
			msg := m.Content
			if strings.HasPrefix(m.Content,"!talk @"){
				indexDesc := strings.Index(msg,"|")
				indexTitle := strings.Index(msg,"@")
				if indexDesc > -1 && indexTitle > -1{
					title := msg[indexTitle+1:indexDesc]
					description := msg[indexDesc+1:]
					createEmbed(s,m,title,description)
				} else {
					fmt.Println("error with msg embed")
					return
				}
			}

			_, err := s.ChannelMessageSend(m.ChannelID,"different command with prefix")
			if err != nil{
				fmt.Println("error occurred ::",err)
				return
			}
		}
	}

}