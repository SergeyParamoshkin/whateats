package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SergeyParamoshkin/version"
	products "github.com/SergeyParamoshkin/whateats/products"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/proxy"
)

func main() {

	var (
		showVersion = flag.Bool("version", false, "Show version information and exit")
		configPath  = flag.String("c", "", "Used for set path to config file.")
		proxyStr    = flag.String("proxy", getEnv("TELEGRAM_PROXY", ""),
			`socks5 proxy example -proxy="103.250.157.43:6667"`)
		proxyUser = flag.String("proxy_user", getEnv("TELEGRAM_PROXY_USER", ""),
			`socks5 proxy user -proxy_user="antip"`)
		proxyPassword = flag.String("proxy_password", getEnv("TELEGRAM_PROXY_PASSWORD", ""),
			`socks5 proxy password -proxy_password="forever1"`)

		token = flag.String("token", getEnv("TELEGRAM_TOKEN", ""), "bot token")
	)
	flag.Parse()

	if *showVersion {
		fmt.Println(version.ShowVersion())
		os.Exit(0)
	}

	if len(*configPath) > 0 {
		_, err := Configure(*configPath)
		if err != nil {
			log.Println(err)
		}
	}

	client := &http.Client{}
	auth := proxy.Auth{}
	if len(*proxyStr) > 0 {
		log.Println("run create proxy", *proxyStr)

		if len(*proxyUser) > 0 && len(*proxyPassword) > 0 {
			auth.User = *proxyUser
			auth.Password = *proxyPassword
		} else {
			log.Fatalf("user %s or password %s is empty", *proxyUser, *proxyPassword)
		}

		dialer, err := proxy.SOCKS5("tcp", *proxyStr, &auth, proxy.Direct)

		if err != nil {
			fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		}

		tr := &http.Transport{Dial: dialer.Dial}
		client.Transport = tr

	}
	bot, err := tgbotapi.NewBotAPIWithClient(*token, client)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
	}

	for update := range updates {
		var msg tgbotapi.MessageConfig
		m := products.NewMenu(update.Message.From.UserName)
		msg.ReplyMarkup = getWeekKeyboard()
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Привет "+update.Message.Chat.FirstName)
				msg.ReplyMarkup = getWeekKeyboard()
			case "open":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Выбирай день")
				msg.ReplyMarkup = getWeekKeyboard()
			case "close":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "В слудующий раз")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			}
		}

		switch update.Message.Text {

		case "понедельник", "monday":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, m.DayOfWeek("monday").GetAll())
		case "вторник", "tuesday":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, m.DayOfWeek("tuesday").GetAll())
		case "среда", "wednesday":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, m.DayOfWeek("wednesday").GetAll())
		case "четверг", "thursday":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, m.DayOfWeek("thursday").GetAll())
		case "пятница", "friday":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, m.DayOfWeek("friday").GetAll())
		case "суббота", "saturday":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, m.DayOfWeek("saturday").GetAll())
		case "воскресенье", "sunday":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, m.DayOfWeek("sunday").GetAll())
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "лучше напиши @sergeyparamoshkin 🤪")
		}
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = "markdown"
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func getWeekKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("понедельник"),
			tgbotapi.NewKeyboardButton("вторник"),
			tgbotapi.NewKeyboardButton("среда"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("четверг"),
			tgbotapi.NewKeyboardButton("пятница"),
			tgbotapi.NewKeyboardButton("суббота"),
			tgbotapi.NewKeyboardButton("воскресенье"),
		),
	)
}
func getEnv(key string, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}
