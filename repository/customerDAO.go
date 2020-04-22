package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/erply/api-go-wrapper/pkg/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (dao *CustomerDao) SaveCustomer(customer *api.Customer) error {
	ctx := context.Background()

	db := dao.dbClient
	coll := db.Database(dao.dbDatabase).Collection(dao.dbCollection)

	// we can inserts many rows at once
	res, err := coll.InsertOne(ctx, customer)
	if err != nil {
		return err
	}

	fmt.Printf("ID: %#v\n", res.InsertedID)
	return nil
}

func (dao *CustomerDao) UpdateCustomer(customer *api.Customer) error {
	ctx := context.Background()

	db := dao.dbClient
	coll := db.Database(dao.dbDatabase).Collection(dao.dbCollection)

	filter := bson.M{"customerid": customer.ID}
	update := bson.M{
		"$set": customer,
	}
	// we can inserts many rows at once
	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("ERROR:UpdateCustomer:", err.Error())
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New(fmt.Sprintf("Customer Id:%d is not available in mongo repository", customer.ID))
	}

	return nil
}

func (dao *CustomerDao) FindCustomer(customerId int) (api.Customer, error) {
	ctx := context.Background()

	db := dao.dbClient
	coll := db.Database(dao.dbDatabase).Collection(dao.dbCollection)

	filter := bson.M{"customerid": customerId}
	// we can inserts many rows at once
	var customer = api.Customer{}
	err := coll.FindOne(ctx, filter).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return api.Customer{}, errors.New(fmt.Sprintf("Customer id %v is not found in customer collection", customerId))
		}
		fmt.Println("ERROR:UpdateCustomer:", err.Error())
		return api.Customer{}, err
	}

	return customer, nil
}
