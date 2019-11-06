package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/minio/minio-go"
)

var listen string
var minioClient *minio.Client

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	listen = os.Getenv("LISTEN")
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioID := os.Getenv("MINIO_ACCESS_KEY")
	minioKey := os.Getenv("MINIO_SECRET_KEY")

	// Minio client
	minioClient, err = minio.New(minioEndpoint, minioID, minioKey, false)
	if err != nil {
		log.Fatal("error loading minio")
	}

	// Create bucket for parsers if it doesn't exist
	err = minioClient.MakeBucket("parser", "us-east-1")
	if err != nil {
		exists, err := minioClient.BucketExists("parser")
		if err == nil && exists {
			log.Printf("bucket parser already exists")
		} else {
			log.Printf("%s", err)
			log.Fatal("error creating bucket")
		}
	} else {
		log.Printf("created bucket parser")
	}

	expiry := time.Second * 120

	// Routes
	router := httprouter.New()
	// Route classes
	m := ModelHandler{minioClient, expiry}
	p := ParserHandler{minioClient, expiry}
	md := ModelDataHandler{minioClient, expiry}
	b := BatchHandler{minioClient, expiry}
	s := ValuesHandler{make(map[string]Session)}
	// Return minio presigned URLs
	// Model routes
	router.GET("/models", m.GetModels)
	router.GET("/model/:id", m.GetModelById)
	router.PUT("/model", m.UploadModel)
	// Parser routes
	router.GET("/parsers", p.GetParsers)
	router.GET("/parser/:id", p.GetParserById)
	router.PUT("/parser", p.UploadParser)
	// Model Data routes
	router.GET("/model/:model/data", md.GetModelData)
	router.GET("/model/:model/data/:id", md.GetModelDataById)
	router.PUT("/model/:model/data", md.UploadModelData)
	router.GET("/model/:model/labels", md.GetModelLabels)
	router.GET("/model/:model/labels/:id", md.GetModelLabelsById)
	router.PUT("/model/:model/labels", md.UploadModelLabels)
	// Batch routes
	router.GET("/model/:model/batch", b.GetBatch)
	router.GET("/model/:model/batch/random", b.GetBatchRand)
	router.GET("/model/:model/batch/:id/data", b.GetBatchData)
	router.GET("/model/:model/batch/:id/labels", b.GetBatchLabels)
	router.POST("/model/:model/data/:id/batch", b.BatchData)
	// Weights/Session routes
	router.GET("/session/:id/loss", s.GetLoss)
	router.POST("/session/:id/loss", s.PostLoss)
	router.POST("/session/:id/weights", s.PostWeights)
	router.POST("/session/:id", s.NewSession)

	// Start server
	log.Printf("starting server on %s", listen)
	log.Fatal(http.ListenAndServe(listen, router))
}
