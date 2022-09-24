package services

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"image"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Ressley/hacknu/internal/app/apiserver/helpers"
	"github.com/Ressley/hacknu/internal/app/apiserver/models"
	"github.com/devedge/imagehash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

var filesCollection *mongo.Collection = client.Database(helpers.PHOTO_DB).Collection(helpers.FS_FILES)
var chunksCollection *mongo.Collection = client.Database(helpers.PHOTO_DB).Collection(helpers.FS_CHUNKS)
var hashedPhotoCollection *mongo.Collection = client.Database(helpers.PHOTO_DB).Collection(helpers.HASHED_PHOTO)

func UploadFile(filename string, data []byte) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var result models.Photo

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatalln(err)
	}
	hexHash, _ := imagehash.Dhash(img, 8)
	hash := hex.EncodeToString(hexHash)

	err = hashedPhotoCollection.FindOne(ctx, bson.D{{Key: "hash", Value: hash}}).Decode(&result)
	if err == nil {
		return *result.Fileid, nil
	}

	bucket, err := gridfs.NewBucket(
		client.Database("HackNUPhoto"),
	)
	if err != nil {
		return "", err
	}

	uploadStream, err := bucket.OpenUploadStream(
		filename,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer uploadStream.Close()
	_, err = uploadStream.Write(data)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	str := fmt.Sprintf("%v", uploadStream.FileID)
	strs := strings.Split(str, `"`)

	photo := models.Photo{
		Fileid: &strs[1],
		Hash:   &hash,
	}
	photo.ID = primitive.NewObjectID()
	_, err = hashedPhotoCollection.InsertOne(ctx, photo)
	if err != nil {
		return "", err
	}
	return strs[1], nil
}

func DownloadFile(id string) ([]byte, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var results bson.D
	fileID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	_err := filesCollection.FindOne(ctx, bson.D{{Key: "_id", Value: fileID}}).Decode(&results)
	if _err != nil {
		return nil, err
	}

	bucket, _ := gridfs.NewBucket(
		client.Database("HackNUPhoto"),
	)
	var buf bytes.Buffer
	_, _err = bucket.DownloadToStream(fileID, &buf)
	if _err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
