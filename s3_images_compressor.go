package main

import (
    "log"
    "os"
    "os/exec"
    "strings"

    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)


func main() {
    aws_session := session.New(&aws.Config{Region: aws.String("eu-central-1")})
    uploader := s3manager.NewUploader(aws_session)
    downloader := s3manager.NewDownloader(aws_session)
    bucket := os.Getenv("BUCKET_NAME")
    svc := s3.New(aws_session)

    result, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket)})
    if err != nil {
      log.Println("Failed to load objects", err)
      return
    }

    for _, file := range result.Contents {
        key := aws.StringValue(file.Key)
        size := aws.Int64Value(file.Size)

        if size == 0 {
            log.Printf("Create directory %s", key)
            os.MkdirAll(key, 0700)
        } else {
            log.Printf("%s", key)
            new_key := strings.Replace(key, ".JPG", "_new.JPG", 1)
            downloadFile(downloader, bucket, key)
            compressFile(key, new_key)
            uploadFile(uploader, bucket, key, new_key)
        }
    }
}


func downloadFile(downloader *s3manager.Downloader, bucket string, key string){
    file, err := os.Create(key)
    if err != nil {
        log.Println("Failed to create file", err)
    }
    defer file.Close()

    numBytes, err := downloader.Download(file,
        &s3.GetObjectInput{
            Bucket: aws.String(bucket),
            Key:    aws.String(key),
        })
    if err != nil {
        log.Println("Failed to download file", err)
        return
    }

    log.Println("Downloaded file", file.Name(), numBytes, "bytes")
}

func compressFile(key string, new_key string){
    outfile, err := os.Create(new_key)
    if err != nil {
      log.Println("Failed to create file", err)
      return
    }

    defer outfile.Close()

    cmd := exec.Command(os.Getenv("CJPEG_PATH"), key)
    cmd.Stdout = outfile
    compile_err := cmd.Run()

    if compile_err != nil {
      log.Println("Failed to compress file", new_key, compile_err)
      return
    }

    log.Println("Compressed file", key)
}


func uploadFile(uploader *s3manager.Uploader, bucket string, key string, new_key string) {
    file, err := os.Open(new_key)

    if err != nil {
      log.Println("Failed to open file", err)
      return
    }

    result, err := uploader.Upload(&s3manager.UploadInput{
      Body:   file,
      Bucket: aws.String(bucket),
      Key:    aws.String(new_key),
    })

    if err != nil {
      log.Println("Failed to upload", err)
      return
    }

    log.Println("Successfully uploaded to", result.Location)
}
