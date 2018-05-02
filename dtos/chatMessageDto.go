package dtos

type ChatMessageDto struct {
	//ID      bson.ObjectId `bson:"_id" json:"id"`
	Author  string `bson:"author" json:"author"`
	Message string `bson:"message" json:"message"`
}
