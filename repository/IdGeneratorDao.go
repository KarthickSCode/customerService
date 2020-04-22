package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// Generate a auto increment version ID for the given key
func (dao *IdGenDao) NextId(key string) (int, error) {

	ctx := context.Background()

	db := dao.dbClient

	coll := db.Database(dao.dbDatabase).Collection(dao.dbCollection)

	result, err := coll.UpdateOne(ctx, bson.M{"key": key}, bson.M{"$inc": bson.M{"counter": 1}})

	if err != nil {
		panic("ERROR:Increment Counter")
	}
	if result.MatchedCount == 0 {
		fmt.Println("ERROR:Key ", key, " is not available")
		var idGen = &IdGenerator{
			Counter: 1,
			Key:     key,
		}
		_, err := coll.InsertOne(ctx, idGen)
		if err != nil {
			log.Fatal(err)
			return -1, err
		}
		return 1, nil
	}

	var idGen IdGenerator
	err = coll.FindOne(ctx, bson.M{"key": key}).Decode(&idGen)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return -1, errors.New("Key is not found in Id Generator collection")
		}
		log.Fatal(err)
		return -1, err
	}

	return idGen.Counter, nil
}
