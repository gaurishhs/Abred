package api

import (
	"github.com/apidev234/abred/database"
	"github.com/apidev234/abred/functions"
	"github.com/apidev234/abred/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func NewGenerator(c *gin.Context, s *discordgo.Session) {
	guildId := c.Query("gid")
	userId := c.Query("uid")
	channelId := c.Query("cid")
	if guildId == "" {
		c.AbortWithStatus(400)
		return
	}
	if userId == "" {
		c.AbortWithStatus(400)
		return
	}
	if channelId == "" {
		c.AbortWithStatus(400)
		return
	}
	perms := utils.HasPerms(s, nil, guildId, userId, discordgo.PermissionManageServer)
	if !perms {
		c.AbortWithStatus(403)
		return
	}
	err := database.Database.NewGenerator(guildId, channelId)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err,
		})
	} else {
		c.JSON(200, gin.H{
			"success": true,
		})
	}
}

func DeleteGenerator(c *gin.Context, s *discordgo.Session) {
	guildId := c.Query("gid")
	userId := c.Query("uid")
	channelId := c.Query("cid")
	if guildId == "" {
		c.AbortWithStatus(400)
		return
	}
	if userId == "" {
		c.AbortWithStatus(400)
		return
	}
	if channelId == "" {
		c.AbortWithStatus(400)
		return
	}
	perms := utils.HasPerms(s, nil, guildId, userId, discordgo.PermissionManageServer)
	if !perms {
		c.AbortWithStatus(403)
		return
	}
	err := database.Database.DeleteGenerator(guildId, channelId)
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"error":   err,
		})
	} else {
		c.JSON(200, gin.H{
			"success": true,
		})
	}
}

func GetGuildData(c *gin.Context, s *discordgo.Session) {
	guildId := c.Query("gid")
	userId := c.Query("uid")
	if guildId == "" {
		c.AbortWithStatus(400)
		return
	}
	if userId == "" {
		c.AbortWithStatus(400)
		return
	}
	perms := utils.HasPerms(s, nil, guildId, userId, discordgo.PermissionManageServer)
	if !perms {
		c.AbortWithStatus(403)
		return
	}
	guilddata, err := database.Database.GetGuildData(guildId)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err,
		})
	}
	c.JSON(200, gin.H{
		"success": true,
		"data":    guilddata,
	})
}

func GetGuilds(c *gin.Context, s *discordgo.Session) {
	var list []string
	for _, guilds := range s.State.Guilds {
		list = append(list, guilds.ID)
	}
	c.JSON(200, gin.H{
		"guilds": list,
	})
}

type Channel struct {
	Id   string
	Name string
}

func GetChannels(c *gin.Context, s *discordgo.Session) {
	guildId := c.Query("gid")
	userId := c.Query("uid")
	if guildId == "" {
		c.AbortWithStatus(400)
		return
	}
	if userId == "" {
		c.AbortWithStatus(400)
		return
	}
	perms := utils.HasPerms(s, nil, guildId, userId, discordgo.PermissionManageServer)
	if !perms {
		c.AbortWithStatus(403)
		return
	}
	guild, err := s.Guild(guildId)
	if err != nil {
		c.JSON(400, gin.H{
			"success": "false",
			"error":   "Invalid Guild",
		})
	}
	var list []Channel
	for _, channel := range guild.Channels {
		if channel.Type == discordgo.ChannelTypeGuildText {
			chonl := Channel{
				Id:   channel.ID,
				Name: channel.Name,
			}
			list = append(list, chonl)
		}
	}
	c.JSON(200, gin.H{
		"success":  true,
		"channels": list,
	})
}

func IncreaseUserLimit(c *gin.Context, s *discordgo.Session) {
	guildId := c.Query("gid")
	userId := c.Query("uid")
	channelId := c.Query("cid")
	if guildId == "" {
		c.AbortWithStatus(400)
		return
	}
	if channelId == "" {
		c.AbortWithStatus(400)
		return
	}
	if userId == "" {
		c.AbortWithStatus(400)
		return
	}
	err := functions.IncreaseLimit(userId, channelId, guildId, s)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

func DecreaseUserLimit(c *gin.Context, s *discordgo.Session) {
	guildId := c.Query("gid")
	userId := c.Query("uid")
	channelId := c.Query("cid")
	if guildId == "" {
		c.AbortWithStatus(400)
		return
	}
	if channelId == "" {
		c.AbortWithStatus(400)
		return
	}
	if userId == "" {
		c.AbortWithStatus(400)
		return
	}
	perms := utils.HasPerms(s, nil, guildId, userId, discordgo.PermissionManageServer)
	if !perms {
		c.AbortWithStatus(403)
		return
	}
	err := functions.DecreaseLimit(userId, channelId, guildId, s)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

func GetChannelName(c *gin.Context, s *discordgo.Session) {
	channelId := c.Query("cid")
	if channelId == "" {
		c.AbortWithStatus(400)
		return
	}
	channel, err := s.Channel(channelId)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
		})
		return
	}
	c.JSON(200, gin.H{
		"success":     true,
		"channelName": channel.Name,
	})
}

func TogglePushToTalk(c *gin.Context, s *discordgo.Session) {
	channelId := c.Query("cid")
	userId := c.Query("uid")
	guildId := c.Query("gid")
	if userId == "" {
		c.AbortWithStatus(400)
		return
	}
	if guildId == "" {
		c.AbortWithStatus(400)
		return
	}
	if channelId == "" {
		c.AbortWithStatus(400)
		return
	}
	err := functions.TogglePushToTalk(userId, channelId, guildId, s)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err,
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

/*
 * Abred
 * Copyright (C) 2022 ApiDev234
 * This software is licensed under Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International
 * For more information, see README.md and LICENSE
 */
