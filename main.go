package main

import (
	"flag"
	"log"

	tgClient "github.com/untibullet/dailyhelper/clients/telegram"
	"github.com/untibullet/dailyhelper/consumer"
	"github.com/untibullet/dailyhelper/events/telegram"
	"github.com/untibullet/dailyhelper/storage/files"
)

const (
	host        = "api.telegram.org"
	storagePath = "~/tg-bot-dailyhelper/users-data"
	batchSize   = 100
)

func main() {
	client := tgClient.NewClient(host, mustToken())

	fileStorage, err := files.NewStrorage(storagePath)
	if err != nil {
		log.Fatal("can`t init file storage", err)
	}

	processor := telegram.NewProcessor(client, fileStorage)

	consumer := consumer.New(processor, processor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
