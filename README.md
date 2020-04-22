# customerService
    This project is used to maintain customer details and developed using the following technologies
    - GoLang
    - Gin
    - MongoDB
    - Swagger
    - gin-swagger
    
# Steps
1. Open the command prompt and navigate to root folder
2. Download swag by using:
    ```console
    $ go get -u github.com/swaggo/swag/cmd/swag
    ```
3. Run swag init in the project's root folder which contains the main.go file
    ```console
    $ swag init
   ```
4. Run your app, and browse to http://localhost:8080/swagger/index.html.

Note: I made a simple authenticate as token. Please use Authorization header as "123-456-789"