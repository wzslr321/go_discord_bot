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
	"time"
)

var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

//  Whole func 'main' contains basic bot configuration, responsible for start,close,status and messageListener
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

//  Sets bot status to '!talk'
func ready(s *discordgo.Session, event *discordgo.Ready){
	 s.UpdateStatus(0, "!talk")
}
/*
	messageEmbed constructor with MessageEmbed type, gives a possibility to
 	use create var and use it  in s.ChannelMessageSendEmbed function

 */
func newMessageEmbed(title string ,description string ) *discordgo.MessageEmbed {
	msg := discordgo.MessageEmbed{
		Title:title,
		Description: description,
	}
	return &msg
}

//   Creates messageEmbed and send it to channel where it was invoked by user`s message
func sendCommands(s * discordgo.Session, m *discordgo.MessageCreate) {
	 commandListsEmbed := newMessageEmbed(
		"Here is list of my commands!",
		"1 : Embed message creator \n 2 : Check time \n 3:  Open calculator \n Type !talk <option> to choose one",
	)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID,commandListsEmbed)
	if err != nil{
		fmt.Println("error occurred ::",err)
		return
	}
}

//   Global vars used as validation in message listener
var isMenuNumberTrue bool
var selectedOption int

/*
     Function which loops for every possible option and checks if
	 user entered correct option, sends message on channel if true

 */
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

//  Invokes one of func based on what menu option user have entered
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

//  Sends instruction of how to create EmbedMessage
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

//  Shows actual time in readable format
func showTime(s *discordgo.Session, m *discordgo.MessageCreate) {
	showTimeMsg := time.Now().Format("2006-01-02 15:04:05")
	_, err := s.ChannelMessageSend(m.ChannelID,"Actual time is : " + showTimeMsg)
	if err != nil{
		fmt.Println("error occurred ::",err)
		return
	}
}

//  Sends instruction of how to use calculator, same way as it was in embedInstruction func
func openCalculator(s *discordgo.Session, m *discordgo.MessageCreate) {
	howToUseCalculator := newMessageEmbed(
		"How to use calculator",
		"Just enter calculation, starting with % \n \n" +
			"Possible calculations: '+' , '-' , '*' , ':' \n \n" +
			"Example: !talk % 2+2 ",
	)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID,howToUseCalculator)
	if err != nil{
		fmt.Println("error occurred ::",err)
		return
	}
}

/*
	Function responsible for creating embed message based on user`s entered data
	It checks if user entered at least one sign in title and description,
	sends created embed message to channel where was invoked

 */
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

/*
	Prototype of calculation function without final calculation based on sign
	It operates on floats to make it possible to have more precised results
	Checks if result is integer by checking decimal numbers and
	display result with different precision based on it

 */
func calcFuncProto(s *discordgo.Session, m *discordgo.MessageCreate, calcResult float64) {
	convCalcResult := strconv.FormatFloat(calcResult,'f',4,32)
	calcIndex := strings.Index(convCalcResult, ".")
	if calcIndex > -1 {
		isEvenIndex := convCalcResult[calcIndex+1:calcIndex+5]
		if isEvenIndex != "0000" {
			convCalcResult := strconv.FormatFloat(calcResult,'f',8,32)
			_, err := s.ChannelMessageSend(m.ChannelID,convCalcResult)
			if err != nil{
				fmt.Println("error occurred ::",err)
				return
			}
			return
		} else {
			convCalcResult := strconv.FormatFloat(calcResult,'f',0,32)
			_, err := s.ChannelMessageSend(m.ChannelID,convCalcResult)
			if err != nil{
				fmt.Println("error occurred ::",err)
				return
			}
			return
		}
	}
	_, err := s.ChannelMessageSend(m.ChannelID,convCalcResult)
	if err != nil{
		fmt.Println("error occurred ::",err)
		return
	}
}

/*
	4 functions below calculate result with proper calculation sign
	and pass it as argument to calcFuncProto

 */
func addValues(s *discordgo.Session, m *discordgo.MessageCreate,firstVal float64, secondVal float64) {
	calculationResult := firstVal + secondVal
	calcFuncProto(s,m,calculationResult)
}

func subtractValues(s *discordgo.Session, m *discordgo.MessageCreate,firstVal float64, secondVal float64)  {
	calculationResult := firstVal - secondVal
	calcFuncProto(s,m,calculationResult)
}

func multiplyValues(s *discordgo.Session, m *discordgo.MessageCreate,firstVal float64, secondVal float64)  {
	calculationResult := firstVal * secondVal
	calcFuncProto(s,m,calculationResult)
}

func divideValues(s *discordgo.Session, m *discordgo.MessageCreate,firstVal float64, secondVal float64) {
	calculationResult := firstVal / secondVal
	calcFuncProto(s,m,calculationResult)
}

/*
	calculate function checks what sign have user entered
	It loops for every calculation option and checks
	if it is equal to calculation sign that user have entered
	Saves user sign properly only if there is no whitespace between sign
	because of strings.Index being precised to one character
	It also separates values using this index, coverts them
	into float number and pass as argument to correct function

 */
func calculate(s *discordgo.Session, m *discordgo.MessageCreate) {
	calcOptions := []string{"+","-","*",":"}

	msg := m.Content

	for _,calcOption:= range calcOptions {
		firstValIndex := strings.Index(msg, "%")
		secondValIndex := strings.Index(msg, calcOption)
		if secondValIndex != -1 && firstValIndex != -1 {
			firstVal := msg[firstValIndex+2 : secondValIndex]
			convFirstVal, _ := strconv.ParseFloat(firstVal, 64)
			secondVal := msg[secondValIndex+1:]
			convSecondVal, _ := strconv.ParseFloat(secondVal, 64)
			switch calcOption {
			case "+":
				addValues(s, m, convFirstVal, convSecondVal)
			case "-":
				subtractValues(s, m, convFirstVal, convSecondVal)
			case "*":
				multiplyValues(s, m, convFirstVal, convSecondVal)
			case ":":
				divideValues(s, m, convFirstVal, convSecondVal)
			}
		}
	}
}

/*
	Message event listener, which is invoked everytime user enters message
	Uses 'strings' library to check if message contains prefix
	It also checks length or message to make it possible to use menu
	without invoking other functions

 */
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
			if strings.HasPrefix(msg,"!talk @"){
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

			if strings.HasPrefix(msg ,"!talk %"){
				calculate(s,m)
			}
		}
	}
}