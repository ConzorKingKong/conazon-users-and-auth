# Ecommerce Users and Auth

This is the auth microservice for the ecommerce project. This application uses [Google for an Oauth provider](https://developers.google.com/identity/protocols/oauth2) and [JSON Web Tokens](https://jwt.io/) (JWT) for authentication.

## Quickstart

To test locally, setup a `.env` file in the root directory with the following variables:

`JWTSECRET` - Secret for JWT REQUIRED
`CLIENTID` - Client ID for Google Oauth REQUIRED
`CLIENTSECRET` - Secret for Google Oauth REQUIRED
`REDIRECTURL` - Redirect url for Google Oauth REQUIRED
`SECURECOOKIE` - If true, enables secure on all cookies. Otherwise, default value of `false` is used

Then run:

`go build .`
`./auth`

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
