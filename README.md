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
11. write code for the transaction, by decalring new store variable with dbtx interface
12. define exectx function for executing any transactions, pay attention to begin, commit and rollback
13. define TransferTx in the order of a transaction operation (deduction, addition, writing entries, transfer operation)
14. understand why ACID propoerties, how to achieve the db isolation and understand its importance (read uncommit, read commit, repeatable read, serializable) and different read phenomenas and how to avoid them.
15. Implement test-driven-development(TDD) for account transaction in the for TransferTx
16. Change sqlc queries for updateaccount with locking (for mysql it is locking and for pistgresql it is dependency management)
17. write an extra function as add account balance for easy balance update in account, use this fuction refactor TransferTx
18. use gin for handling http requests
19. use viper to load config

