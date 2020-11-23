package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"todo_web_service/src/telegram_bot/utils"
)

var weekdayScheduleKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Понедельник", "add_mon"),
		tgbotapi.NewInlineKeyboardButtonData("Четверг", "add_thu"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вторник", "add_tue"),
		tgbotapi.NewInlineKeyboardButtonData("Пятница", "add_fri"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Среда", "add_wed"),
		tgbotapi.NewInlineKeyboardButtonData("Суббота", "add_sat"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Воскресенье", "add_sun"),
	),
)

var scheduleDeleteKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(utils.EmojiTitle+
			"Удалить задачу", "delete_schedule_task"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(utils.EmojiWeekday+
			"Очистить расписание на день", "clear_weekday_schedule"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(utils.EmojiFire+
			"Полностью очистить расписание", "clear_schedule"),
	),
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
