# Go Stripe

Sample Golang application with stripe payment

## How to start

### First Scaffolding

1. Start backend `make start_back` (to stop run `make stop_back`) 
2. Start frontend `air`
3. you can run both using single command `make start` or `make stop` for stopping
4. open browser `http://localhost:4001/api/payment-intent` to test backend
5. open browser `http://localhost:4000/virtual-terminal` to test frontend
6. do some test (see https://stripe.com/docs/testing)
7. check your stripe dashboard (https://dashboard.stripe.com/test/payments)
8. for login testing use credentials `admin@example.com` and `password` 

### Database setup

1. build mysqldb with docker `docker-compose -p mysql -f ./docker-compose.yml up -d`
2. connect to db using root user
3. run query `GRANT ALL ON widgets.* TO 'ananto'@'%' IDENTIFIED BY 'secret'` or `CREATE USER 'ananto'@'%' IDENTIFIED BY 'secret'; && GRANT ALL PRIVILEGES ON widgets.* TO 'ananto'@'%';` for newest version


### Migrate Database
- Opsi 1
1. seed database using `go run ./cmd/seed/*.go`
2. test widget api `http://localhost:4001/api/widget/1`
- Opsi 2
1. run `soda migrate`


## Note
1. if database failure to save transaction, we should cancel the order
2. if one of database saving process fail, we should rollback the transaction
3. parse int amout: input 12.345 shall output 12.34

