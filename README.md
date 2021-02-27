# Разработка веб-сервиса по планировке пользовательских задач

<img alt="Go" src="https://img.shields.io/badge/go-%2300ADD8.svg?&style=for-the-badge&logo=go&logoColor=white"/><img alt="Postgres" src ="https://img.shields.io/badge/postgres-%23316192.svg?&style=for-the-badge&logo=postgresql&logoColor=white"/><img alt="Docker" src="https://img.shields.io/badge/docker%20-%230db7ed.svg?&style=for-the-badge&logo=docker&logoColor=white"/><img alt="TravisCI" src="https://img.shields.io/badge/travisci%20-%232B2F33.svg?&style=for-the-badge&logo=travis&logoColor=white"/>

[![Build Status](https://travis-ci.com/bmstu-iu8-g4-2020-project/todo_web_service.svg?branch=master)](https://travis-ci.com/bmstu-iu8-g4-2020-project/todo_web_service)
[![Go Report Card](https://goreportcard.com/badge/github.com/bmstu-iu8-g4-2020-project/todo_web_service)](https://goreportcard.com/report/github.com/bmstu-iu8-g4-2020-project/todo_web_service)

<img alt="Telegram" src="https://img.shields.io/badge/aaaaaaaalesha-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white" />


### Задание:

***

#### Разработать сервис, реализующий создание пользовательских задач и списков дел. Создать подпрограмму для мессенджера Telegram, которая посредством HTTP запросов будет взаимодействовать с сервисом. С помощью этой подпрограммы, пользователь должен иметь возможность составлять регулярное расписание, создавать/редактировать списки текущих дел, получать необходимые напоминания об установленных планах и задачах, получать полезную информацию от сторонних сервисов для правильного планирования распорядка дня.

***

## Архитектура проекта

![Project Architecture](assets/images/project_architecture.png)

## Конечный автомат для состояний пользователей

Пакет [__telegram-bot-api__](https://github.com/go-telegram-bot-api/telegram-bot-api) для работы с ботом «из коробки» не предоставлял
возможность выстраивать цепочки сообщений с пользователями, из-за чего
пришлось хранить состояния пользователей в виде уникальных кодов на
каждое состояние, а переходы между ними осуществлять в виде абстрактной
модели [конечного автомата](https://ru.wikipedia.org/wiki/%D0%9A%D0%BE%D0%BD%D0%B5%D1%87%D0%BD%D1%8B%D0%B9_%D0%B0%D0%B2%D1%82%D0%BE%D0%BC%D0%B0%D1%82). Все состояния сохраняются в базе
данных для каждого пользователя, таким образом падение сервера или
перезапуск бота не приводит к потере состояния для пользователя.

![FMS](assets/images/FMS_for_user_states.png)

## Основные команды

[![Bot and commands](assets/images/bot&commands.png)](https://t.me/todownik_bot "t.me/todownik_bot")

```
Copyright 2020 aaaaaaaalesha 
```
