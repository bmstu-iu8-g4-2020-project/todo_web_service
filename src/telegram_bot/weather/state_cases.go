package weather

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	owm "github.com/briandowns/openweathermap"

	"todo_web_service/src/telegram_bot/user"
	"todo_web_service/src/telegram_bot/utils"
)

func FillWeatherFuncs(stateFuncDict *map[int]user.StateFunc) {
	(*stateFuncDict)[user.WEATHER_SEND_LOCATION] = SendLocation
}

func SendLocation(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	msg := update.Message
	if update.Message.Location == nil {
		_, _ = (*bot).Send(tgbotapi.NewMessage(msg.Chat.ID,
			fmt.Sprintf("%sНеобходимо отправить геопозицию для получения данных о погоде. Воспользуйтесь %s для отправки геопозиции.",
				utils.EmojiWarning, utils.EmojiPaperclip)))
		return
	}

	weatherApiKey := os.Getenv("WEATHER_API_KEY")

	w, err := owm.NewCurrent("C", "ru", weatherApiKey)
	if err != nil {
		log.Fatalln(err)
	}

	err = w.CurrentByCoordinates(&owm.Coordinates{
		Longitude: msg.Location.Longitude,
		Latitude:  msg.Location.Latitude,
	})
	if err != nil {
		log.Fatal(err)
	}

	var output strings.Builder

	fmt.Fprintf(&output, "%s Погода в городе %s.\n\n", utils.EmojiLocation, w.Name)
	fmt.Fprintf(&output, "На улице: %s\n", w.Weather[0].Description)
	fmt.Fprintf(&output, "Температура воздуха: %v °C\n", w.Main.Temp)
	fmt.Fprintf(&output, "Ощущается как: %v °C\n", w.Main.FeelsLike)
	fmt.Fprintf(&output, "Скорость ветра: %v м/c\n", w.Wind.Speed)
	fmt.Fprintf(&output, "Облачность: %v%% \n", w.Clouds.All)

	_, _ = (*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, output.String()))

	_ = user.ResetState(msg.From.ID, msg.From.UserName, userStates)
}
