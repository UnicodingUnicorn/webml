package main

import (
  "io"
  "log"
  "net/http"
  "os"
  "strconv"
  "time"

  "github.com/joho/godotenv"
  "github.com/julienschmidt/httprouter"
  "github.com/yuin/gopher-lua"
  "github.com/minio/minio-go"
)

var listen string
var minioClient *minio.Client

func main() {
  // Load .env
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  listen = os.Getenv("LISTEN")
  minioEndpoint := os.Getenv("MINIO_ENDPOINT")
  minioID := os.Getenv("MINIO_ACCESS_KEY")
  minioKey := os.Getenv("MINIO_SECRET_KEY")

  // Minio client
  minioClient, err = minio.New(minioEndpoint, minioID, minioKey, false)
  if err != nil {
    log.Fatal("Error loading minio")
  }

  // Create bucket if it doesn't exist
  err = minioClient.MakeBucket("parser", "us-east-1")
  if err != nil {
    exists, err := minioClient.BucketExists("parser")
    if err == nil && exists {
      log.Printf("Bucket %s already exists", "parser")
    } else {
      log.Printf("%s", err)
      log.Fatal("Error creating bucket")
    }
  } else {
    log.Printf("Created bucket %s", "parser")
  }

  // Routes
  router := httprouter.New()
  // Return minio presigned URLs
  router.POST("/model", UploadModel)
  router.POST("/data", UploadData)
  router.POST("/data_parser", UploadDataParser)

  router.POST("/parse", TestParse)

  // Start server
  log.Printf("starting server on %s", listen)
  log.Fatal(http.ListenAndServe(listen, router))
}

func UploadModel(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  bucketName := r.FormValue("id")

  // Create bucket if it doesn't exist
  err := minioClient.MakeBucket(bucketName, "us-east-1")
  if err != nil {
    exists, err := minioClient.BucketExists(bucketName)
    if !(err == nil && exists) {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
  		return
    }
  }

  expiry := time.Second * 120;
  presignedURL, err := minioClient.PresignedPutObject(bucketName, "model", expiry)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
  }

  w.Write([]byte(presignedURL.String()))
}

func UploadData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  bucketName := r.FormValue("id")
  hash := r.FormValue("hash")

  exists, err := minioClient.BucketExists(bucketName)
  if !(err == nil && exists) {
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
    return
  }

  expiry := time.Second * 120;
  presignedURL, err := minioClient.PresignedPutObject(bucketName, hash, expiry)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
  }

  w.Write([]byte(presignedURL.String()))
}

func UploadDataParser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  hash := r.FormValue("hash")

  expiry := time.Second * 120;
  presignedURL, err := minioClient.PresignedPutObject("parser", hash, expiry)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
  }

  w.Write([]byte(presignedURL.String()))
}

func TestParse(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  file, _, err := r.FormFile("file")
  if err != nil {
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
  }
  defer file.Close()

  L := lua.NewState()
  defer L.Close()
  err = L.DoFile("../mnist_data_parser.lua")
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
  }

  num := 0

  buf := make([]byte, 512)
  for {
    n, err := file.Read(buf)
    if err == io.EOF {
      break
    } else if err != nil {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
  		return
    }

    err = L.CallByParam(lua.P {
      Fn: L.GetGlobal("parse"),
      NRet: 1,
      Protect: true,
      }, lua.LString(buf), lua.LNumber(n))
    if err != nil {
      log.Printf("%s", err)
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
  		return
    }

    lv := L.Get(-1)
    L.Pop(1)
    if table, ok := lv.(*lua.LTable); ok {
      num += table.Len()
    }
  }

  w.Write([]byte(strconv.Itoa(num)))
}