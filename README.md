# README
Purpose of this script is to download images from S3 specific bucket, compress them and upload them back to S3 with
suffix _new

## Install glide
https://glide.sh/

## Install dependencies
```bash
glide install
```

## Build
```bash
go build
```

## Run
```bash
CJPEG_PATH=/usr/local/Cellar/mozjpeg/3.1_1/bin/cjpeg BUCKET_NAME=imagecompress ./s3_images_compressor
```
