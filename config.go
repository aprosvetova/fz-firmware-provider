package main

type config struct {
	TelegramToken string `env:"TG_TOKEN,required"`
	UserID        int    `env:"USER_ID,required"`
}
