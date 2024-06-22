# Conazon Users and Auth

This is the users and auth microservice for the ecommerce project. This application uses [Google for an Oauth provider](https://developers.google.com/identity/protocols/oauth2) and [JSON Web Tokens](https://jwt.io/) (JWT) for authentication.

## Quickstart

To test locally, setup a `.env` file in the root directory with the following variables:

`JWTSECRET` - Secret for JWT REQUIRED - Must be the same across microservices
`CLIENTID` - Client ID for Google Oauth REQUIRED
`CLIENTSECRET` - Secret for Google Oauth REQUIRED
`REDIRECTURL` - Redirect url for Google Oauth REQUIRED
`DATABASEURL` - Url to postgres database. REQUIRED
`SECURECOOKIE` - If true, enables secure on all cookies. Otherwise, default value of `false` is used
`PORT` - Port to run server on. Defaults to 8080

Datbase url should bne formatted this - 'host=postgres port=5432 user=postgres dbname=conazon sslmode=disable'

Then run:

`go build .`
`./conazon-users-and-auth`

Do not expose to the internet without setting the `SECURE` environment variable to `true` and setting up https.

## Endpoints (later will have swagger)

- /

GET - generic hello world. useless endpoint

- /users

DELETE - delete user

- /users/{id}

GET - get user info

- /auth/google/login

GET - starts oauth flow and redirects

- /auth/google/callback

GET - generates oauth token

- /logout

DELETE - removes cookie - later will add token invalidation
