package database

import (
	"context"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Channel struct {
	OwnerID string
	Id      string
}

type Guild struct {
	Generators       []string
	Id               string
	DefaultCategory  string
	DefaultBitrate   int
	DefaultState     string
	DefaultUserLimit int
	HelpChannel      string
	Channels         []Channel
}

func (db *MongoDB) CreateGuild(guild *discordgo.Guild) error {
	_, err := db.Collection.InsertOne(context.Background(), bson.M{
		"id":               guild.ID,
		"defaultState":     "unlocked",
		"defaultUserLimit": 1,
		"defaultBitrate":   64,
		"generators":       bson.A{},
		"channels":         bson.A{},
	})
	if err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) DeleteGuild(guildId string) error {
	_, err := db.Collection.DeleteOne(context.Background(), bson.M{"id": guildId})
	if err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) NewChannel(guildId string, ownerId string, channelId string) error {
	if _, err = db.Collection.UpdateOne(context.Background(), bson.M{"id": guildId}, bson.M{"$push": bson.M{"channels": bson.M{"id": channelId, "ownerid": ownerId}}}, &options.UpdateOptions{}); err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) TransferOwner(guildId string, channelId string, id string) error {
	query := bson.M{
		"id":          guildId,
		"channels.id": channelId,
	}
	update := bson.M{
		"$set": bson.M{
			"channels.$.ownerid": id,
		},
	}
	if _, err = db.Collection.UpdateOne(context.Background(), query, update); err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) SetDefaultState(guildId string, state string) error {
	if state == "locked" {
		if _, err = db.Collection.UpdateOne(context.Background(), bson.M{"id": guildId}, bson.M{"$set": bson.M{"defaultState": "locked"}}, &options.UpdateOptions{}); err != nil {
			return err
		}
	} else {
		if _, err = db.Collection.UpdateOne(context.Background(), bson.M{"id": guildId}, bson.M{"$set": bson.M{"defaultState": "unlocked"}}, &options.UpdateOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func (db *MongoDB) SetDefaultUl(guildId string, limit int) error {
	if _, err = db.Collection.UpdateOne(context.Background(), bson.M{"id": guildId}, bson.M{"$set": bson.M{"defaultUserLimit": limit}}, &options.UpdateOptions{}); err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) SetDefaultCategory(guildId string, cat string) error {
	if _, err = db.Collection.UpdateOne(context.Background(), bson.M{"id": guildId}, bson.M{"$set": bson.M{"defaultCategory": cat}}, &options.UpdateOptions{}); err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) SetDefaultBitrate(guildId string, limit int) error {
	if _, err = db.Collection.UpdateOne(context.Background(), bson.M{"id": guildId}, bson.M{"$set": bson.M{"defaultBitrate": limit}}, &options.UpdateOptions{}); err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) DeleteChannel(guildId string, channelId string) error {
	if _, err = db.Collection.UpdateOne(context.Background(), bson.M{"id": guildId}, bson.M{"$pull": bson.M{"channels": bson.M{"id": channelId}}}, &options.UpdateOptions{}); err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) IsAChannel(guildId string, channelId string) bool {
	guilddata, err := db.GetGuildData(guildId)
	if err != nil {
		return false
	}
	for _, channel := range guilddata.Channels {
		if channel.Id == channelId {
			return true
		}
	}
	return false
}

func (db *MongoDB) NewGenerator(guildid string, channelid string) error {
	if _, err = db.Collection.UpdateOne(context.Background(), bson.M{"id": guildid}, bson.M{"$push": bson.M{"generators": channelid}}, &options.UpdateOptions{}); err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) DeleteGenerator(guildid string, channelid string) error {
	if _, err = db.Collection.UpdateOne(context.Background(), bson.M{"id": guildid}, bson.M{"$pull": bson.M{"generators": channelid}}, &options.UpdateOptions{}); err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) GetChannelData(guildid string, channelid string, s *discordgo.Session) (*Channel, error) {
	guilddata, err := Database.GetGuildData(guildid)
	if err != nil {
		return nil, err
	}
	for _, chonl := range guilddata.Channels {
		if chonl.Id == channelid {
			return &chonl, nil
		}
	}
	return nil, nil
}

func (db *MongoDB) SetupChannel(guildId string, channelId string) error {
	if _, err = db.Collection.UpdateOne(context.Background(), bson.M{"id": guildId}, bson.M{"$set": bson.M{"helpChannel": channelId}}, &options.UpdateOptions{}); err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) GetGuildData(guildid string) (*Guild, error) {
	result := Guild{}
	if err = db.Collection.FindOne(context.TODO(), bson.M{"id": guildid}).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func SetupDB() MongoDB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var db = MongoDB{}

	db.Client, err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URL")))

	if err != nil {
		panic(err)
	}

	err = db.Client.Connect(db.Ctx)
	if err != nil {
		panic(err)
	}

	db.Database = db.Client.Database("abred")
	db.Collection = db.Database.Collection("guilds")

	return db
}

var (
	cancel   func()
	Database = SetupDB()
	err      error
)

type (
	MongoDB struct {
		Collection *mongo.Collection
		Client     *mongo.Client
		Ctx        context.Context
		Database   *mongo.Database
	}
)

/*
 * Abred
 * Copyright (C) 2022 ApiDev234
 * This software is licensed under Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International
 * For more information, see README.md and LICENSE
 */
