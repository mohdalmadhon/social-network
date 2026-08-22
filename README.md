# SOCIAL NETWORK

REQUIREMENTS
===========
    * golang should be downloaded in the running server
    * docker should be avilable and running
    * golang-migrate should be installed with CGO enabled: "go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    * bcrypt package "go install golang.org/x/crypto/bcrypt"
        - NOTE: just run the project and go packages they will be downloaded automatically
    * if you are using VC CODE recommended to download vue official extention

RUN
============
 - make sure to do "npm install" before lunching to download all the required dependencies into node_modules
 - apply database migrations from the project root before starting the backend
 - start the backend from the `backend` folder with `go run ./cmd/main.go`
 - start the frontend from the `frontend` folder with `npm run dev`
 - open `http://localhost:5173/login`

POSTS API
============
`GET /api/posts` returns the signed-in user's filtered personal feed.

`POST /api/posts` creates a post. Both routes require the existing `token` cookie.

```json
{
  "content": "Hello, Orbit!",
  "privacy": "public",
  "selectedFollowerIds": []
}
```

Privacy can be `public`, `followers`, or `selected`. Selected posts require one or more IDs belonging to users who follow the author.
