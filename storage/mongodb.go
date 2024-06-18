package storage

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tmv/project"
	"tmv/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	Client            *mongo.Client
	UserCollection    *mongo.Collection
	ProjectCollection *mongo.Collection
	TaskCollection    *mongo.Collection
}

func NewMongoStorage(uri string, dbName string, userCollectionName, projectCollectionName, taskCollectionName string) (*MongoStorage, error) {
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

	userCollection := client.Database(dbName).Collection(userCollectionName)
	taskCollection := client.Database(dbName).Collection(taskCollectionName)
	projectCollection := client.Database(dbName).Collection(projectCollectionName)

	return &MongoStorage{
		Client:            client,
		UserCollection:    userCollection,
		ProjectCollection: projectCollection,
		TaskCollection:    taskCollection,
	}, nil
}

func (m *MongoStorage) GetAllUsers() map[primitive.ObjectID]user.User {
	users := make(map[primitive.ObjectID]user.User)

	cursor, err := m.UserCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return users // return empty map if there's an error
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var us user.User
		err := cursor.Decode(&us)
		if err != nil {
			continue // skip decoding errors
		}
		users[us.Id] = us
	}

	return users
}
func (m *MongoStorage) GetUser(userId primitive.ObjectID) (user.User, error) {
	var usr user.User

	filter := bson.D{{Key: "_id", Value: userId}}
	err := m.UserCollection.FindOne(context.TODO(), filter).Decode(&usr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return usr, errors.New("user not found")
		}
		return usr, err
	}

	return usr, nil
}
func (m *MongoStorage) InsertUser(u *user.User) error {
	u.Id = primitive.NewObjectID()

	_, err := m.UserCollection.InsertOne(context.TODO(), u)
	return err
}
func (m *MongoStorage) UpdateUser(userId primitive.ObjectID, e *user.User) error {
	filter := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$set", Value: e}}

	_, err := m.UserCollection.UpdateOne(context.TODO(), filter, update)
	return err
}
func (m *MongoStorage) DeleteUser(userId primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: userId}}

	_, err := m.UserCollection.DeleteOne(context.TODO(), filter)
	return err
}

func (m *MongoStorage) GetAllProjects() map[primitive.ObjectID]project.Project {
	projects := make(map[primitive.ObjectID]project.Project)

	cursor, err := m.ProjectCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return projects // return empty map if there's an error
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var proj project.Project
		err := cursor.Decode(&proj)
		if err != nil {
			continue // skip decoding errors
		}
		projects[proj.Id] = proj
	}
	return projects
}
func (m *MongoStorage) GetProjectByUser(userId primitive.ObjectID) ([]project.Project, error) {
	var projects []project.Project

	// Создаем фильтр для поиска проектов по userId
	filter := bson.D{{Key: "userId", Value: userId}}

	// Выполняем запрос к коллекции проектов
	cursor, err := m.ProjectCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Обрабатываем результаты запроса
	for cursor.Next(context.TODO()) {
		var proj project.Project
		if err := cursor.Decode(&proj); err != nil {
			return nil, err
		}
		projects = append(projects, proj)
	}

	// Проверяем наличие ошибок при работе с курсором
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// Возвращаем результаты
	return projects, nil
}
func (m *MongoStorage) GetProject(userId, projectId primitive.ObjectID) (*project.Project, error) {
	var proj project.Project

	// Создаем фильтр для поиска проекта по userId и projectId
	filter := bson.D{
		{Key: "_id", Value: projectId},
		{Key: "userId", Value: userId},
	}

	// Выполняем запрос к коллекции проектов
	err := m.ProjectCollection.FindOne(context.TODO(), filter).Decode(&proj)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("проект не найден")
		}
		return nil, err
	}

	// Возвращаем найденный проект
	return &proj, nil
}
func (m *MongoStorage) DeleteProject(id primitive.ObjectID) error {
	// Найти проект по ID, чтобы получить userID
	var project project.Project
	err := m.ProjectCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&project)
	if err != nil {
		return err
	}
	// Удалить проект из коллекции проектов
	_, err = m.ProjectCollection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}
	// Удалить ID проекта из массива projects в документе пользователя
	filter := bson.D{{Key: "_id", Value: project.UserID}}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "projects", Value: id},
		}},
	}

	_, err = m.UserCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
