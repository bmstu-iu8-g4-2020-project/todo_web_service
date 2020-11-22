package weather

import (
	"fmt"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	owm "github.com/briandowns/openweathermap"

	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/user"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/utils"
)

func FillWeatherFuncs(stateFuncDict *map[int]user.StateFunc) {
	(*stateFuncDict)[user.WEATHER_CURRENT_SEND_LOCATION] = CurrentSendLocation
	(*stateFuncDict)[user.WEATHER_CURRENT_SEND_NAME] = CurrentByName
	(*stateFuncDict)[user.WEATHER_FORECAST_CHOOSE_INPUT] = ForecastChooseInput
	(*stateFuncDict)[user.WEATHER_FORECAST_SEND_LOCATION] = ForecastByLocation
	(*stateFuncDict)[user.WEATHER_FORECAST_SEND_NAME] = ForecastByName
}

func CurrentSendLocation(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
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
		SendInternalWeatherAPIError(update, bot)
		_ = user.ResetState(update.Message.From.ID, msg.From.UserName, userStates)

		return
	}

	err = w.CurrentByCoordinates(&owm.Coordinates{
		Longitude: msg.Location.Longitude,
		Latitude:  msg.Location.Latitude,
	})
	if err != nil {
		SendWrongLocation(update, bot)
		_ = user.ResetState(msg.From.ID, msg.From.UserName, userStates)
		return
	}

	_, _ = (*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, WeatherOutput(w)))

	_ = user.ResetState(msg.From.ID, msg.From.UserName, userStates)
}

func CurrentByName(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	msg := update.Message
	if update.Message.Text == "" {
		_, _ = (*bot).Send(tgbotapi.NewMessage(msg.Chat.ID,
			fmt.Sprintf("%sВы прислали не название места. Попробуйте снова.", utils.EmojiWarning)))
		return
	}

	weatherApiKey := os.Getenv("WEATHER_API_KEY")

	w, err := owm.NewCurrent("C", "ru", weatherApiKey)
	if err != nil {
		SendInternalWeatherAPIError(update, bot)
		_ = user.ResetState(update.Message.From.ID, msg.From.UserName, userStates)
		return
	}

	err = w.CurrentByName(msg.Text)
	if err != nil {
		SendWrongPlaceName(update, bot)
		return
	}

	_, _ = (*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, WeatherOutput(w)))

	_ = user.ResetState(msg.From.ID, msg.From.UserName, userStates)
}

func ForecastChooseInput(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	msg := update.Message

	if msg.Text == "" {
		_, _ = (*bot).Send(tgbotapi.NewMessage(msg.Chat.ID,
			fmt.Sprintf("%sВы прислали не текстовое сообщение. Попробуйте снова.", utils.EmojiWarning)))
		return
	}

	switch msg.Text {
	case "1":
		_, _ = (*bot).Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf(
			"Пришлите мне свою геопозицию. \n(нажмите на %s и выберите \"Геопозиция\")", utils.EmojiPaperclip)))
		_ = user.SetState(msg.From.ID, msg.From.UserName, userStates,
			user.State{Code: user.WEATHER_FORECAST_SEND_LOCATION, Request: "{}"})
	case "2":
		_, _ = (*bot).Send(tgbotapi.NewMessage(msg.Chat.ID, utils.EmojiLocation+
			"Введите место, где вы бы хотели узнать данные о погоде."))
		_ = user.SetState(msg.From.ID, msg.From.UserName, userStates,
			user.State{Code: user.WEATHER_FORECAST_SEND_NAME, Request: "{}"})
	default:
		_, _ = (*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("%sВыберите 1 или 2 чтобы выбрать способ ввода. Попробуйте снова.", utils.EmojiWarning)))
	}
}

func ForecastByLocation(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	msg := update.Message
	if update.Message.Location == nil {
		_, _ = (*bot).Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf(
			"%sНеобходимо отправить геопозицию для получения данных о погоде. Воспользуйтесь %s для отправки геопозиции.",
			utils.EmojiWarning, utils.EmojiPaperclip)))
		return
	}

	weatherApiKey := os.Getenv("WEATHER_API_KEY")

	w, err := owm.NewForecast("5", "C", "ru", weatherApiKey)
	if err != nil {
		SendInternalWeatherAPIError(update, bot)
		_ = user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)
		return
	}

	err = w.DailyByCoordinates(&owm.Coordinates{
		Longitude: msg.Location.Longitude,
		Latitude:  msg.Location.Latitude,
	}, 0)
	if err != nil {
		SendWrongLocation(update, bot)
		_ = user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)
		return
	}

	SendForecast(w, update, bot)

	_ = user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)
}

func ForecastByName(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	msg := update.Message
	if update.Message.Text == "" {
		_, _ = (*bot).Send(tgbotapi.NewMessage(msg.Chat.ID,
			fmt.Sprintf("%sВы прислали не название места. Попробуйте снова.", utils.EmojiWarning)))
		return
	}

	weatherApiKey := os.Getenv("WEATHER_API_KEY")
	w, err := owm.NewForecast("5", "C", "ru", weatherApiKey)
	if err != nil {
		SendInternalWeatherAPIError(update, bot)
		_ = user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)
		return
	}

	err = w.DailyByName(msg.Text, 0)

	if err != nil {
		SendWrongPlaceName(update, bot)
		return
	}

	SendForecast(w, update, bot)

	_ = user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)
}
