package main

import tgbotapi "github.com/Syfaro/telegram-bot-api"

var weatherChooseCurrentKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("По геопозиции", "curr_location"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("По введённому месту", "curr_place_name"),
	),
)
