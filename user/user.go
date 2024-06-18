package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name     string               `bson:"name" json:"name"`
	Work     string               `bson:"work" json:"work"`
	Age      int                  `bson:"age" json:"age"`
	Salary   int                  `bson:"salary" json:"salary"`
	Email    string               `bson:"email" json:"email"`
	Projects []primitive.ObjectID `bson:"projects" json:"projects"`
}

func NewUser(name, work string, age, salary int, email string, projects []primitive.ObjectID) *User {
	return &User{
		Name:     name,
		Work:     work,
		Age:      age,
		Salary:   salary,
		Email:    email,
		Projects: projects,
	}
}
