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
	router.HEAD("/model/:id", m.HeadModelById)
	router.PUT("/model/:id", m.UploadModel)
	// Parser routes
	router.GET("/parsers", p.GetParsers)
	router.GET("/parser/:id", p.GetParserById)
	router.HEAD("/parser/:id", p.HeadParserById)
	router.PUT("/parser", p.UploadParser)
	// Model Data routes
	router.GET("/model/:id/data", md.GetModelData)
	router.GET("/model/:id/data/:dataid", md.GetModelDataById)
	router.HEAD("/model/:id/data/:dataid", md.HeadModelDataById)
	router.PUT("/model/:id/data/:dataid", md.UploadModelData)
	router.GET("/model/:id/labels", md.GetModelLabels)
	router.GET("/model/:id/labels/:labelsid", md.GetModelLabelsById)
	router.HEAD("/model/:id/labels/:labelsid", md.HeadModelLabelsById)
	router.PUT("/model/:id/labels/:labelsid", md.UploadModelLabels)
	// Batch routes
	router.GET("/model/:id/batches", b.GetBatch)
	router.GET("/model/:id/batch", b.GetBatchRand)
	router.GET("/model/:id/batch/:batchid/data", b.GetBatchData)
	router.HEAD("/model/:id/batch/:batchid/data", b.HeadBatchData)
	router.GET("/model/:id/batch/:batchid/labels", b.GetBatchLabels)
	router.HEAD("/model/:id/batch/:batchid/labels", b.HeadBatchLabels)
	router.POST("/model/:id/data/:batchid/batch", b.BatchData)
	// Weights/Session routes
	router.GET("/session/:id/loss", s.GetLoss)
	router.POST("/session/:id/loss", s.PostLoss)
	router.POST("/session/:id/weights", s.PostWeights)
	router.POST("/session/:id", s.NewSession)

	// Start server
	log.Printf("starting server on %s", listen)
	log.Fatal(http.ListenAndServe(listen, router))
}
