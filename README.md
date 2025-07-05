# Go Stripe

Sample Golang application with stripe payment

## How to start

### First Scaffolding

1. setup database via docker `docker-compose -p mysql -f ./docker-compose.yml up -d`
2. migrate database using [soda](https://gobuffalo.io/pt/documentation/database/soda/), run `soda migrate`
3. Start backend `make start_back`, to stop run `make stop_back`
4. Start frontend `air`
5. alternatively, you can run both using single command `make start` or `make stop` for stopping
6. open browser `http://localhost:4001/api/payment-intent` to test backend
7. open browser `http://localhost:4000/virtual-terminal` to test frontend
8. do some test (see <https://stripe.com/docs/testing>)
9. check your stripe dashboard (<https://dashboard.stripe.com/test/payments>)
10. for login testing use credentials `admin@example.com` and `password`

## Note

### TODO

1. if database failure to save transaction, we should cancel the order
2. if one of database saving process fail, we should rollback the transaction
3. parse int amout: input 12.345 shall output 12.34
4. input email should be lowercase
5. Error 

```js
Uncaught (in promise) TypeError: Cannot read properties of undefined (reading 'addEventListener')
    at paginator (all-sales:177:28)
    at all-sales:241:17
```

### Ref

1. <https://developer.mozilla.org/en-US/docs/Web/API/Window/localStorage>
