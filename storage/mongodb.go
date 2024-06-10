package storage

import (
	"context"
	"errors"
	"time"
	"tmv/employee"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	Client     *mongo.Client
	Collection *mongo.Collection
	Counter    int
}

func NewMongoStorage(uri string, dbName string, collectionName string) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &MongoStorage{
		Client:     client,
		Collection: collection,
		Counter:    1,
	}, nil
}

func (m *MongoStorage) GetAll() map[int]employee.Employee {
	employees := make(map[int]employee.Employee)

	cursor, err := m.Collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return employees // return empty map if there's an error
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var emp employee.Employee
		err := cursor.Decode(&emp)
		if err != nil {
			continue // skip decoding errors
		}
		employees[emp.Id] = emp
	}

	return employees
}

func (m *MongoStorage) Get(id int) (employee.Employee, error) {
	var emp employee.Employee

	filter := bson.D{{Key: "id", Value: id}}
	err := m.Collection.FindOne(context.TODO(), filter).Decode(&emp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return emp, errors.New("employee not found")
		}
		return emp, err
	}

	return emp, nil
}

func (m *MongoStorage) Insert(e *employee.Employee) {
	e.Id = m.Counter
	m.Counter++

	_, _ = m.Collection.InsertOne(context.TODO(), e)
}

func (m *MongoStorage) Update(id int, e *employee.Employee) {
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{{Key: "$set", Value: e}}

	_, _ = m.Collection.UpdateOne(context.TODO(), filter, update)
}

func (m *MongoStorage) Delete(id int) {
	filter := bson.D{{Key: "id", Value: id}}

	_, _ = m.Collection.DeleteOne(context.TODO(), filter)
}
