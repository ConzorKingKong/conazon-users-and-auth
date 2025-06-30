# Conazon Users and Auth

This is the users and auth microservice for the ecommerce project. This application uses [Google for an Oauth provider](https://developers.google.com/identity/protocols/oauth2) and [JSON Web Tokens](https://jwt.io/) (JWT) for authentication.

## Quickstart

To test locally, setup a `.env` file in the root directory with the following variables:

```
JWTSECRET - Secret for JWT REQUIRED - Must be the same across microservices
CLIENTID - Client ID for Google Oauth REQUIRED
CLIENTSECRET - Secret for Google Oauth REQUIRED
REDIRECTURL - Redirect url for Google Oauth REQUIRED
DATABASEURL - Url to postgres database. REQUIRED
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

- /users

DELETE - delete user

- /users/{id}

GET - get user info

- /auth/google/login

GET - starts oauth flow and redirects to google

- /auth/google/callback

GET - generates oauth token, start cookie session, and redirects to homepage 

- /logout

DELETE - removes cookie - later will add token invalidation
