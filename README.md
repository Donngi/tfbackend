# tfbackend
`tfbackend` enable you to create terraform backend at a blistering speed🚀🚀

## Usage
### AWS
```
$ tfbackend aws --s3 YOUR_BUCKET_NAME --dynamodb YOUR_TABLE_NAME
```

### Other
TBD

## Why is tfbackend worth?
It is the best practice of terraform to create backend by means other than using terraform project itself.

Many of terraform users make it by
- console of cloud
- cli command
, however to apply appropriate configuration to backend is a bit complicated.

How we should do? `tfbackend` solve this problem!🚀

By only using a single command like `'tfbackend aws --s3 backend-bucket'`, you can create resources which is optimized for terraform backend.

For example, in AWS, you can make S3 bucket with
- enabled public access block
- enabled transparent encryption
- enabled versioning

## Installation
### Homebrew
TBD

### The others
```
go 1.16~
$ go install github.com/Jimon-s/tfbackend@latest
```
