package utils

import (
	"fmt"
	"time"

	"github.com/apidev234/abred/database"
	"github.com/bwmarrin/discordgo"
)

func GetGuildOwner(s *discordgo.Session, guildID string) string {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return ""
	}
	return guild.OwnerID
}

func HasPerms(s *discordgo.Session, m *discordgo.Message, guildID, userID string, permissions ...int) bool {
	if GetGuildOwner(s, guildID) == userID {
		return true
	}

	var (
		err   error
		guild *discordgo.Guild
		perms int64
	)

	switch m != nil {
	case true:

		perms, err = s.State.MessagePermissions(m)
		if err != nil {
			return false
		}

	case false:

		guild, err = s.State.Guild(guildID)
		if err != nil {
			return false
		}

		if len(guild.Channels) == 0 {
			return false
		}

		perms, err = s.State.UserChannelPermissions(userID, guild.Channels[0].ID)
		if err != nil {
			return false
		}
	}

	for _, perm := range permissions {
		if perms&(int64(perm)) != int64(perm) {
			continue
		}
		return true
	}

	return false
}

func IsGenerator(data *database.Guild, cid string) bool {
	var result bool = false
	for _, gen := range data.Generators {
		if gen == cid {
			result = true
			break
		}
	}
	return result
}

func CheckIfChannelEmpty(channelid string, guildid string, s *discordgo.Session) bool {
	guild, err := s.State.Guild(guildid)
	if err != nil {
		return false
	}
	var list []string
	for _, vS := range guild.VoiceStates {
		if vS.ChannelID == channelid {
			list = append(list, vS.UserID)
		}
	}
	if len(list) == 0 {
		return true
	}
	return false
}

type UserVoiceData struct {
	GuildId     string `json:"guildId"`
	ChannelId   string `json:"channelId"`
	GuildName   string `json:"guildName"`
	ChannelName string `json:"channelName"`
	GuildIcon   string `json:"guildIcon"`
	PushToTalk  bool   `json:"pushToTalk"`
	Limit       int    `json:"userLimit"`
	Locked      bool   `json:"locked"`
}

func PushToTalkEnabled(channelId string, s *discordgo.Session) bool {
	channel, err := s.Channel(channelId)
	if err != nil {
		return false
	}
	for _, p := range channel.PermissionOverwrites {
		if p.ID == channel.GuildID {
			if p.Deny&discordgo.PermissionVoiceUseVAD == 0 {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func Locked(channelId string, s *discordgo.Session) bool {
	channel, err := s.Channel(channelId)
	if err != nil {
		return false
	}
	for _, p := range channel.PermissionOverwrites {
		if p.ID == channel.GuildID {
			if p.Deny&discordgo.PermissionVoiceConnect == 0 {
				return false
			} else {
				return true
			}
		}
	}
	return false
}

func GetUserVoiceData(userid string, s *discordgo.Session) *UserVoiceData {
	for _, guild := range s.State.Guilds {
		for _, vS := range guild.VoiceStates {
			if vS.UserID == userid {
				ischan := database.Database.IsAChannel(vS.GuildID, vS.ChannelID)
				if !ischan {
					return nil
				}
				guild, err := s.Guild(vS.GuildID)
				if err != nil {
					return nil
				}
				channel, err := s.Channel(vS.ChannelID)
				if err != nil {
					return nil
				}
				ptt := PushToTalkEnabled(channel.ID, s)
				lcked := Locked(channel.ID, s)
				buaia := UserVoiceData{
					GuildId:     vS.GuildID,
					ChannelId:   vS.ChannelID,
					GuildName:   guild.Name,
					GuildIcon:   guild.IconURL(),
					ChannelName: channel.Name,
					PushToTalk:  ptt,
					Limit:       channel.UserLimit,
					Locked:      lcked,
				}
				return &buaia
			}
		}
	}
	return nil
}

func GetUserVoiceChannel(s *discordgo.Session, userId string, guildId string) *string {
	guild, err := s.State.Guild(guildId)
	if err != nil {
		return nil
	}
	for _, vS := range guild.VoiceStates {
		if vS.UserID == userId {
			return &vS.ChannelID
		}
	}
	return nil
}

func ComposeErrorMessage(s *discordgo.Session, channelId string, message string, usrmsg string) {
	msg, err := s.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "<:Abred_Fail:940541495491051560> __**ERROR**__ <:Abred_Fail:940541495491051560>",
				Description: message,
				Color:       16711680,
			},
		},
	})
	if err != nil {
		return
	}
	time.AfterFunc(30*time.Second, func() {
		err := s.ChannelMessageDelete(msg.ChannelID, msg.ID)
		if err != nil {
			return
		}
		err = s.ChannelMessageDelete(channelId, usrmsg)
		if err != nil {
			return
		}
	})
}

func ComposeSuccessMessage(s *discordgo.Session, channelId string, message string, usrmsg string) {
	msg, err := s.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "<:Abred_Success:940541058192908350> __**SUCCESS**__ <:Abred_Success:940541058192908350>",
				Description: message,
				Color:       3342080,
			},
		},
	})
	if err != nil {
		return
	}
	time.AfterFunc(30*time.Second, func() {
		err := s.ChannelMessageDelete(msg.ChannelID, msg.ID)
		if err != nil {
			return
		}
		err = s.ChannelMessageDelete(channelId, usrmsg)
		if err != nil {
			return
		}
	})
}

func ComposeIntErrorMessage(s *discordgo.Session, interaction *discordgo.Interaction, msg string) {
	s.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 64,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "<:Abred_Fail:940541495491051560> __**ERROR**__ <:Abred_Fail:940541495491051560>",
					Description: msg,
					Color:       16711680,
				},
			},
		},
	})
}

func ComposeIntSuccessMessage(s *discordgo.Session, i *discordgo.Interaction, msg string) {
	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 64,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "<:Abred_Success:940541058192908350> __**SUCCESS**__ <:Abred_Success:940541058192908350>",
					Description: msg,
					Color:       3342080,
				},
			},
		},
	})
}

func GetChannelName(cid string, s *discordgo.Session) string {
	channel, err := s.Channel(cid)
	if err != nil {
		return "Invalid Channel"
	}
	return channel.Name
}

func WrapGenerators(arr []string) []string {
	var rearr []string
	for _, id := range arr {
		rearr = append(rearr, fmt.Sprintf("<#%s>", id))
	}
	return rearr
}

/*
 * Abred
 * Copyright (C) 2022 ApiDev234
 * This software is licensed under Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International
 * For more information, see README.md and LICENSE
 */
