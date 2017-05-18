package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"meanbot/constants"
	"meanbot/settings"
	"os"
	"strconv"
	"strings"
	"tgbot/bot"
	"tgbot/methods"
	"tgbot/tgtype"
	"time"
)

var meanBot bot.Bot
var insults []string

func main() {
	botSettings := settings.GetSavedSettings()

	insults = loadInsults(constants.BotPath + "insults.dat")
	fmt.Println("Installed insults:")
	for _, v := range insults {
		fmt.Println("\t", v)
	}

	meanBot = bot.Bot{APIKey: constants.BotAPIKey}

	getMeMethod := methods.GetMe{}
	meanBotID := getMeMethod.CallMethod(meanBot.GetBotURL()).(tgtype.User).ID

	fmt.Println("Bot ID:", meanBotID)

	//
	//Configure the command processor
	//
	meanBot.CommandProcessor.BotName = constants.BotUserName
	meanBot.CommandProcessor.SwallowRegisteredBotCommands = true
	meanBot.CommandProcessor.SwallowOtherCommands = true

	// /set setting value
	meanBot.CommandProcessor.RegisterCommad("set", func(args string, m tgtype.Message) {
		arguments := strings.Fields(args)

		method := methods.SendMessage{
			ChatID:           m.Chat.ID,
			ReplyToMessageID: m.MessageID,
		}

		if len(arguments) < 2 {
			method.Text = fmt.Sprintf(constants.ReplyNotEnoughArguments)
			method.CallMethod(meanBot.GetBotURL())
			return
		}

		setting := botSettings[getKey(m.Chat.ID)]
		switch arguments[0] {
		case "insult_interval":
			if val, err := strconv.ParseInt(arguments[1], 10, 64); err == nil {
				if val < 1 {
					val = 1
				}
				setting.InsultInterval = val
			} else {
				method.Text = fmt.Sprintf(constants.ReplyNotANumber, arguments[1])
				method.CallMethod(meanBot.GetBotURL())
				return
			}
		default:
			method.Text = fmt.Sprintf(constants.ReplyNotASetting, arguments[0])
			method.CallMethod(meanBot.GetBotURL())
			return
		}
		botSettings[getKey(m.Chat.ID)] = setting

		method.Text = fmt.Sprintf(constants.ReplySettingSet, arguments[0], arguments[1])
		method.CallMethod(meanBot.GetBotURL())

		settings.SaveSettings(botSettings)
	})

	// /insult
	meanBot.CommandProcessor.RegisterCommad("insult", func(_ string, m tgtype.Message) {
		if m.ReplyToMessage == nil {
			method := methods.SendMessage{
				ChatID:           m.Chat.ID,
				ReplyToMessageID: m.MessageID,
				Text:             constants.ReplyMissingTarget,
			}

			method.CallMethod(meanBot.GetBotURL())
			return
		}

		if strings.ToLower(m.ReplyToMessage.From.Username) == constants.BotUserName {
			method := methods.SendMessage{
				ChatID:           m.Chat.ID,
				ReplyToMessageID: m.MessageID,
				Text:             constants.ReplyInsultItself,
			}

			method.CallMethod(meanBot.GetBotURL())
			return
		}

		sendInsult(m.ReplyToMessage)
	})

	// /suggest suggestion
	meanBot.CommandProcessor.RegisterCommad("suggest", func(args string, m tgtype.Message) {
		if strings.TrimSpace(args) == "" {
			//Tell the user to fuck off
			method := methods.SendMessage{
				ChatID:           m.Chat.ID,
				ReplyToMessageID: m.MessageID,
				Text:             constants.ReplyMissingSuggestion,
			}

			method.CallMethod(meanBot.GetBotURL())
			return
		}

		//Save the suggestion to a file
		file, fileErr := os.OpenFile(constants.BotPath+"suggestions.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		defer file.Close()

		if fileErr != nil {
			fmt.Println("Failed to open suggestions file!")
		}

		file.WriteString(args + "\n")

		//Send a lovely thanks message
		method := methods.SendMessage{
			ChatID:           m.Chat.ID,
			ReplyToMessageID: m.MessageID,
			Text:             fmt.Sprintf(constants.ReplyAcceptSuggestion, m.From.FirstName),
		}

		method.CallMethod(meanBot.GetBotURL())
	})

	// Other targeted command
	meanBot.CommandProcessor.OnUnregisteredTragetedCommand = func(command string, _ string, m tgtype.Message) {
		// Tell the user that such message doesn't exist
		method := methods.SendMessage{
			ChatID:           m.Chat.ID,
			ReplyToMessageID: m.MessageID,
			Text:             fmt.Sprintf(constants.ReplyNonexistantCommand, command),
		}

		method.CallMethod(meanBot.GetBotURL())
	}

	//Log the recieved commands
	meanBot.CommandProcessor.OnBeforeCommand = func(m tgtype.Message) {
		fmt.Printf("%v Recieved command \"%v\" from %v %v (%v/%v) in %v chat %v %v (%v/%v)\n",
			time.Now().UTC(),
			m.Text,
			m.From.FirstName,
			m.From.LastName,
			m.From.ID,
			m.From.Username,
			m.Chat.TypeString,
			m.Chat.FirstName,
			m.Chat.LastName,
			m.Chat.ID,
			m.Chat.Username)
	}

	//
	//Configure bot callbacks
	//
	meanBot.OnMessage = func(message *tgtype.Message) {
		// If a current group doesn't have settings generate them and save them
		if _, ok := botSettings[getKey(message.Chat.ID)]; !ok {
			botSettings[getKey(message.Chat.ID)] = settings.GetDefaultSettings()
			fmt.Println(botSettings)
			settings.SaveSettings(botSettings)
		}

		if message.NewChatMember != nil {
			if message.NewChatMember.ID == meanBotID {
				//The bot was added to a new group
				onAddToGroup(message)
			} else {
				//Possible TODO: Welcome message
				//Not really sure if its necessary as a lot of bots have them and it gets a bit annoying to receive what is essentially spam whenever a new user joins
				//(Maybe a setting for this?)
			}
		} else if message.Text != "" {
			onTimeForInsult(botSettings[getKey(message.Chat.ID)], message)
		}
	}

	meanBot.RunBot()
}

func loadInsults(location string) []string {
	insultsFile, insultsFileErr := os.Open(location)
	defer insultsFile.Close()

	if insultsFileErr != nil {
		fmt.Println("Couldn't open the insults file located at ", location)
	}

	//Scanner splits by the line by default
	scanner := bufio.NewScanner(insultsFile)
	scanner.Scan()
	insultCount, insultCountErr := strconv.ParseInt(scanner.Text(), 10, 64)

	if insultCountErr != nil {
		fmt.Printf("failed to convert %v to a number", scanner.Text())
	}

	result := make([]string, insultCount)

	var i int64
	for scanner.Scan() {
		if scanner.Text() == "\n" {
			continue
		}

		if i >= insultCount {
			fmt.Println("The insult count in the file doesn't match the count specified on top")
		}

		result[i] = scanner.Text()

		i++
	}

	if i < insultCount-1 {
		fmt.Println("The insult count in the the file is smaller than the one specified in it")
	}

	return result
}

func getKey(i int64) string {
	return strconv.FormatInt(i, 10)
}

func onAddToGroup(m *tgtype.Message) {
	joinmessage := fmt.Sprintf(constants.ReplyJoinChat, m.Chat.TypeString)
	method := methods.SendMessage{
		ChatID: m.Chat.ID,
		Text:   joinmessage,
	}
	method.CallMethod(meanBot.GetBotURL())

	fmt.Printf("%v bot as added to %v chat %v %v (%v/%v)\n",
		time.Now().UTC(),
		m.Chat.TypeString,
		m.Chat.FirstName,
		m.Chat.LastName,
		m.Chat.ID,
		m.Chat.Username)
}

func sendInsult(m *tgtype.Message) {
	insult := insults[rand.Intn(len(insults))]
	method := methods.SendMessage{
		ChatID:           m.Chat.ID,
		Text:             insult,
		ReplyToMessageID: m.MessageID,
	}

	method.CallMethod(meanBot.GetBotURL())

	fmt.Printf("%v: insulted %v %v (%v/%v) in %v chat %v (%v/%v) with \"%v\" in response to \"%v\"\n",
		time.Now().UTC(),
		m.From.FirstName,
		m.From.LastName,
		m.From.ID,
		m.From.Username,
		m.Chat.TypeString,
		m.Chat.FirstName,
		m.Chat.ID,
		m.Chat.Username,
		insult,
		m.Text)
}

func onTimeForInsult(s settings.Settings, m *tgtype.Message) {
	var shouldInsult bool

	switch m.Chat.TypeString {
	case "private":
		shouldInsult = true
	case "group", "supergroup":
		insultInterval := s.InsultInterval

		if insultInterval <= 1 || rand.Int63n(insultInterval) == 0 {
			shouldInsult = true
		}
	}

	if shouldInsult {
		sendInsult(m)
	}
}
