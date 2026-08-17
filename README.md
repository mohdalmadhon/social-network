# SOCIAL NETWORK

REQUARMENTS
===========
    * golang should be downloaded in the running server
    * docker should be avilable and running
    * golang-magirate package should be installed by this command: "go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"     
    * bcrypt package "go install golang.org/x/crypto/bcrypt"
        - NOTE: just run the project and go packages they will be downloaded automatically
    * if you are using VC CODE recommended to download vue official extention

RUN
============
 - make sure to do "npm install" before lunching to download all the required dependencies into node_modules
 - to run the Application use "go run ./backend/cmd/main.go"