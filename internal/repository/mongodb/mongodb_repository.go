package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/pkg/logger"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongoDBGenericRepository[T any] interface {
	Create(ctx context.Context, collectionName string, entity *T) error
	CreateMany(ctx context.Context, collectionName string, entities []T) error
	GetByID(ctx context.Context, collectionName string, id string) (*T, error)
	Find(ctx context.Context, collectionName string, filter interface{}, limit, offset int64) ([]T, int64, error)
	FindOneByFilter(ctx context.Context, collectionName string, filter interface{}) (*T, error)
	FindOneAndUpdate(ctx context.Context, collectionName string, filter interface{}, data interface{}) (*T, error)
	Update(ctx context.Context, collectionName string, id string, updateData interface{}) error
	DeleteById(ctx context.Context, collectionName string, id string) error
	BulkUpsert(ctx context.Context, collectionName string, models []domain.BulkUpsert) error
	FindPaginated(ctx context.Context, collectionName string, filter interface{}, findOptions *options.FindOptions) ([]T, int64, error)
}

type mongoGenericRepository[T any] struct {
	db *mongo.Database
}

// NewMongoGenericRepository สร้าง Repository อิสระสำหรับ Entity ใดๆ
func NewMongoGenericRepository[T any](db *mongo.Database) IMongoDBGenericRepository[T] {
	return &mongoGenericRepository[T]{
		db: db,
	}
}

func (r *mongoGenericRepository[T]) Create(ctx context.Context, collectionName string, entity *T) error {
	coll := r.db.Collection(collectionName)

	result, err := coll.InsertOne(ctx, entity)
	if err != nil {
		logger.Error(
			"Failed to insert document into collection: "+collectionName,
			zap.Error(err),
		)
		return domain.ErrDatabase
	}

	logger.Info("Successfully inserted document with ID: " + result.InsertedID.(primitive.ObjectID).Hex())
	return nil
}

func (r *mongoGenericRepository[T]) CreateMany(ctx context.Context, collectionName string, entities []T) error {
	if len(entities) == 0 {
		return nil
	}

	// แปลง []T ให้เป็น []interface{} เนื่องจาก InsertMany ของ Mongo Driver รับ []interface{}
	documents := make([]interface{}, len(entities))
	for i, entity := range entities {
		documents[i] = entity
	}

	collection := r.db.Collection(collectionName)
	_, err := collection.InsertMany(ctx, documents)
	if err != nil {
		logger.Error(
			"Failed to insert many document into collection: "+collectionName,
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (r *mongoGenericRepository[T]) GetByID(ctx context.Context, collectionName string, id string) (*T, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	coll := r.db.Collection(collectionName)
	filter := bson.M{"_id": objID}

	var result T
	err = coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		logger.Error(
			"Failed to find document into collection: "+collectionName,
			zap.Error(err),
		)
		return nil, domain.ErrDatabase
	}

	return &result, nil
}

func (r *mongoGenericRepository[T]) Find(ctx context.Context, collectionName string, filter interface{}, limit, offset int64) ([]T, int64, error) {
	coll := r.db.Collection(collectionName)

	if filter == nil {
		filter = bson.M{}
	}

	total, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error(
			"Failed to count document into collection: "+collectionName,
			zap.Error(err),
		)
		return nil, 0, domain.ErrDatabase
	}

	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(limit)
	}
	if offset >= 0 {
		findOptions.SetSkip(offset)
	}

	cursor, err := coll.Find(ctx, filter, findOptions)
	if err != nil {
		logger.Error(
			"Failed to find document into collection: "+collectionName,
			zap.Error(err),
		)
		return nil, 0, domain.ErrDatabase
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		logger.Error(
			"Failed to decode document into collection: "+collectionName,
			zap.Error(err),
		)
		return nil, 0, domain.ErrDatabase
	}

	return results, total, nil
}

func (r *mongoGenericRepository[T]) FindOneByFilter(ctx context.Context, collectionName string, filter interface{}) (*T, error) {
	coll := r.db.Collection(collectionName)

	if filter == nil {
		filter = bson.M{}
	}

	var result T
	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		logger.Error(
			"Failed to find document in collection: "+collectionName,
			zap.Error(err),
		)
		return nil, domain.ErrDatabase
	}

	return &result, nil
}

