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

func invokeMenuFuncs(opt int) {
	switch opt {
	case 1:
		iconConverter()
	case 2:
		showTime()
	case 3:
		openCalculator()
	default:
		fmt.Println("Error occurred while invoking menu func")
		return
	}
}

func iconConverter() {
	fmt.Println("icon converter")
}
func showTime() {
	fmt.Println("show time")
}
func openCalculator() {
	fmt.Println("open calculator")
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
				invokeMenuFuncs(selectedOption)
			}
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