func (m *MongoStorage) DeleteProjects(userID primitive.ObjectID, projectIDs []primitive.ObjectID) error {
	// Удалить проекты из коллекции проектов
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: projectIDs}}}}
	_, err := m.ProjectCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}

	// Обновить документ пользователя, удалив ID проектов из массива projects
	userFilter := bson.D{{Key: "_id", Value: userID}}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "projects", Value: bson.D{{Key: "$in", Value: projectIDs}}},
		}},
	}

	_, err = m.UserCollection.UpdateOne(context.TODO(), userFilter, update)
	if err != nil {
		return err
	}

	return nil
}
func (m *MongoStorage) UpdateProject(projectID primitive.ObjectID, updateFields bson.M) error {
	filter := bson.D{{Key: "_id", Value: projectID}}
	update := bson.D{{Key: "$set", Value: updateFields}}

	_, err := m.ProjectCollection.UpdateOne(context.TODO(), filter, update)
	return err
}
func (m *MongoStorage) InsertProject(p *project.Project, userID primitive.ObjectID) error {
	// Генерируем новый ObjectID для проекта
	p.Id = primitive.NewObjectID()
	// Присваиваем ObjectID пользователя проекту
	p.UserID = userID

	// Вставляем документ проекта в коллекцию проектов (ProjectCollection)
	_, err := m.ProjectCollection.InsertOne(context.TODO(), p)
	if err != nil {
		return err
	}

	// Обновляем документ пользователя в коллекции пользователей (UserCollection),
	// чтобы добавить ID проекта в массив projects
	filter := bson.D{{Key: "_id", Value: userID}}

	// Проверяем, существует ли поле projects
	userUpdateResult := m.UserCollection.FindOne(context.TODO(), filter)

	var userDoc map[string]interface{}
	err = userUpdateResult.Decode(&userDoc)
	if err != nil {
		// Если пользователь не найден, возвращаем ошибку
		return err
	}
	if _, ok := userDoc["projects"]; !ok {
		// Если поле projects не существует, инициализируем его как пустой массив
		_, err = m.UserCollection.UpdateOne(
			context.TODO(),
			filter,
			bson.D{
				{Key: "$set", Value: bson.D{{Key: "projects", Value: []primitive.ObjectID{}}}},
			},
		)
		if err != nil {
			return err
		}
	}
	// Добавляем новый проект в массив projects
	_, err = m.UserCollection.UpdateOne(
		context.TODO(),
		filter,
		bson.D{
			{Key: "$addToSet", Value: bson.D{{Key: "projects", Value: p.Id}}},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoStorage) GetTasksByProject(projectId primitive.ObjectID) ([]project.Task, error) {
	var tasks []project.Task
	filter := bson.D{{Key: "projectId", Value: projectId}}

	cursor, err := m.TaskCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var task project.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
func (m *MongoStorage) InsertTask(t *project.Task, projectId primitive.ObjectID) error {

	t.ID = primitive.NewObjectID()
	t.ProjectID = projectId

	_, err := m.TaskCollection.InsertOne(context.TODO(), t)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: projectId}}

	projectUpdateResult := m.ProjectCollection.FindOne(context.TODO(), filter)

	var projectDoc map[string]interface{}
	err = projectUpdateResult.Decode(&projectDoc)
	if err != nil {
		// Если пользователь не найден, возвращаем ошибку
		return err
	}
	if _, ok := projectDoc["tasks"]; !ok || projectDoc["tasks"] == nil {
		// Если поле tasks не существует, инициализируем его как пустой массив
		_, err = m.ProjectCollection.UpdateOne(
			context.TODO(),
			filter,
			bson.D{
				{Key: "$set", Value: bson.D{{Key: "tasks", Value: []primitive.ObjectID{}}}},
			},
		)
		if err != nil {
			return err
		}
	}

	// Добавляем новый проект в массив tasks
	_, err = m.ProjectCollection.UpdateOne(
		context.TODO(),
		filter,
		bson.D{
			{Key: "$addToSet", Value: bson.D{{Key: "tasks", Value: t.ID}}},
		},
	)
	if err != nil {
		return err
	}
	return nil
}
func (m *MongoStorage) GetTask(projectId, taskId primitive.ObjectID) (*project.Task, error) {
	// Фильтр для поиска задачи по taskId и projectId
	filter := bson.D{
		{Key: "_id", Value: taskId},
		{Key: "projectId", Value: projectId},
	}

	var task project.Task
	err := m.TaskCollection.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}
func (m *MongoStorage) GetAllTasks() map[primitive.ObjectID]project.Task {
	tasks := make(map[primitive.ObjectID]project.Task)

	cursor, err := m.TaskCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return tasks // return empty map if there's an error
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var task project.Task
		err := cursor.Decode(&task)
		if err != nil {
			continue // skip decoding errors
		}
		tasks[task.ID] = task
	}
	return tasks
}
func (m *MongoStorage) DeleteTask(projectId, taskId primitive.ObjectID) error {
	filter := bson.D{
		{Key: "_id", Value: taskId},
		{Key: "projectId", Value: projectId},
	}

	// Удаление задачи из коллекции задач
	_, err := m.TaskCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	// Фильтр для обновления проекта
	projectFilter := bson.D{{Key: "_id", Value: projectId}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "tasks", Value: taskId}}}}

	// Удаление ID задачи из массива tasks в проекте
	_, err = m.ProjectCollection.UpdateOne(context.TODO(), projectFilter, update)
	if err != nil {
		return err
	}

	return nil
}
func (m *MongoStorage) DeleteTasks(projectId primitive.ObjectID, taskIds []primitive.ObjectID) error {
	// Фильтр для удаления задач по projectId и массиву taskIds
	filter := bson.D{
		{Key: "projectId", Value: projectId},
		{Key: "_id", Value: bson.D{{Key: "$in", Value: taskIds}}},
	}

	// Удаление задач из коллекции задач
	_, err := m.TaskCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}

	// Фильтр для обновления проекта
	projectFilter := bson.D{{Key: "_id", Value: projectId}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "tasks", Value: bson.D{{Key: "$in", Value: taskIds}}}}}}

	// Удаление ID задач из массива tasks в проекте
	_, err = m.ProjectCollection.UpdateOne(context.TODO(), projectFilter, update)
	if err != nil {
		return err
	}

	return nil
}
func (m *MongoStorage) UpdateTask(projectId, taskId primitive.ObjectID, updateFields bson.M) error {
	filter := bson.D{
		{Key: "_id", Value: taskId},
		{Key: "projectId", Value: projectId},
	}
	update := bson.D{{Key: "$set", Value: updateFields}}

	_, err := m.TaskCollection.UpdateOne(context.TODO(), filter, update)
	return err
}
