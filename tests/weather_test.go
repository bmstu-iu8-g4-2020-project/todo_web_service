package tests

import (
	"fmt"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/utils"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/weather"
	"testing"
)

func TestPressure(t *testing.T) {
	results := []string{
		fmt.Sprintf("%.2f мм рт. ст", weather.TransferToMmHg(900)),
		fmt.Sprintf("%.2f мм рт. ст", weather.TransferToMmHg(1001)),
		fmt.Sprintf("%.2f мм рт. ст", weather.TransferToMmHg(989.98)),
	}

	expected := []string{
		"675.06 мм рт. ст",
		"750.81 мм рт. ст",
		"742.55 мм рт. ст",
	}

	for i := range results {
		if results[i] != expected[i] {
			t.Fatalf("Expected: %v. Result: %v.", expected[i], results[i])
		}
	}
}

func TestClouds(t *testing.T) {
	results := []string{
		weather.CloudsAllToEmoji(0),
		weather.CloudsAllToEmoji(28),
		weather.CloudsAllToEmoji(45),
		weather.CloudsAllToEmoji(100),
	}

	expected := []string{
		utils.EmojiWeatherClear,
		utils.EmojiWeatherFewClouds,
		utils.EmojiWeatherScatteredClouds,
		utils.EmojiWeatherOvercastClouds,
	}

	for i := range results {
		if results[i] != expected[i] {
			t.Fatalf("Expected: %v. Result: %v.", expected[i], results[i])
		}
	}
}

func TestWeatherToEmoji(t *testing.T) {
	results := []string{
		weather.WeatherIdToEmoji(800),
		weather.WeatherIdToEmoji(615),
		weather.WeatherIdToEmoji(311),
		weather.WeatherIdToEmoji(520),
		weather.WeatherIdToEmoji(511),
		weather.WeatherIdToEmoji(721),
		weather.WeatherIdToEmoji(804),
		weather.WeatherIdToEmoji(212),
	}

	expected := []string{
		utils.EmojiWeatherClear,
		utils.EmojiWeatherSnow,
		utils.EmojiWeatherDrizzleRain,
		utils.EmojiWeatherRain,
		utils.EmojiWeatherFreezingRain,
		utils.EmojiWeatherMist,
		utils.EmojiWeatherOvercastClouds,
		utils.EmojiWeatherThunderRain,
	}

	for i := range results {
		if results[i] != expected[i] {
			t.Fatalf("Expected: %v. Result: %v.", expected[i], results[i])
		}
	}
}

func TestDegToDirection(t *testing.T) {
	results := []string{
		weather.DegToDirection(0),
		weather.DegToDirection(35),
		weather.DegToDirection(100),
		weather.DegToDirection(120),
		weather.DegToDirection(175),
		weather.DegToDirection(210),
		weather.DegToDirection(284),
		weather.DegToDirection(330),
	}

	expected := []string{
		utils.EmojiWeatherNorth + "С",
		utils.EmojiWeatherNorthEast + "СВ",
		utils.EmojiWeatherEast + "В",
		utils.EmojiWeatherSouthEast + "ЮВ",
		utils.EmojiWeatherSouth + "Ю",
		utils.EmojiWeatherSouthWest + "ЮЗ",
		utils.EmojiWeatherWest + "З",
		utils.EmojiWeatherNorthWest + "СЗ",
	}

	for i := range results {
		if results[i] != expected[i] {
			t.Fatalf("Expected: %v. Result: %v.", expected[i], results[i])
		}
	}
}
