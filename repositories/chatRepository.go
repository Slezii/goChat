package repositories

import (
	"log"

	. "../dtos"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ChatRepository struct {
	Server   string
	Database string
}

var db *mgo.Database

func (dao *ChatRepository) Connect() {
	session, err := mgo.Dial(dao.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(dao.Database)
}

func (dao *ChatRepository) Insert(message ChatMessageDto) error {
	err := db.C("Messages").Insert(&message)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (dao *ChatRepository) GetLast() ([]ChatMessageDto, error) {
	var messages []ChatMessageDto
	err := db.C("Messages").Find(bson.M{}).All(&messages)
	return messages, err
}
