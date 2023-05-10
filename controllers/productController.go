package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"task1-ch2-api-product/database"
	"task1-ch2-api-product/models"
)

var productCollection *mongo.Collection = database.OpenCollection(database.Client, "product")
var productValidate = validator.New()

// CreateProduct is the api used to create a single product
func CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var product models.Product
		defer cancel()
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := productValidate.Struct(product)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		product.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		product.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		product.ID = primitive.NewObjectID()
		product.Product_id = product.ID.Hex()

		// Upload file to MinIO
		file, err := c.FormFile("image")
		if err == nil {
			// Generate random file name
			randName := uuid.New().String() + path.Ext(file.Filename)
			// Connect to MinIO server
			endpoint := os.Getenv("MINIO_ENDPOINT")
			accessKey := os.Getenv("MINIO_ACCESS_KEY")
			secretKey := os.Getenv("MINIO_SECRET_KEY")
			useSSL := true // Set to true if using HTTPS
			minioClient, err := minio.New(endpoint, &minio.Options{
				Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
				Secure: useSSL,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Open the file for reading
			fileReader, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			defer fileReader.Close()

			// Read the contents of the file into a byte slice
			fileBytes, err := ioutil.ReadAll(fileReader)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Upload file to MinIO
			_, err = minioClient.PutObject(ctx, "product-images", randName, bytes.NewReader(fileBytes), int64(len(fileBytes)), minio.PutObjectOptions{})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Set image URL in product
			product.ImageUrl = "https://" + endpoint + "/product-images/" + randName
		}

		resultInsertionNumber, insertErr := productCollection.InsertOne(ctx, product)
		if insertErr != nil {
			msg := fmt.Sprintf("Product item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

// GetProduct is the api used to get a single product
func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var product models.Product
		product_id := c.Param("id")
		defer cancel()

		err := productCollection.FindOne(ctx, bson.M{"product_id": product_id}).Decode(&product)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "product not found"})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

// GetProducts is the api used to get all products
func GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Define filter and options for MongoDB find operation
		filter := bson.M{}
		options := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

		// Execute MongoDB find operation
		cursor, err := productCollection.Find(ctx, filter, options)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
			return
		}
		defer cursor.Close(ctx)

		// Iterate over the cursor and append products to the results slice
		var results []models.Product
		for cursor.Next(ctx) {
			var product models.Product
			if err := cursor.Decode(&product); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode product"})
				return
			}
			results = append(results, product)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

// UpdateProduct is the api used to update a single product
func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Get product ID from URL parameter
		productID := c.Param("id")

		// Validate product ID
		id, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
			return
		}

		// Get request body and bind it to Product model
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate request body
		validationErr := validate.Struct(product)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Create filter and update values for MongoDB update operation
		filter := bson.M{"_id": id}
		update := bson.M{"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"category":    product.Category,
			"price":       product.Price,
			"rate":        product.Rate,
			"location":    product.Location,
			"quantity":    product.Quantity,
			"store":       product.Store,
			"updated_at":  time.Now().UTC(),
		}}

		// Upload new image to MinIO and update image URL in product
		file, err := c.FormFile("image")
		if err == nil {
			// Generate random file name
			randName := uuid.New().String() + path.Ext(file.Filename)

			// Connect to MinIO server
			endpoint := os.Getenv("MINIO_ENDPOINT")
			accessKey := os.Getenv("MINIO_ACCESS_KEY")
			secretKey := os.Getenv("MINIO_SECRET_KEY")
			useSSL := true // Set to true if using HTTPS
			minioClient, err := minio.New(endpoint, &minio.Options{
				Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
				Secure: useSSL,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Open the file for reading
			fileReader, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			defer fileReader.Close()

			// Read the contents of the file into a byte slice
			fileBytes, err := ioutil.ReadAll(fileReader)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Upload file to MinIO
			_, err = minioClient.PutObject(ctx, "product-images", randName, bytes.NewReader(fileBytes), int64(len(fileBytes)), minio.PutObjectOptions{})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Set image URL in product
			update["$set"].(bson.M)["image"] = "https://" + endpoint + "/product-images/" + randName
		}

		// Perform MongoDB update operation
		result, err := productCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
			return
		}

		// Check if product was updated
		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "product updated successfully"})
	}
}

// DeleteProduct is the api used to delete a single product
func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// get product id from url parameter
		id := c.Param("id")

		// convert product id to ObjectID
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
			return
		}

		// delete product from database
		result, err := productCollection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
	}
}

// SearchProduct is the api used to search for products by name
func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Get query parameter "name" from URL
		name := c.Query("name")

		// Define filter and options for MongoDB find operation
		filter := bson.M{"name": primitive.Regex{Pattern: name, Options: "i"}}
		options := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

		// Execute MongoDB find operation
		cursor, err := productCollection.Find(ctx, filter, options)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search for products"})
			return
		}
		defer cursor.Close(ctx)

		// Iterate over the cursor and append products to the results slice
		var results []models.Product
		for cursor.Next(ctx) {
			var product models.Product
			if err := cursor.Decode(&product); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode product"})
				return
			}
			results = append(results, product)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search for products"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}