func (r *mongoGenericRepository[T]) FindOneAndUpdate(ctx context.Context, collectionName string, filter interface{}, data interface{}) (*T, error) {
	coll := r.db.Collection(collectionName)

	if filter == nil {
		filter = bson.M{}
	}

	if data == nil {
		data = bson.M{}
	}

	var result T
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := coll.FindOneAndUpdate(ctx, filter, data, opts).Decode(&result)
	if err != nil {
		logger.Error(
			"Failed to find document and updating in collection: "+collectionName,
			zap.Error(err),
		)
		return nil, domain.ErrDatabase
	}

	return &result, nil
}

func (r *mongoGenericRepository[T]) Update(ctx context.Context, collectionName string, id string, updateData interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	coll := r.db.Collection(collectionName)
	filter := bson.M{"_id": objID}
	updating := bson.M{
		"$set": updateData,
	}
	res, err := coll.UpdateOne(ctx, filter, updating)
	if err != nil {

		logger.Error(
			"Failed to update document into collection: "+collectionName,
			zap.Error(err),
		)
		return domain.ErrDatabase
	}

	if res.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	logger.Info("Successfully updated document in collection: " + collectionName)
	return nil
}

func (r *mongoGenericRepository[T]) DeleteById(ctx context.Context, collectionName string, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	coll := r.db.Collection(collectionName)
	filter := bson.M{"_id": objID}

	res, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		logger.Error(
			"Failed to delete document into collection: "+collectionName,
			zap.Error(err),
		)
		return domain.ErrDatabase
	}

	if res.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	logger.Info("Successfully deleted document from collection: " + collectionName)
	return nil
}

func (r *mongoGenericRepository[T]) BulkUpsert(ctx context.Context, collectionName string, models []domain.BulkUpsert) error {
	if len(models) == 0 {
		return nil
	}

	writes := make([]mongo.WriteModel, len(models))

	for i, m := range models {
		model := mongo.NewUpdateOneModel().
			SetFilter(m.Filter).
			SetUpdate(m.Update).
			SetUpsert(true)

		writes[i] = model
	}

	coll := r.db.Collection(collectionName)

	res, err := coll.BulkWrite(ctx, writes)
	if err != nil {
		logger.Error(
			"Failed to BulkUpsert document for collection: "+collectionName,
			zap.Error(err),
		)
		return domain.ErrDatabase
	}

	logger.Info(fmt.Sprintf("BulkUpsert finished in %s. Inserted: %d, Modified: %d, Upserted: %d",
		collectionName, res.InsertedCount, res.ModifiedCount, res.UpsertedCount))

	return nil
}

// FindPaginated รองรับ Generic Search, Filtering, และ Sorting ทุกแบบ
func (r *mongoGenericRepository[T]) FindPaginated(ctx context.Context, collectionName string, filter interface{}, findOpts *options.FindOptions) ([]T, int64, error) {
	coll := r.db.Collection(collectionName)

	if filter == nil {
		filter = bson.M{}
	}

	// 1. นับจำนวน Document ทั้งหมดที่ตรงตาม Filter
	total, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error(
			"Failed to count document into collection: "+collectionName,
			zap.Error(err),
		)
		return nil, 0, domain.ErrDatabase
	}

	// 2. ดึง ข้อมูลตาม Filter และ Options (Limit, Skip, Sort)
	cursor, err := coll.Find(ctx, filter, findOpts)
	if err != nil {
		logger.Error(
			"Failed to find document into collection: "+collectionName,
			zap.Error(err),
		)
		return nil, 0, domain.ErrDatabase
	}
	defer cursor.Close(ctx)

	results := make([]T, 0)
	if err := cursor.All(ctx, &results); err != nil {
		logger.Error(
			"Failed to decode document into collection: "+collectionName,
			zap.Error(err),
		)
		return nil, 0, domain.ErrDatabase
	}

	return results, total, nil
}
