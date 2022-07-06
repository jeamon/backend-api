package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jeamon/sample-rest-api/pkg/domain"
	mongodb "github.com/jeamon/sample-rest-api/pkg/infrastructure/mongo"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

type MongoScanInfosRepository struct {
	logger *zap.Logger
	mgo    *mongodb.Handler
	dbname string
}

// NewMongoScanInfosRepository provides an instance of MongoScanInfosRepository structure.
func NewMongoScanInfosRepository(logger *zap.Logger, h *mongodb.Handler, dbname string) *MongoScanInfosRepository {
	return &MongoScanInfosRepository{
		logger: logger,
		mgo:    h,
		dbname: dbname,
	}
}

func (repo *MongoScanInfosRepository) Save(ctx context.Context, s domain.ScanInfos) (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return s.ID, errors.Wrap(err, "could not save scan infos. unable to generate uuid")
	}
	s.ID = uid.String()
	collection := repo.mgo.Client.Database(repo.dbname).Collection("scan_infos")
	_, err = collection.InsertOne(ctx, s)
	if err != nil {
		return s.ID, errors.Wrapf(err, "could not save scan infos")
	}
	return s.ID, nil
}

func (repo *MongoScanInfosRepository) FindByID(ctx context.Context, id string) (domain.ScanInfos, error) {
	s := domain.ScanInfos{}
	collection := repo.mgo.Client.Database(repo.dbname).Collection("scan_infos")
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(s)
	return s, errors.Wrapf(err, "could not find scan infos with ID: %s", id)
}

func (repo *MongoScanInfosRepository) FindAll(ctx context.Context) ([]domain.ScanInfos, error) {
	var res []domain.ScanInfos
	collection := repo.mgo.Client.Database(repo.dbname).Collection("scan_infos")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return res, errors.Wrapf(err, "could not find all scan infos")
	}

	defer cursor.Close(ctx)
	for cursor.Next(context.TODO()) {
		var s domain.ScanInfos
		if err := cursor.Decode(&s); err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

func (repo *MongoScanInfosRepository) UpdateByID(ctx context.Context, id string, s domain.ScanInfos) error {
	collection := repo.mgo.Client.Database(repo.dbname).Collection("scan_infos")
	update := bson.M{"$set": bson.M{"username": s.Username, "company_id": s.CompanyID}}
	_, err := collection.UpdateOne(
		ctx, bson.M{"_id": s.ID},
		update,
	)
	return errors.Wrapf(err, "cannot update scan infos with ID: %s", id)
}

func (repo *MongoScanInfosRepository) DeleteByID(ctx context.Context, id string) error {
	collection := repo.mgo.Client.Database(repo.dbname).Collection("scan_infos")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return errors.Wrapf(err, "could not delete scan infos with ID: %s", id)
}
