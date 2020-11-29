package weather

import (
	"fmt"
	"math"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	owm "github.com/briandowns/openweathermap"

	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/utils"
)

const (
	PressureTransferConst = 0.7500637554192
	MinWindDeg            = 0
	EastDeg               = 90
	SouthDeg              = 180
	WestDeg               = 270
	MaxWindDeg            = 360
	LayoutTime            = "15:40"
)

func TransferToMmHg(pressure float64) float64 {
	return pressure * PressureTransferConst
}

func WeatherIdToEmoji(weatherID int) string {
	switch weatherID {
	// Ясно.
	case 800:
		return utils.EmojiWeatherClear
	// Снег.
	case 600, 601, 602, 611, 612, 613, 615, 616, 620, 621, 622:
		return utils.EmojiWeatherSnow
	// Моросящий дождь.
	case 300, 301, 302, 310, 311, 312, 313, 314, 321:
		return utils.EmojiWeatherDrizzleRain
	// Дождь.
	case 500, 501, 502, 503, 504, 511, 520, 521, 522, 531:
		if weatherID == 511 {
			return utils.EmojiWeatherFreezingRain
		} else if weatherID < 502 {
			return utils.EmojiWeatherDrizzleRain
		} else {
			return utils.EmojiWeatherRain
		}
	// Облачно.
	case 801, 802, 803, 804:
		if weatherID == 803 || weatherID == 804 {
			return utils.EmojiWeatherOvercastClouds
		} else if weatherID == 801 {
			return utils.EmojiWeatherFewClouds
		} else if weatherID == 802 {
			return utils.EmojiWeatherScatteredClouds
		}
	// Гроза.
	case 200, 201, 202, 210, 211, 212, 221, 230, 231, 232:
		if 202 <= weatherID || weatherID >= 230 {
			return utils.EmojiWeatherThunderRain
		} else {
			return utils.EmojiWeatherThunder
		}
	// Туман, задымлённость.
	case 701, 711, 721, 731, 741, 751, 761, 762, 771, 781:
		return utils.EmojiWeatherMist
	}

	return ""
}

func DegToDirection(deg float64) string {
	halfDelta := 22.5 //

	if deg >= MinWindDeg && deg <= halfDelta || deg <= MaxWindDeg && deg >= MaxWindDeg-halfDelta {
		return utils.EmojiWeatherNorth + "С"
	} else if deg > halfDelta && deg <= EastDeg-halfDelta {
		return utils.EmojiWeatherNorthEast + "СВ"
	} else if deg > EastDeg-halfDelta && deg <= EastDeg+halfDelta {
		return utils.EmojiWeatherEast + "В"
	} else if deg > EastDeg+halfDelta && deg <= SouthDeg-halfDelta {
		return utils.EmojiWeatherSouthEast + "ЮВ"
	} else if deg > SouthDeg-halfDelta && deg <= SouthDeg+halfDelta {
		return utils.EmojiWeatherSouth + "Ю"
	} else if deg > SouthDeg+halfDelta && deg <= WestDeg-halfDelta {
		return utils.EmojiWeatherSouthWest + "ЮЗ"
	} else if deg > WestDeg-halfDelta && deg <= WestDeg+halfDelta {
		return utils.EmojiWeatherWest + "З"
	} else if deg > WestDeg+halfDelta && deg <= MaxWindDeg-halfDelta {
		return utils.EmojiWeatherNorthWest + "СЗ"
	}

	return "Штиль"
}

func CloudsAllToEmoji(cloudsAll int) string {
	if cloudsAll >= 0 && cloudsAll <= 10 {
		return utils.EmojiWeatherClear
	} else if cloudsAll <= 30 {
		return utils.EmojiWeatherFewClouds
	} else if cloudsAll <= 60 {
		return utils.EmojiWeatherScatteredClouds
	}
	return utils.EmojiWeatherOvercastClouds
}

