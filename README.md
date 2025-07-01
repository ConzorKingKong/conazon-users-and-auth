# Conazon Users and Auth

This is the users and auth microservice for the ecommerce project. This application uses [Google for an Oauth provider](https://developers.google.com/identity/protocols/oauth2) and [JSON Web Tokens](https://jwt.io/) (JWT) for authentication.

Coupled at first for simplicity, this will be split into 2 different services soon and will add m2m to auth service

## Quickstart

To test locally, setup a `.env` file in the root directory with the following variables:

```
JWTSECRET - Secret for JWT REQUIRED - Must be the same across microservices
CLIENTID - Client ID for Google Oauth REQUIRED
CLIENTSECRET - Secret for Google Oauth REQUIRED
REDIRECTURL - Redirect url for Google Oauth REQUIRED
DATABASEURL - Url to postgres database. REQUIRED
PROTOCOL - http or https in prod
HOSTNAME - hostname you're project will run on
SECURECOOKIE - If true, enables secure on all cookies (only use cookie on https). Otherwise, default value of `FALSE` is used. MUST BE TRUE IN PROD
PORT - Port to run server on. Defaults to 8080
```

Datbase url should be formatted this - `'host=postgres port=5432 user=postgres dbname=conazon sslmode=disable'`

Then run:

`docker-compose up`

Do not expose to the internet without setting the `SECURE` environment variable to `true` and setting up https.

## Endpoints (later will have swagger)

- /

GET - Catch all 404

### USERS

- /users

GET - Lists users
POST - create user (to be implemented. currently created by /auth/google/callback)

- /users/{id}

GET - get public user info
POST - Update user info PROTECTED
DELETE - delete user PROTECTED

- /me

GET - Pulls non-public user info PROTECTED

### AUTH

- /auth/google/login

GET - starts oauth flow and redirects to google

- /auth/google/callback

GET - generates oauth token, creates user, start cookie session, and redirects to homepage 

- /verify

GET - Verifies if user has a session when front-end first loads PROTECTED

- /logout

DELETE - removes cookie - later will add token invalidation PROTECTED

### HEALTH

- /healthz

GET - returns 200

