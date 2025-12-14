package repositories

import (
    "context"
    "errors"
    "time"
    
    "ecommerce-backend/internal/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
    collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
    return &UserRepository{
        collection: db.Collection("users"),
    }
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
    existingUser, _ := r.FindByEmail(ctx, user.Email)
    if existingUser != nil {
        return nil, errors.New("user with this email already exists")
    }
    
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    
    result, err := r.collection.InsertOne(ctx, user)
    if err != nil {
        return nil, err
    }
    
    user.ID = result.InsertedID.(primitive.ObjectID)
    return user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    filter := bson.M{"email": email}
    
    err := r.collection.FindOne(ctx, filter).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil
        }
        return nil, err
    }
    
    return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
    var user models.User
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }
    
    filter := bson.M{"_id": objectID}
    err = r.collection.FindOne(ctx, filter).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil
        }
        return nil, err
    }
    
    return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id string, updateData *models.UpdateUserRequest) (*models.User, error) {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }
    
    update := bson.M{
        "$set": bson.M{
            "name":       updateData.Name,
            "email":      updateData.Email,
            "updated_at": time.Now(),
        },
    }
    
    _, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
    if err != nil {
        return nil, err
    }
    
    return r.FindByID(ctx, id)
}
