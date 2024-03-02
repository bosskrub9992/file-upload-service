# file-upload-service

#### getting started

1. Copy the attached key.json in the sent email, and paste it in the folder ./file-upload-service
   
2. Run
```sh
podman run --rm -it -p 3300:80 -p 2525:25 rnwood/smtp4dev
```

3. Open browser and navigate to localhost:3300 to open testing email server
```
localhost:3300
```

4. Run
```sh
go run main.go
```
It will start the REST API server on port 8080.

5. Call API by run the below curl command, Please give your own full file path to the '-F' flag. The file can be in jpeg, png or heic format. Or you can import the curl command into to Postman to be more convenience to test
```sh
curl -L 'http://localhost:8080/api/v1/files' \
-F 'file=@"/Users/s.thunmanuthum/Downloads/please.jpeg"'
```
