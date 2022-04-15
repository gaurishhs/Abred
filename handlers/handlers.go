package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/apidev234/abred/database"
	"github.com/apidev234/abred/functions"
	"github.com/apidev234/abred/utils"
	"github.com/bwmarrin/discordgo"
)

func VoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if v.BeforeUpdate == nil {
		targetuser, err := s.User(v.UserID)
		if err != nil {
			return
		}

		data, err := database.Database.GetGuildData(v.GuildID)
		if err != nil {
			return
		}
		isgen := utils.IsGenerator(data, v.ChannelID)
		if !isgen {
			return
		}
		if data.DefaultState == "locked" {
			Channel, err := s.GuildChannelCreateComplex(v.GuildID, discordgo.GuildChannelCreateData{
				Name:      fmt.Sprintf("%s's VC", targetuser.Username),
				UserLimit: data.DefaultUserLimit,
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{
						ID:    v.GuildID,
						Type:  discordgo.PermissionOverwriteTypeRole,
						Allow: discordgo.PermissionViewChannel,
						Deny:  discordgo.PermissionVoiceConnect,
					},
				},
				ParentID: data.DefaultCategory,
				Bitrate:  int(data.DefaultBitrate) * 1000,
				Type:     discordgo.ChannelTypeGuildVoice,
			})
			if err != nil {
				return
			}
			err = s.GuildMemberMove(v.GuildID, v.UserID, &Channel.ID)
			if err != nil {
				return
			}
			err = database.Database.NewChannel(v.GuildID, v.UserID, Channel.ID)
			if err != nil {
				return
			}
		} else {
			Channel, err := s.GuildChannelCreateComplex(v.GuildID, discordgo.GuildChannelCreateData{
				Name:      fmt.Sprintf("%s's VC", targetuser.Username),
				UserLimit: data.DefaultUserLimit,
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{
						ID:    v.GuildID,
						Type:  discordgo.PermissionOverwriteTypeRole,
						Allow: discordgo.PermissionViewChannel | discordgo.PermissionVoiceConnect,
						Deny:  0,
					},
				},
				ParentID: data.DefaultCategory,
				Bitrate:  int(data.DefaultBitrate) * 1000,
				Type:     discordgo.ChannelTypeGuildVoice,
			})
			if err != nil {
				return
			}
			s.GuildMemberMove(v.GuildID, v.UserID, &Channel.ID)
			err = database.Database.NewChannel(v.GuildID, v.UserID, Channel.ID)
			if err != nil {
				return
			}
		}
	} else if v.BeforeUpdate.ChannelID != "" && v.ChannelID == "" {
		ischannel := database.Database.IsAChannel(v.GuildID, v.BeforeUpdate.ChannelID)
		if !ischannel {
			return
		}
		empty := utils.CheckIfChannelEmpty(v.BeforeUpdate.ChannelID, v.GuildID, s)
		if !empty {
			return
		}
		err := database.Database.DeleteChannel(v.GuildID, v.BeforeUpdate.ChannelID)
		if err != nil {
			return
		}
		_, err = s.ChannelDelete(v.BeforeUpdate.ChannelID)
		if err != nil {
			return
		}
	}
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	data, err := database.Database.GetGuildData(m.GuildID)
	if err != nil {
		return
	}
	if m.ChannelID != data.HelpChannel {
		return
	}
	var (
		fields = strings.Fields(strings.TrimSpace(m.Content))
	)

	if len(fields) == 0 {
		return
	}

	switch fields[0] {
	case "lock":
		channel := utils.GetUserVoiceChannel(s, m.Author.ID, m.GuildID)
		if channel == nil {
			utils.ComposeErrorMessage(s, m.ChannelID, "Couldn't retrieve your voice channel", m.ID)
			return
		}
		err := functions.LockChannel(m.Author.ID, *channel, m.GuildID, s)
		if err != nil {
			utils.ComposeErrorMessage(s, m.ChannelID, err.Error(), m.ID)
			return
		}
		utils.ComposeSuccessMessage(s, m.ChannelID, fmt.Sprintf("Successfully Locked <#%s>", *channel), m.ID)
		return
	case "unlock":
		channel := utils.GetUserVoiceChannel(s, m.Author.ID, m.GuildID)
		if channel == nil {
			utils.ComposeErrorMessage(s, m.ChannelID, "Couldn't retrieve your voice channel", m.ID)
			return
		}
		err := functions.UnlockChannel(m.Author.ID, *channel, m.GuildID, s)
		if err != nil {
			utils.ComposeErrorMessage(s, m.ChannelID, err.Error(), m.ID)
			return
		}
		utils.ComposeSuccessMessage(s, m.ChannelID, fmt.Sprintf("Successfully unlocked <#%s>", *channel), m.ID)
		return
	case "hide":
		channel := utils.GetUserVoiceChannel(s, m.Author.ID, m.GuildID)
		if channel == nil {
			utils.ComposeErrorMessage(s, m.ChannelID, "Couldn't retrieve your voice channel", m.ID)
			return
		}
		err := functions.HideChannel(m.Author.ID, *channel, m.GuildID, s)
		if err != nil {
			utils.ComposeErrorMessage(s, m.ChannelID, err.Error(), m.ID)
			return
		}
		utils.ComposeSuccessMessage(s, m.ChannelID, fmt.Sprintf("Successfully hid <#%s>", *channel), m.ID)
		return
	case "show":
		channel := utils.GetUserVoiceChannel(s, m.Author.ID, m.GuildID)
		if channel == nil {
			utils.ComposeErrorMessage(s, m.ChannelID, "Couldn't retrieve your voice channel", m.ID)
			return
		}
		err := functions.ShowChannel(m.Author.ID, *channel, m.GuildID, s)
		if err != nil {
			utils.ComposeErrorMessage(s, m.ChannelID, err.Error(), m.ID)
			return
		}
		utils.ComposeSuccessMessage(s, m.ChannelID, fmt.Sprintf("Successfully changed <#%s>'s state to publicly visible", *channel), m.ID)
		return
	case "incl":
		channel := utils.GetUserVoiceChannel(s, m.Author.ID, m.GuildID)
		if channel == nil {
			utils.ComposeErrorMessage(s, m.ChannelID, "Couldn't retrieve your voice channel", m.ID)
			return
		}
		err := functions.IncreaseLimit(m.Author.ID, *channel, m.GuildID, s)
		if err != nil {
			utils.ComposeErrorMessage(s, m.ChannelID, err.Error(), m.ID)
			return
		}
		utils.ComposeSuccessMessage(s, m.ChannelID, fmt.Sprintf("Successfully increased <#%s>'s user limit", *channel), m.ID)
		return
	case "decl":
		channel := utils.GetUserVoiceChannel(s, m.Author.ID, m.GuildID)
		if channel == nil {
			utils.ComposeErrorMessage(s, m.ChannelID, "Couldn't retrieve your voice channel", m.ID)
			return
		}
		err := functions.DecreaseLimit(m.Author.ID, *channel, m.GuildID, s)
		if err != nil {
			utils.ComposeErrorMessage(s, m.ChannelID, err.Error(), m.ID)
			return
		}
		utils.ComposeSuccessMessage(s, m.ChannelID, fmt.Sprintf("Successfully decreased <#%s>'s user limit", *channel), m.ID)
		return
	case "bitrate":
		channel := utils.GetUserVoiceChannel(s, m.Author.ID, m.GuildID)
		if channel == nil {
			utils.ComposeErrorMessage(s, m.ChannelID, "Couldn't retrieve your voice channel", m.ID)
			return
		}
		if len(fields) < 2 {
			utils.ComposeErrorMessage(s, m.ChannelID, "This Command requires a bitrate value as well", m.ID)
			return
		}
		bit, err := strconv.Atoi(fields[1])
		if err != nil {
			return
		}
		if bit < 8 || bit > 96 {
			utils.ComposeErrorMessage(s, m.ChannelID, "The bitrate should be 8 or more than that and not more than 96", m.ID)
			return
		}
		err = functions.SetBitrate(m.Author.ID, *channel, m.GuildID, s, bit)
		if err != nil {
			utils.ComposeErrorMessage(s, m.ChannelID, err.Error(), m.ID)
			return
		}
		utils.ComposeSuccessMessage(s, m.ChannelID, fmt.Sprintf("Successfully changed <#%s>'s bitrate", *channel), m.ID)
		return
	}
	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		return
	}
}

func Ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStreamingStatus(0, fmt.Sprintf("/help | Shard %d", s.ShardID), "https://twitch.tv/shinchanforever")
}

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if i.ApplicationCommandData().Name == "setup" {
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				return
			}
			if channel.Type != discordgo.ChannelTypeGuildText {
				utils.ComposeIntErrorMessage(s, i.Interaction, "This Command is guild only")
				return
			}
			Channel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
				Name: "tempvcs-help",
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{
						ID:    i.GuildID,
						Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
						Deny:  discordgo.PermissionCreatePublicThreads | discordgo.PermissionCreatePrivateThreads | discordgo.PermissionAddReactions,
						Type:  discordgo.PermissionOverwriteTypeRole,
					},
				},
				RateLimitPerUser: 5,
			})
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "<:Cofo_Fail:940541495491051560> __**ERROR**__ <:Cofo_Fail:940541495491051560>",
								Description: "Couldn't create the channel, Try checking the permissions.",
								Color:       16711680,
							},
						},
					},
				})
				return
			}
			err = database.Database.SetupChannel(i.GuildID, Channel.ID)
			if err != nil {
				return
			}
			msg, err := s.ChannelMessageSendComplex(Channel.ID, &discordgo.MessageSend{
				Embed: &discordgo.MessageEmbed{
					Title:       "Abred VC Controller",
					Description: "[Invite](https://abred.bar/invite) | [Support](https://abred.bar/support) | [Documentation](https://docs.abred.bar) | [Website](https://abred.bar)\n\n **Note:** If you just type in the commands here without any other thing it will work!\n\n <:lock:950769296752132177> [`Locks`](https://docs.abredbot.bar/controller#lock) your vc\n<:unlock:950769348337877032> [`Unlocks`](https://docs.abredbot.bar/controller#unlock) your vc\n<:push_to_talk:950769457796612156> Toggles [`Push to talk`](https://docs.abred.bar/controller#push-to-talk) for your vc\n<:increase:950769641737818162> Increases [`user limit`](https://docs.abred.bar/controller#increase-limit) for your vc\n<:decrease:950769557323259914> Decreases [`user limit`](https://docs.abred.bar/controller#decrease-limit) for your vc\n<:hide:950769392063508501> [`Hides`](https://docs.abred.bar/controller#hide) your voice channel from other members\n<:show:957916957422788618> [`Show`](https://docs.abred.bar/controller#show) your voice channel to other members\n<:question:948855713541812224> Get [`help`](https://docs.abred.bar/controller#help) with usage of the controller",
					Color:       65523,
				},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Style: discordgo.SecondaryButton,
								Emoji: discordgo.ComponentEmoji{
									Name: "lock",
									ID:   "950769296752132177",
								},
								CustomID: "lock",
							},
							discordgo.Button{
								Style: discordgo.SecondaryButton,
								Emoji: discordgo.ComponentEmoji{
									Name: "unlock",
									ID:   "950769348337877032",
								},
								CustomID: "unlock",
							},
							discordgo.Button{
								Style: discordgo.SecondaryButton,
								Emoji: discordgo.ComponentEmoji{
									Name: "push_to_talk",
									ID:   "950769457796612156",
								},
								CustomID: "push-to-talk",
							},
							discordgo.Button{
								Style: discordgo.SecondaryButton,
								Emoji: discordgo.ComponentEmoji{
									Name: "increase",
									ID:   "950769641737818162",
								},
								CustomID: "increase-limit",
							},
							discordgo.Button{
								Emoji: discordgo.ComponentEmoji{
									ID:   "950769557323259914",
									Name: "decrease",
								},
								Style:    discordgo.SecondaryButton,
								CustomID: "decrease-limit",
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Style: discordgo.SecondaryButton,
								Emoji: discordgo.ComponentEmoji{
									ID:   "950769392063508501",
									Name: "hide",
								},
								CustomID: "hide",
							},
							discordgo.Button{
								Style: discordgo.SecondaryButton,
								Emoji: discordgo.ComponentEmoji{
									ID:   "957916957422788618",
									Name: "show",
								},
								CustomID: "show",
							},
							discordgo.Button{
								Style: discordgo.SecondaryButton,
								Emoji: discordgo.ComponentEmoji{
									ID:   "948855713541812224",
									Name: "question",
								},
								CustomID: "controller-help",
							},
							discordgo.Button{
								Style: discordgo.LinkButton,
								URL:   "https://abred.bar/manage",
								Label: "Manage Voice",
							},
						},
					},
				},
			})
			if err != nil {
				fmt.Print(err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "<:Cofo_Fail:940541495491051560> __**ERROR**__ <:Cofo_Fail:940541495491051560>",
								Description: "Couldn't send the message to channel, Try checking the permissions.",
								Color:       16711680,
							},
						},
					},
				})
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Description: "Successfully created the controller",
							Color:       8275449,
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "Jump To Channel",
									Value: fmt.Sprintf("[`Click Here`](https://discord.com/channels/%s/%s/%s)", i.GuildID, Channel.ID, msg.ID),
								},
							},
						},
					},
				},
			})
		} else if i.ApplicationCommandData().Name == "help" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 64,
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Abred Help Panel",
							Description: "Quickly get started with abred!\n\nThe bot has only a few commands\n`/help` - Returns this help message\n`/config` - Returns the guild's configuration\n`/newgen` - Create a new generator for the guild\n`/delgen` - Delete a previously registered generator\n`/default-state` - Configure the default state to use while creating channels\n`/default-user-limit` - Configure the default user limit to use while creating channels\n`/default-bitrate` - Configure the default bitrate to use while creating channels\n`/default-category` - Configure the default parent category under which all channels will be created",
							Footer: &discordgo.MessageEmbedFooter{
								Text: "Thanks For Using Abred :)",
							},
							Color: 1357731,
						},
					},
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Label: "Invite",
									Style: discordgo.LinkButton,
									URL:   "https://abred.bar/invite",
								},
								discordgo.Button{
									Label: "Support",
									Style: discordgo.LinkButton,
									URL:   "https://abred.bar/support",
								},
								discordgo.Button{
									Label: "Manage",
									Style: discordgo.LinkButton,
									URL:   "https://abred.bar/login",
								},
								discordgo.Button{
									Label: "Documentation",
									Style: discordgo.LinkButton,
									URL:   "https://docs.abred.bar",
								},
								discordgo.Button{
									Label: "Commands",
									Style: discordgo.LinkButton,
									URL:   "https://docs.abred.bar/commands",
								},
							},
						},
					},
				},
			})
		} else if i.ApplicationCommandData().Name == "config" {
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				return
			}
			if channel.Type != discordgo.ChannelTypeGuildText {
				utils.ComposeIntErrorMessage(s, i.Interaction, "This Command is guild only")
				return
			}
			guildData, err := database.Database.GetGuildData(i.GuildID)
			if err != nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Could Not Retrieve Data For This Guild")
				return
			}
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Configuration For Guild",
							Description: fmt.Sprintf("__***Statistics***__: \nChannels: %d\nGenerators: %d\nGenerator Channels: %s\n\n__***Config***__:\nDefault State: %s\nDefault User Limit: %d\nDefault Category: %s\nController channel: %s\nDefault Bitrate: %d", len(guildData.Channels), len(guildData.Generators), strings.Join(utils.WrapGenerators(guildData.Generators), ","), guildData.DefaultState, guildData.DefaultUserLimit, utils.GetChannelName(guildData.DefaultCategory, s), utils.GetChannelName(guildData.HelpChannel, s), guildData.DefaultBitrate),
							Color:       3787460,
						},
					},
				},
			})
			if err != nil {
				return
			}
		} else if i.ApplicationCommandData().Name == "newgen" {
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				return
			}
			if channel.Type != discordgo.ChannelTypeGuildText {
				utils.ComposeIntErrorMessage(s, i.Interaction, "This Command is guild only")
				return
			}
			cid := i.ApplicationCommandData().Options[0].ChannelValue(s)
			guildData, err := database.Database.GetGuildData(i.GuildID)
			if err != nil {
				return
			}
			if utils.IsGenerator(guildData, cid.ID) {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Channel is already registered!")
				return
			}
			if len(guildData.Generators) >= 5 {
				utils.ComposeIntErrorMessage(s, i.Interaction, "You cannot add more than 5 generators!")
				return
			}
			err = database.Database.NewGenerator(i.GuildID, cid.ID)
			if err != nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Could not register a new generator")
				return
			}
			utils.ComposeIntSuccessMessage(s, i.Interaction, fmt.Sprintf("Successfully registered <#%s> as a generator", cid.ID))
		} else if i.ApplicationCommandData().Name == "delgen" {
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				return
			}
			if channel.Type != discordgo.ChannelTypeGuildText {
				utils.ComposeIntErrorMessage(s, i.Interaction, "This Command is guild only")
				return
			}
			cid := i.ApplicationCommandData().Options[0].ChannelValue(s)
			guildData, err := database.Database.GetGuildData(i.GuildID)
			if err != nil {
				return
			}
			if !utils.IsGenerator(guildData, cid.ID) {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Channel is not registered as a generator!")
				return
			}
			err = database.Database.DeleteGenerator(i.GuildID, cid.ID)
			if err != nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Could not delete generator")
				return
			}
			utils.ComposeIntSuccessMessage(s, i.Interaction, fmt.Sprintf("Successfully removed <#%s> as a generator", cid.ID))
		} else if i.ApplicationCommandData().Name == "default-state" {
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				return
			}
			if channel.Type != discordgo.ChannelTypeGuildText {
				utils.ComposeIntErrorMessage(s, i.Interaction, "This Command is guild only")
				return
			}
			strong := i.ApplicationCommandData().Options[0].StringValue()
			err = database.Database.SetDefaultState(i.GuildID, strong)
			if err != nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Could not change the default state")
				return
			}
			utils.ComposeIntSuccessMessage(s, i.Interaction, "Successfully changed the default state")
		} else if i.ApplicationCommandData().Name == "default-user-limit" {
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				return
			}
			if channel.Type != discordgo.ChannelTypeGuildText {
				utils.ComposeIntErrorMessage(s, i.Interaction, "This Command is guild only")
				return
			}
			limit := i.ApplicationCommandData().Options[0].IntValue()
			err = database.Database.SetDefaultUl(i.GuildID, int(limit))
			if err != nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Could not change the default user limit")
				return
			}
			utils.ComposeIntSuccessMessage(s, i.Interaction, "Successfully changed the default user limit")
		} else if i.ApplicationCommandData().Name == "default-bitrate" {
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				return
			}
			if channel.Type != discordgo.ChannelTypeGuildText {
				utils.ComposeIntErrorMessage(s, i.Interaction, "This Command is guild only")
				return
			}
			limit := i.ApplicationCommandData().Options[0].IntValue()
			err = database.Database.SetDefaultBitrate(i.GuildID, int(limit))
			if err != nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Could not change the default bitrate")
				return
			}
			utils.ComposeIntSuccessMessage(s, i.Interaction, "Successfully changed the default bitrate")
		} else if i.ApplicationCommandData().Name == "default-category" {
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				return
			}
			if channel.Type != discordgo.ChannelTypeGuildText {
				utils.ComposeIntErrorMessage(s, i.Interaction, "This Command is guild only")
				return
			}
			cid := i.ApplicationCommandData().Options[0].ChannelValue(s)
			err = database.Database.SetDefaultCategory(i.GuildID, cid.ID)
			if err != nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Could not set default category")
				return
			}
			utils.ComposeIntSuccessMessage(s, i.Interaction, fmt.Sprintf("Successfully made <#%s> the default category", cid.ID))
		}
	case discordgo.InteractionMessageComponent:
		switch i.MessageComponentData().CustomID {
		case "hide":
			channel := utils.GetUserVoiceChannel(s, i.Member.User.ID, i.GuildID)
			if channel == nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Couldn't retrieve your voice channel")
				return
			}
			err := functions.HideChannel(i.Member.User.ID, *channel, i.GuildID, s)
			if err != nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, err.Error())
				return
			}
			utils.ComposeIntSuccessMessage(s, i.Interaction, fmt.Sprintf("Successfully hid <#%s>", *channel))
		case "show":
			channel := utils.GetUserVoiceChannel(s, i.Member.User.ID, i.GuildID)
			if channel == nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Couldn't retrieve your voice channel")
				return
			}
			err := functions.ShowChannel(i.Member.User.ID, *channel, i.GuildID, s)
			if err != nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, err.Error())
				return
			}
			utils.ComposeIntSuccessMessage(s, i.Interaction, fmt.Sprintf("Successfully made <#%s> publicly visible.", *channel))
		case "lock":
			channel := utils.GetUserVoiceChannel(s, i.Member.User.ID, i.GuildID)
			if channel == nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Couldn't retrieve your voice channel")
				return
			}
			err := functions.LockChannel(i.Member.User.ID, *channel, i.GuildID, s)
			if err != nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, err.Error())
				return
			}
			utils.ComposeIntSuccessMessage(s, i.Interaction, fmt.Sprintf("Successfully locked <#%s>", *channel))
		case "unlock":
			channel := utils.GetUserVoiceChannel(s, i.Member.User.ID, i.GuildID)
			if channel == nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Couldn't retrieve your voice channel")
				return
			}
			err := functions.UnlockChannel(i.Member.User.ID, *channel, i.GuildID, s)
			if err != nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, err.Error())
				return
			}
			utils.ComposeIntSuccessMessage(s, i.Interaction, fmt.Sprintf("Successfully unlocked <#%s>", *channel))
		case "increase-limit":
			channel := utils.GetUserVoiceChannel(s, i.Member.User.ID, i.GuildID)
			if channel == nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Couldn't retrieve your voice channel")
				return
			}
			err := functions.IncreaseLimit(i.Member.User.ID, *channel, i.GuildID, s)
			if err != nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, err.Error())
				return
			}
			utils.ComposeIntSuccessMessage(s, i.Interaction, fmt.Sprintf("Successfully increased <#%s>'s limit", *channel))
		case "decrease-limit":
			channel := utils.GetUserVoiceChannel(s, i.Member.User.ID, i.GuildID)
			if channel == nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Couldn't retrieve your voice channel")
				return
			}
			err := functions.DecreaseLimit(i.Member.User.ID, *channel, i.GuildID, s)
			if err != nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, err.Error())
				return
			}
			utils.ComposeIntSuccessMessage(s, i.Interaction, fmt.Sprintf("Successfully decreased <#%s>'s limit", *channel))
		case "controller-help":
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 64,
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Controller Help",
							Color:       9396735,
							Description: "For Help regarding vc controllers, Please click [here](https://abred.bar/topics/controller)",
						},
					},
				},
			})
			if err != nil {
				return
			}
		case "push-to-talk":
			channel := utils.GetUserVoiceChannel(s, i.Member.User.ID, i.GuildID)
			if channel == nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, "Couldn't retrieve your voice channel")
				return
			}
			err := functions.TogglePushToTalk(i.Member.User.ID, *channel, i.GuildID, s)
			if err != nil {
				utils.ComposeIntErrorMessage(s, i.Interaction, err.Error())
				return
			}
			utils.ComposeIntSuccessMessage(s, i.Interaction, fmt.Sprintf("Successfully toggled push to talk for <#%s>", *channel))
		}
	}
}

func GuildCreate(s *discordgo.Session, evt *discordgo.GuildCreate) {
	_, err := database.Database.GetGuildData(evt.Guild.ID)
	if err != nil {
		s.State.GuildAdd(evt.Guild)
		database.Database.CreateGuild(evt.Guild)
	}
}

func GuildDelete(s *discordgo.Session, evt *discordgo.GuildDelete) {
	if evt.Guild.Unavailable {
		return
	}
	database.Database.DeleteGuild(evt.Guild.ID)
}

var (
	muteX = &sync.RWMutex{}
)

/*
 * Abred
 * Copyright (C) 2022 ApiDev234
 * This software is licensed under Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International
 * For more information, see README.md and LICENSE
 */
