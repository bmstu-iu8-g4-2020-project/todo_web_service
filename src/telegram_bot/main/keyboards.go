package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"todo_web_service/src/telegram_bot/utils"
)

var weatherKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("На текущий момент времени", "current_weather"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Прогноз на ближайшие 5 дней", "weather_forecast"),
	),
)

var weatherChooseCurrentKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(utils.EmojiLocation+"По геопозиции", "curr_location"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(utils.EmojiMap+"По введённому месту", "curr_place_name"),
	),
)

var weatherChooseForecastKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(utils.EmojiLocation+"По геопозиции", "forecast_location"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(utils.EmojiMap+"По введённому месту", "forecast_place_name"),
	),
)
