package mongo

import (
	"companies/internal/domain/interfaces"
	"companies/internal/entity"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	db         = "companies"
	companyCol = "company"
)

type CompaniesRepoMongo struct {
	Mongo *mongo.Client
}

func NewCompaniesRepo(mongo *mongo.Client) *CompaniesRepoMongo {
	if mongo == nil {
		return nil
	}
	return &CompaniesRepoMongo{
		Mongo: mongo,
	}
}

func (r *CompaniesRepoMongo) Create(ctx context.Context, createParams entity.CreateCompany) error {
	set := bson.D{
		{"uid", createParams.UId},
		{"name", createParams.Name},
		{"type", createParams.Type},
		{"employees", createParams.Employees},
		{"registered", createParams.Registered}}

	if createParams.Description != nil {
		set = append(set, bson.E{Key: "description", Value: *createParams.Description})
	}
	doc := bson.D{
		{"$set", set},
		{"$setOnInsert", bson.D{
			{"created_at", primitive.NewDateTimeFromTime(time.Now())},
		}},
	}

	if r == nil { // might be nil, as initially the interface method is called
		return errors.New("companies repo not initialized")
	}
	if r.Mongo == nil {
		return errors.New("companies mongo client not initialized")
	}
	collection := r.Mongo.Database(db).Collection(companyCol)
	if collection == nil {
		return errors.New("company collection not initialized")
	}

	filter := bson.D{{"uid", createParams.UId}}
	res, err := collection.UpdateOne(ctx, filter, doc, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 && res.UpsertedCount == 0 {
		return interfaces.ErrStorageNonRetryable // not modified and not inserted somehow
	}

	return nil
}

func (r *CompaniesRepoMongo) FetchByUid(ctx context.Context, uid string) (*entity.Company, error) {
	filter := bson.M{"uid": uid}

	var company entity.Company
	if r == nil { // might be nil, as initially the interface method is called
		return nil, errors.New("companies repo not initialized")
	}
	if r.Mongo == nil {
		return nil, errors.New("companies mongo client not initialized")
	}
	collection := r.Mongo.Database(db).Collection(companyCol)
	if collection == nil {
		return nil, errors.New("company collection not initialized")
	}
	err := collection.FindOne(ctx, filter).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if the company is not found
		}
		return nil, err
	}

	return &company, nil
}

func (r *CompaniesRepoMongo) UpdateByUId(ctx context.Context, uid string, updateParams entity.UpdateCompany) error {
	filter := bson.M{"uid": uid}

	set := bson.D{}

	if updateParams.Description != nil {
		set = append(set, bson.E{Key: "description", Value: *updateParams.Description})
	}
	if updateParams.Name != nil {
		set = append(set, bson.E{Key: "name", Value: *updateParams.Name})
	}
	if updateParams.Employees != nil {
		set = append(set, bson.E{Key: "employees", Value: *updateParams.Employees})
	}
	if updateParams.Type != nil {
		set = append(set, bson.E{Key: "type", Value: *updateParams.Type})
	}
	if updateParams.Registered != nil {
		set = append(set, bson.E{Key: "registered", Value: *updateParams.Registered})
	}
	doc := bson.D{
		{"$set", set},
		{"$currentDate", bson.D{
			{"updated_at", true},
		}},
	}
	if r == nil { // might be nil, as initially the interface method is called
		return errors.New("companies repo not initialized")
	}
	if r.Mongo == nil {
		return errors.New("companies mongo client not initialized")
	}
	collection := r.Mongo.Database(db).Collection(companyCol)
	if collection == nil {
		return errors.New("company collection not initialized")
	}

	res, err := collection.UpdateOne(ctx, filter, doc)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 && res.ModifiedCount == 0 {
		return interfaces.ErrStorageNonRetryable // not modified and not inserted somehow
	}

	return nil
}

func (r *CompaniesRepoMongo) DeleteByUId(ctx context.Context, uid string) error {
	filter := bson.M{"uid": uid}
	if r == nil { // might be nil, as initially the interface method is called
		return errors.New("companies repo not initialized")
	}
	if r.Mongo == nil {
		return errors.New("companies mongo client not initialized")
	}
	collection := r.Mongo.Database(db).Collection(companyCol)
	if collection == nil {
		return errors.New("company collection not initialized")
	}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return interfaces.ErrStorageNonRetryable // not modified and not inserted somehow
	}

	return nil
}
