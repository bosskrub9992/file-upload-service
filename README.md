# file-upload-service

file-upload-service is a RESTful API server. It has an endpoint that can upload files. All of the uploaded files will be stored in Google Cloud Storage (GCS) in my account. You can use this email account, which has only view permission, to view the files that is already uploaded.

email: publicemailforjob@gmail.com
pw: Gmail-temp-password-2024

#### Prerequisite

1. Go 1.22
2. Docker or Podman

#### Getting started

1. clone project by run
```sh
git clone https://github.com/bosskrub9992/file-upload-service.git
```
2. From the sent email, Download the attached file named 'key.json' then paste it at the root of the project folder ./file-upload-service
3. Run this command
```sh
docker run --rm -it -p 3300:80 -p 2525:25 rnwood/smtp4dev
```
4. Open browser and navigate to localhost:3300 to open testing email server
```
localhost:3300
```
5. Run the server, It will start a REST API server on port 8080.
```sh
go run main.go
```
6. Call API by run the below curl command, Please give your own full file path to the '-F' flag. The file can be in jpeg, png or heic format.
```sh
curl -L 'http://localhost:8080/api/v1/files' \
-F 'file=@"/Users/s.thunmanuthum/Downloads/please.jpeg"'
```

#### Assignment Concerns that needed to be addressed

##### Security
1. Set the maximum size of the upload file
2. Whitelist the allowed content types and file extensions of the upload file

##### Scalability
1. The service is developed with Go programming language, a modern language that is suited for deploying on cloud.
2. The service is stateless.
3. The uploaded files are stored in Google Cloud Storage which is highly scalable.
4. The email notification is designed to be async with the upload function.

Usually, you can do vertical scaling. But with the above reasons, horizontal scaling is also possible.

##### Testability
1. There is unit test implemented at the business logic layer.
2. You can do end-to-end testing like in the 'Getting started' part.

##### API standard (RESTful)
1. Resource based naming endpoint
