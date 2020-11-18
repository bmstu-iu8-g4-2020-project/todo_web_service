package weather

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	owm "github.com/briandowns/openweathermap"
	"log"
	"os"

	"todo_web_service/src/telegram_bot/user"
	"todo_web_service/src/telegram_bot/utils"
)

func FillWeatherFuncs(stateFuncDict *map[int]user.StateFunc) {
	(*stateFuncDict)[user.WEATHER_CURRENT_SEND_LOCATION] = SendLocation
	(*stateFuncDict)[user.WEATHER_CURRENT_SEND_NAME] = SendName
	(*stateFuncDict)[user.WEATHER_FORECAST] = Forecast
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
		_, _ = (*bot).Send(tgbotapi.NewMessage(msg.Chat.ID,
			fmt.Sprintf("%s Что-то не так с вашей геопозицией, данные отыскать не удалось,"+
				" воспользуйтесь /place_weather и введите название места, где вы находитесь.",
				utils.EmojiWarning)))
		_ = user.ResetState(msg.From.ID, msg.From.UserName, userStates)
		return
	}

	_, _ = (*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
		WeatherOutput(w)))

	_ = user.ResetState(msg.From.ID, msg.From.UserName, userStates)
}

func SendName(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	msg := update.Message
	if update.Message.Text == "" {
		_, _ = (*bot).Send(tgbotapi.NewMessage(msg.Chat.ID,
			fmt.Sprintf("%sВы прислали не название места. Попробуйте снова.", utils.EmojiWarning)))
		return
	}

	weatherApiKey := os.Getenv("WEATHER_API_KEY")

	w, err := owm.NewCurrent("C", "ru", weatherApiKey)
	if err != nil {
		log.Fatalln(err)
	}

	err = w.CurrentByName(msg.Text)
	if err != nil {
		_, _ = (*bot).Send(tgbotapi.NewMessage(msg.Chat.ID, utils.EmojiWarning+
			"Что-то не так с названием места, данные отыскать не удалось.\n"+
			"Воспользуйтесь \n/geopos_weather и пришлите геопозицию необходимого места.",
		))
		_ = user.ResetState(msg.From.ID, msg.From.UserName, userStates)
		return
	}

	_, _ = (*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
		WeatherOutput(w)))

	_ = user.ResetState(msg.From.ID, msg.From.UserName, userStates)
}

func Forecast(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	msg := update.Message
	if update.Message.Location == nil {
		_, _ = (*bot).Send(tgbotapi.NewMessage(msg.Chat.ID,
			fmt.Sprintf("%sНеобходимо отправить геопозицию для получения данных о погоде. Воспользуйтесь %s для отправки геопозиции.",
				utils.EmojiWarning, utils.EmojiPaperclip)))
		return
	}

	weatherApiKey := os.Getenv("WEATHER_API_KEY")

	w, err := owm.NewForecast("5", "C", "ru", weatherApiKey)
	if err != nil {
		log.Fatal(err)
	}

	err = w.DailyByCoordinates(&owm.Coordinates{
		Longitude: msg.Location.Longitude,
		Latitude:  msg.Location.Latitude,
	}, 0)
	if err != nil {
		_, _ = (*bot).Send(tgbotapi.NewMessage(msg.Chat.ID,
			fmt.Sprintf("%s Что-то не так с вашей геопозицией, данные отыскать не удалось,",
				utils.EmojiWarning)))
		_ = user.ResetState(msg.From.ID, msg.From.UserName, userStates)
		return
	}

	if forecast, ok := w.ForecastWeatherJson.(*owm.Forecast5WeatherData); ok {
		_, _ = (*bot).Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Прогноз погоды: %s (%s)",
			forecast.City.Name, forecast.City.Country)))

		for i := 0; i < len(forecast.List); i++ {
			if forecast.List[i].DtTxt.Format(LayoutTime) == "15:00" {
				output := ForecastOutput(&forecast.List[i])
				_, _ = (*bot).Send(tgbotapi.NewMessage(msg.Chat.ID, output))
				i += 7
			}
		}

	}

	_ = user.ResetState(msg.From.ID, msg.From.UserName, userStates)
}
