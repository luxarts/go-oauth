# go-oauth

## How to use the login
1. Set the following envs
    ```
    OPENID_CONFIG_URL=<Endpoint to get the openid configuration>
    CLIENT_ID=<client ID in your auth server>
    CLIENT_SECRET=<client secret>
    REDIRECT_URI=<url where the service is up>/callback
    JWKS_URL=<Endpoint to get jwks>
    ```
2. Get the login URL
    ```
    curl --location 'localhost:8080/login'
    ```
3. After login completes it will redirect to callback and will return the token

## Validate token
1. Import the package
   ```go
   import "github.com/luxarts/go-auth/jwt"
   ```
2. Use it as follows
   ```go
   isValid := jwt.NewValidator().IsValid(token)
   ```