func WeatherOutput(w *owm.CurrentWeatherData) string {
	var output strings.Builder
	var weather owm.Weather
	if len(w.Weather) > 1 {
		weather = w.Weather[len(w.Weather)-1]
	} else {
		weather = w.Weather[0]
	}

	_, _ = fmt.Fprintf(&output, "%sПогода: %s.\n", utils.EmojiLocation, w.Name)
	_, _ = fmt.Fprintf(&output, "Сейчас на улице:\n%s%s\n",
		WeatherIdToEmoji(weather.ID), strings.Title(weather.Description))
	_, _ = fmt.Fprintf(&output, "Температура воздуха: %v°C\n", math.Round(w.Main.Temp))
	_, _ = fmt.Fprintf(&output, "Ощущается как: %v°C\n", math.Round(w.Main.FeelsLike))
	_, _ = fmt.Fprintf(&output, "Влажность воздуха: %v%%\n", w.Main.Humidity)
	_, _ = fmt.Fprintf(&output, "Атмосферное давление: \n%.2f мм рт. ст.\n", TransferToMmHg(w.Main.Pressure))
	_, _ = fmt.Fprintf(&output, "Ветер: %s %v м/c\n", DegToDirection(w.Wind.Deg), math.Round(w.Wind.Speed))
	_, _ = fmt.Fprintf(&output, "Облачность:%s %v%% \n", CloudsAllToEmoji(w.Clouds.All), w.Clouds.All)

	return output.String()
}

func ForecastOutput(w *owm.Forecast5WeatherList) string {
	var output strings.Builder

	date := w.DtTxt.Format("02.01.2006")

	weather := w.Weather[0]
	_, _ = fmt.Fprintf(&output, "%sПогода на %s.\n", utils.EmojiWeekday, date)
	_, _ = fmt.Fprintf(&output, "%s%s\n", WeatherIdToEmoji(weather.ID), strings.Title(weather.Description))
	_, _ = fmt.Fprintf(&output, "Температура воздуха: %v°C\n", math.Round(w.Main.Temp))
	_, _ = fmt.Fprintf(&output, "Влажность воздуха: %v%%\n", w.Main.Humidity)
	_, _ = fmt.Fprintf(&output, "Атмосферное давление: \n%.2f мм рт. ст.\n", TransferToMmHg(w.Main.Pressure))
	_, _ = fmt.Fprintf(&output, "Ветер: %s %v м/c\n", DegToDirection(w.Wind.Deg), math.Round(w.Wind.Speed))
	_, _ = fmt.Fprintf(&output, "Облачность:%s %v%% \n", CloudsAllToEmoji(w.Clouds.All), w.Clouds.All)

	return output.String()
}

func SendForecast(w *owm.ForecastWeatherData, update *tgbotapi.Update, bot **tgbotapi.BotAPI) {
	if forecast, ok := w.ForecastWeatherJson.(*owm.Forecast5WeatherData); ok {
		_, _ = (*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Прогноз погоды: %s (%s)",
			forecast.City.Name, forecast.City.Country)))

		for i := 0; i < len(forecast.List); i++ {
			if forecast.List[i].DtTxt.Format(LayoutTime) == "15:00" {
				output := ForecastOutput(&forecast.List[i])
				_, _ = (*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, output))
				i += 7
			}
		}
	}
}

func SendWrongPlaceName(update *tgbotapi.Update, bot **tgbotapi.BotAPI) {
	_, _ = (*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf("%sВы ввели непонятное название места, данные отыскать не удалось.",
			utils.EmojiWarning)))
}

func SendWrongLocation(update *tgbotapi.Update, bot **tgbotapi.BotAPI) {
	_, _ = (*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf("%s Что-то не так с вашей геопозицией, данные отыскать не удалось.",
			utils.EmojiWarning)))
}

func SendInternalWeatherAPIError(update *tgbotapi.Update, bot **tgbotapi.BotAPI) {
	_, _ = (*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf("%sВнутренняя ошибка openweathermap.org. Попробуйте позднее.",
			utils.EmojiWarning)))
}
