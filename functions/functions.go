package functions

import (
	"errors"

	"github.com/apidev234/abred/database"
	"github.com/bwmarrin/discordgo"
)

func LockChannel(userid string, channelId string, guildId string, s *discordgo.Session) error {
	channeldata, err := database.Database.GetChannelData(guildId, channelId, s)
	if err != nil {
		return errors.New("No such channel")
	}

	if channeldata.OwnerID != userid {
		return errors.New("Unauthorized, You aren't the owner of this channel")
	}

	_, err = s.ChannelEditComplex(channeldata.Id, &discordgo.ChannelEdit{
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: 0,
				Deny:  discordgo.PermissionVoiceConnect,
				ID:    guildId,
			},
		},
	})
	if err != nil {
		return errors.New("Couldn't Change The Channel's State")
	}
	return nil
}

func UnlockChannel(userid string, channelId string, guildId string, s *discordgo.Session) error {
	channelData, err := database.Database.GetChannelData(guildId, channelId, s)
	if err != nil {
		return errors.New("No such channel")
	}

	if channelData.OwnerID != userid {
		return errors.New("Unauthorized, You aren't the owner of this channel")
	}
	_, err = s.ChannelEditComplex(channelData.Id, &discordgo.ChannelEdit{
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionVoiceConnect | discordgo.PermissionViewChannel,
				Deny:  0,
				ID:    guildId,
			},
		},
	})
	if err != nil {
		return errors.New("Couldn't Change The Channel's State")
	}
	return nil
}

func ShowChannel(userid string, channelId string, guildId string, s *discordgo.Session) error {
	channeldata, err := database.Database.GetChannelData(guildId, channelId, s)
	if err != nil {
		return errors.New("No such channel")
	}
	if channeldata.OwnerID != userid {
		return errors.New("Unauthorized, You aren't the owner of this channel")
	}

	_, err = s.ChannelEditComplex(channeldata.Id, &discordgo.ChannelEdit{
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel,
				Deny:  0,
				ID:    guildId,
			},
		},
	})
	if err != nil {
		return errors.New("Couldn't Change The Channel's State")
	}
	return nil
}

func HideChannel(userid string, channelId string, guildId string, s *discordgo.Session) error {
	channeldata, err := database.Database.GetChannelData(guildId, channelId, s)
	if err != nil {
		return errors.New("No such channel")
	}
	if channeldata.OwnerID != userid {
		return errors.New("Unauthorized, You aren't the owner of this channel")
	}

	_, err = s.ChannelEditComplex(channeldata.Id, &discordgo.ChannelEdit{
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: 0,
				Deny:  discordgo.PermissionViewChannel,
				ID:    guildId,
			},
		},
	})
	if err != nil {
		return errors.New("Couldn't Change The Channel's State")
	}

	return nil
}

func IncreaseLimit(userid string, channelId string, guildId string, s *discordgo.Session) error {
	channeldata, err := database.Database.GetChannelData(guildId, channelId, s)
	if err != nil {
		return errors.New("No such channel")
	}
	if channeldata.OwnerID != userid {
		return errors.New("Unauthorized, You aren't the owner of this channel")
	}

	channel, err := s.Channel(channelId)
	if err != nil {
		return errors.New("No such channel")
	}
	_, err = s.ChannelEditComplex(channeldata.Id, &discordgo.ChannelEdit{
		UserLimit: channel.UserLimit + 1,
	})
	if err != nil {
		return errors.New("Couldn't Change The Channel's State")
	}

	return nil
}

func DecreaseLimit(userid string, channelId string, guildId string, s *discordgo.Session) error {
	channeldata, err := database.Database.GetChannelData(guildId, channelId, s)
	if err != nil {
		return errors.New("No such channel")
	}
	if channeldata.OwnerID != userid {
		return errors.New("Unauthorized, You aren't the owner of this channel")
	}

	channel, err := s.Channel(channelId)
	if err != nil {
		return errors.New("No such channel")
	}
	_, err = s.ChannelEditComplex(channeldata.Id, &discordgo.ChannelEdit{
		UserLimit: channel.UserLimit - 1,
	})
	if err != nil {
		return errors.New("Couldn't Change The Channel's State")
	}

	return nil
}

func TogglePushToTalk(userid string, channelId string, guildId string, s *discordgo.Session) error {
	channeldata, err := database.Database.GetChannelData(guildId, channelId, s)
	if err != nil {
		return errors.New("No such channel")
	}
	if channeldata.OwnerID != userid {
		return errors.New("Unauthorized, You aren't the owner of this channel")
	}

	channel, err := s.Channel(channelId)
	if err != nil {
		return errors.New("No such channel")
	}
	for _, p := range channel.PermissionOverwrites {
		if p.ID == channel.GuildID {
			if p.Deny&discordgo.PermissionVoiceUseVAD == 0 {
				_, err := s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{
					PermissionOverwrites: []*discordgo.PermissionOverwrite{
						{
							Allow: 0,
							Deny:  discordgo.PermissionVoiceUseVAD,
							ID:    guildId,
							Type:  discordgo.PermissionOverwriteTypeRole,
						},
					},
				})
				if err != nil {
					return errors.New("Could not enable push to talk")
				}
			} else {
				_, err := s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{
					PermissionOverwrites: []*discordgo.PermissionOverwrite{
						{
							Allow: discordgo.PermissionVoiceUseVAD,
							Deny:  0,
							ID:    guildId,
							Type:  discordgo.PermissionOverwriteTypeRole,
						},
					},
				})
				if err != nil {
					return errors.New("Could not disable push to talk")
				}
			}
		}
	}
	return nil
}

func SetBitrate(userid string, channelId string, guildId string, s *discordgo.Session, newBit int) error {
	channeldata, err := database.Database.GetChannelData(guildId, channelId, s)
	if err != nil {
		return errors.New("No such channel")
	}
	if channeldata.OwnerID != userid {
		return errors.New("Unauthorized, You aren't the owner of this channel")
	}

	_, err = s.ChannelEditComplex(channeldata.Id, &discordgo.ChannelEdit{
		Bitrate: newBit * 1000,
	})
	if err != nil {
		return errors.New("Couldn't Change The Channel's State")
	}

	return nil
}

/*
 * Abred
 * Copyright (C) 2022 ApiDev234
 * This software is licensed under Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International
 * For more information, see README.md and LICENSE
 */
