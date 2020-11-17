package weather

import "todo_web_service/src/telegram_bot/utils"

const (
	PressureTransferConst = 0.7500637554192
	MinWindDeg            = 0
	EastDeg               = 90
	SouthDeg              = 180
	WestDeg               = 270
	MaxWindDeg            = 360
)

func TransferToMmHg(pressure float64) float64 {
	return pressure * PressureTransferConst
}

func WeatherIdToEmoji(weatherID int) string {
	switch weatherID {
	// Ясно.
	case 800:
		return utils.EmojiWeatherClear
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
	// Снег.
	case 600, 601, 602, 611, 612, 613, 615, 616, 620, 621, 622:
		return utils.EmojiWeatherSnow
	// Туман, задымлённость.
	case 701, 711, 721, 731, 741, 751, 761, 762, 771, 781:
		return utils.EmojiWeatherMist
	}

	return ""
}

// Шпора: https://dpva.ru/Guide/GuideUnitsAlphabets/GuideUnitsAlphabets/WindRoseRuEng/
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