package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
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

	dg.AddHandler(ready)
}

func ready(s *discordgo.Session){
	 err := s.UpdateStatus(0, "!talk")
	 if err != nil {
		 fmt.Println("error with status", err)
	 }
}


func newMessageEmbed(title string ,description string,color int, ) *discordgo.MessageEmbed {
	msg := discordgo.MessageEmbed{
		Title:title,
		Description: description,
		Color: color,
	}
	return &msg
}


var commandListsEmbed = newMessageEmbed(
	"Here is list of my commands!",
	"1 : Unicode converter \n 2 : Check time \n 3:  Open calculator \n Type !talk <option> to choose one",
	0000,
)

func sendCommands(s * discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSendEmbed(m.ChannelID,commandListsEmbed)
	if err != nil{
		fmt.Println("error occurred ::",err)
		return
	}
}



func getEnteredOption(option string,s *discordgo.Session, m *discordgo.MessageCreate ){
	index := strings.Index(option,"")
	if index == 0 {
		opt := option[index + 6:]
		if opt == "1" || opt == "2" || opt == "3"{
			msg  := "You have selected option:" + opt
			_,err := s.ChannelMessageSend(m.ChannelID,msg)
			if err != nil{
				fmt.Println("error occurred ::",err)
				return
			}
		} else {
			_,err := s.ChannelMessageSend(m.ChannelID,"Please enter correct number")
			if err != nil{
				fmt.Println("error occurred ::",err)
				return
			}
		}
	} else {
		return
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
		}
		if len(m.Content) >= 8 {
			_, err := s.ChannelMessageSend(m.ChannelID,"different command with prefix")
			if err != nil{
				fmt.Println("error occurred ::",err)
				return
			}
		}
	}

}