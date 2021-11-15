# goMasterClass
SimpleBank project from YouTube

// creating migration files
migrate create -ext sql -dir db/migration -seq init_schema

steps
1. run postgres in docker
2. create db model using dbdiagram, export as postgresql command
3. install golang migrate
4. create migration files using command above
5. init sqlc
6. add query, schema, and destination paths in the sqlc, add necessary tags by checking the doc
7. start docker postgres and createdb using docker exec shell
8. add queries by checking sqlc documentation and sqlc generate. Check for errors
9. check sqlc output inside path given under sqlc yaml
10. write functional test using custom made helper package utils
