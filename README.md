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

### Database setup

8. build mysqldb with docker `docker-compose -p mysql -f ./docker-compose.yml up -d`
9. connect to db using root user
9. run query `GRANT ALL ON widgets.* TO 'ananto'@'%' IDENTIFIED BY 'secret'` or `CREATE USER 'ananto'@'%' IDENTIFIED BY 'secret'; && GRANT ALL PRIVILEGES ON widgets.* TO 'ananto'@'%';` for newest version

