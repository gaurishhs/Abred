package registry

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func RegisterCommands(s *discordgo.Session) {
	var y int64 = 8
	var x float64 = float64(y)
	var minul int64 = 1
	var minul2 float64 = float64(minul)
	var commands = []*discordgo.ApplicationCommand{
		{
			Name:        "setup",
			Description: "Setup the bot",
		},
		{
			Name:        "help",
			Description: "Get help regarding the usage of bot",
		},

		{
			Name:        "config",
			Description: "Get Guild's Configuration",
		},
		{
			Name:        "newgen",
			Description: "Register a new generator",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "channel",
					Description: "The new generator channel",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionChannel,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildVoice,
					},
				},
			},
		},
		{
			Name:        "delgen",
			Description: "Remove a already registered generator",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "channel",
					Description: "The generator channel",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionChannel,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildVoice,
					},
				},
			},
		},
		{
			Name:        "default-state",
			Description: "Configure the default state for the guild",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name: "state",
					Type: discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "locked",
							Value: "locked",
						},
						{
							Name:  "unlocked",
							Value: "unlocked",
						},
					},
					Required:    true,
					Description: "The new state of the guild",
				},
			},
		},
		{
			Name:        "default-bitrate",
			Description: "Configure the default bitrate for the guild",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "bitrate",
					Description: "The bitrate to be used",
					Type:        discordgo.ApplicationCommandOptionInteger,
					MaxValue:    96,
					MinValue:    &x,
					Required:    true,
				},
			},
		},
		{
			Name:        "default-user-limit",
			Description: "Configure the default user limit for the guild",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "limit",
					Description: "The limit to be used",
					Type:        discordgo.ApplicationCommandOptionInteger,
					MaxValue:    99,
					MinValue:    &minul2,
					Required:    true,
				},
			},
		},
		{
			Name:        "default-category",
			Description: "Configure the default category for the guild",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionChannel,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildCategory,
					},
					Name:        "channel",
					Description: "The category",
					Required:    true,
				},
			},
		},
	}
	err := godotenv.Load()
	if err != nil {
		fmt.Print("Could not lo ad .env in registry")
	}
	_, err = s.ApplicationCommandBulkOverwrite(os.Getenv("CLIENT_ID"), "", commands)
	if err != nil {
		fmt.Println("Bulk Overwrite Error:", err)
	}
}

/*
 * Abred
 * Copyright (C) 2022 ApiDev234
 * This software is licensed under Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International
 * For more information, see README.md and LICENSE
 */
