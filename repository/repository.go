package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"time"
)

type Repository struct {
	dbClient *mongo.Client
	dbServer string
}

type IdGenerator struct {
	Counter int `bson:"counter"`
	Key     string
}

//type Dao struct {
//	dbClient     *mongo.Client
//	dbServer     string
//	dbDatabase   string
//	dbCollection string
//}

type IdGenDao struct {
	dbClient     *mongo.Client
	dbServer     string
	dbDatabase   string
	dbCollection string
}

type CustomerDao struct {
	dbClient     *mongo.Client
	dbServer     string
	dbDatabase   string
	dbCollection string
}

// Setup initializes a mongo client
func Setup(address string) (*Repository, error) {
	ctx := context.Background()
	repo := new(Repository)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	// cancel will be called when setup exits
	defer cancel()
	client, err := mongo.NewClient(options.Client().ApplyURI(address))
	if err != nil {
		return nil, err
	}
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	repo.dbClient = client
	repo.dbServer = address

	return repo, nil
}

func (repo *Repository) GetIdGenDao(dbDatabase string, dbCollection string) *IdGenDao {
	dao := new(IdGenDao)
	dao.dbClient = repo.dbClient
	dao.dbServer = repo.dbServer
	dao.dbDatabase = dbDatabase
	dao.dbCollection = dbCollection

	return dao
}

func (repo *Repository) GetCustomerDao(dbDatabase string, dbCollection string) *CustomerDao {
	dao := new(CustomerDao)
	dao.dbClient = repo.dbClient
	dao.dbServer = repo.dbServer
	dao.dbDatabase = dbDatabase
	dao.dbCollection = dbCollection

	return dao
}
