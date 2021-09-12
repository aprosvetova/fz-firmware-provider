package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"strconv"
	"time"
)

var cfg config
var firmware []byte
var checksum uint32

func main() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Config", err)
	}

	b, err := tb.NewBot(tb.Settings{
		Token: cfg.TelegramToken,
		Poller: tb.NewMiddlewarePoller(&tb.LongPoller{Timeout: 10 * time.Second},
			func(upd *tb.Update) bool {
				if upd.Message == nil {
					return true
				}
				return upd.Message.Sender.ID == cfg.UserID
			}),
		ParseMode: tb.ModeHTML,
	})

	if err != nil {
		log.Fatalln("Telegram", err)
		return
	}

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "<b>( ´ ▽ ` )ﾉ</b> Welcome!")
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		err := loadFirmware(m.Text)
		if err != nil {
			b.Reply(m, "<b>(︶︹︺)</b> Oops! "+err.Error())
		} else {
			b.Reply(m, "<b>(*^‿^*)</b> Loaded successfully!")
		}
	})

	go b.Start()

	log.Println("Server started")

	r := gin.New()

	r.GET("/checksum", serveChecksum)
	r.GET("/firmware.bin", serveBinary)

	log.Fatal(r.Run(":8080"))
}

func loadFirmware(url string) error {
	var err error
	_, firmware, err = fasthttp.Get(nil, url)
	if err != nil {
		return err
	}
	checksum, err = stm32crc(firmware)
	return err
}

func serveChecksum(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("%08x", checksum))
}

func serveBinary(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename=firmware.bin")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.Itoa(len(firmware)))
	c.Writer.Write(firmware)
}
