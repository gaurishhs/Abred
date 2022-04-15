package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apidev234/abred/api"
	"github.com/apidev234/abred/database"
	"github.com/apidev234/abred/handlers"
	"github.com/apidev234/abred/registry"
	"github.com/apidev234/abred/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wshandler(w http.ResponseWriter, r *http.Request, s *discordgo.Session) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				uid := r.URL.Query().Get("uid")
				if uid == "" {
					conn.Close()
					return
				}

				data := utils.GetUserVoiceData(uid, s)
				if data == nil {
					err = conn.WriteMessage(1, []byte("none"))
					if err != nil {
						return
					}
					return
				}
				reqbodyBytes := new(bytes.Buffer)
				json.NewEncoder(reqbodyBytes).Encode(data)
				err = conn.WriteMessage(1, []byte(reqbodyBytes.Bytes()))
				if err != nil {
					return
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(CORSMiddleware())
	err := godotenv.Load()
	abred, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		fmt.Println("Error creating a new discord session", err)
		os.Exit(1)
	}
	abred.Identify.Intents = discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuilds | discordgo.IntentsGuildMessages
	abred.AddHandler(handlers.VoiceStateUpdate)
	abred.AddHandler(handlers.GuildCreate)
	abred.AddHandler(handlers.GuildDelete)
	abred.AddHandler(handlers.InteractionCreate)
	abred.AddHandler(handlers.MessageCreate)
	abred.AddHandler(handlers.Ready)
	registry.RegisterCommands(abred)
	err = abred.Open()
	if err != nil {
		fmt.Println("Error opening the connection: ", err)
		os.Exit(1)
	}
	database.SetupDB()
	r.GET("/manage", func(ctx *gin.Context) {
		wshandler(ctx.Writer, ctx.Request, abred)
	})
	// r.GET("/guilds", func(c *gin.Context) { api.GetGuilds(c, abred) })
	// r.GET("/channels", func(c *gin.Context) { api.GetChannels(c, abred) })
	r.GET("/cname", func(ctx *gin.Context) { api.GetChannelName(ctx, abred) })
	// r.POST("/generator", func(ctx *gin.Context) { api.NewGenerator(ctx, abred) })
	r.POST("/inc-limit", func(ctx *gin.Context) { api.IncreaseUserLimit(ctx, abred) })
	r.POST("/ptt", func(ctx *gin.Context) { api.TogglePushToTalk(ctx, abred) })
	r.POST("/dec-limit", func(ctx *gin.Context) { api.DecreaseUserLimit(ctx, abred) })
	// r.DELETE("/generator", func(ctx *gin.Context) { api.DeleteGenerator(ctx, abred) })
	// r.GET("/guild", func(ctx *gin.Context) { api.GetGuildData(ctx, abred) })
	r.Run()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	abred.Close()
}
