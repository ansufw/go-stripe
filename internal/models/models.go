package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// DBModel is the type for database connnection values
type DBModel struct {
	DB *sql.DB
}

// Models is the wrapper for all models
type Models struct {
	DB DBModel
}

// New returns a model type with database connection pool
func New(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// Widget is the type for all widgets
type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	IsRecurring    bool      `json:"is_recurring"`
	PlanID         string    `json:"plan_id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

// Order is the type for all orders
type Order struct {
	ID            int         `json:"id"`
	WidgetID      int         `json:"widget_id"`
	TransactionID int         `json:"transaction_id"`
	CustomerID    int         `json:"customer_id"`
	StatusID      int         `json:"status_id"`
	Quantity      int         `json:"quantity"`
	Amount        int         `json:"amount"`
	CreatedAt     time.Time   `json:"-"`
	UpdatedAt     time.Time   `json:"-"`
	Widget        Widget      `json:"widget"`
	Transaction   Transaction `json:"transaction"`
	Customer      Customer    `json:"customer"`
}

// Status is the type for all statuses
type Status struct {
	ID        int       `json:"id"`
	Name      int       `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// TransactionStatus is the type for all transaction status
type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      int       `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Transaction is the type for all transaction status
type Transaction struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	ExpiryMonth         int       `json:"expiry_month"`
	ExpiryYear          int       `json:"expiry_year"`
	PaymentIntent       string    `json:"payment_intent"`
	PaymentMethod       string    `json:"payment_method"`
	BankReturnCode      string    `json:"bank_return_code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	CreatedAt           time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
}

// User is the type for all users
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Customer is the type for all Customers
type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var widget Widget

	row := m.DB.QueryRowContext(ctx, `
	SELECT 
		id, name, description, inventory_level, price, coalesce(image, ''),
		is_recurring, plan_id,
		created_at, updated_at
	FROM 
		widgets 
	WHERE id = ?`, id)
	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Price,
		&widget.Image,
		&widget.IsRecurring,
		&widget.PlanID,
		&widget.CreatedAt,
		&widget.UpdatedAt,
	)
	if err != nil {
		return widget, err
	}

	return widget, nil
}

// InsertTransaction insert a new txn, and returns its id
func (m *DBModel) InsertTransaction(txn Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO transactions 
			(amount, currency, last_four, bank_return_code, expiry_month, expiry_year,
				payment_intent, payment_method,
			transaction_status_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.BankReturnCode,
		txn.ExpiryMonth,
		txn.ExpiryYear,
		txn.PaymentIntent,
		txn.PaymentMethod,
		txn.TransactionStatusID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// InsertOrder insert a new order, and returns its id
func (m *DBModel) InsertOrder(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO orders 
			(widget_id, transaction_id, status_id, quantity, customer_id,
			amount, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		order.WidgetID,
		order.TransactionID,
		order.StatusID,
		order.Quantity,
		order.CustomerID,
		order.Amount,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// InsertCustomer insert a new customer, and returns its id
func (m *DBModel) InsertCustomer(c Customer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO customers 
			(first_name, last_name, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		c.FirstName,
		c.LastName,
		c.Email,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetUserByEmail gets a user by email address
func (m *DBModel) GetUserByEmail(email string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	email = strings.ToLower(email)
	var u User

	row := m.DB.QueryRowContext(ctx, `
		SELECT
			id, first_name, last_name, email, password, created_at, updated_at
		FROM
			users
		WHERE
			email = ?
	`, email)

	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (m *DBModel) Authenticate(email, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = ?", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, err
	}

	return id, nil

}

func (m *DBModel) UpdatePasswordForUser(u User, hash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE users SET password = ? WHERE id = ?`

	_, err := m.DB.ExecContext(ctx, stmt, hash, u.ID)

	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) GetAllOrders() ([]*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var orders []*Order

	query := `SELECT
		o.id, o.widget_id, o.transaction_id, o.customer_id,
		o.status_id, o.quantity, o.amount, o.created_at, 
		o.updated_at, w.id, w.name, t.id, t.amount, t.currency,
		t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
		t.bank_return_code, c.id, c.first_name, c.last_name, c.email 
	FROM
		orders o
		LEFT JOIN widgets w ON (o.widget_id = w.id)
		LEFT JOIN transactions t ON (o.transaction_id = t.id)
		LEFT JOIN customers c ON (o.customer_id = c.id)
	WHERE 
		w.is_recurring = 0
	ORDER BY
		o.created_at DESC;`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.WidgetID,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.Widget.ID,
			&o.Widget.Name,
			&o.Transaction.ID,
			&o.Transaction.Amount,
			&o.Transaction.Currency,
			&o.Transaction.LastFour,
			&o.Transaction.ExpiryMonth,
			&o.Transaction.ExpiryYear,
			&o.Transaction.PaymentIntent,
			&o.Transaction.BankReturnCode,
			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	return orders, nil
}

func (m *DBModel) GetAllOrdersPaginated(pageSize, currentPage int) ([]*Order, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var orders []*Order

	offset := (currentPage - 1) * pageSize

	query := `SELECT
		o.id, o.widget_id, o.transaction_id, o.customer_id,
		o.status_id, o.quantity, o.amount, o.created_at, 
		o.updated_at, w.id, w.name, t.id, t.amount, t.currency,
		t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
		t.bank_return_code, c.id, c.first_name, c.last_name, c.email 
	FROM
		orders o
		LEFT JOIN widgets w ON (o.widget_id = w.id)
		LEFT JOIN transactions t ON (o.transaction_id = t.id)
		LEFT JOIN customers c ON (o.customer_id = c.id)
	WHERE 
		w.is_recurring = 0
	LIMIT ? 
	OFFSET ?`

	rows, err := m.DB.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.WidgetID,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.Widget.ID,
			&o.Widget.Name,
			&o.Transaction.ID,
			&o.Transaction.Amount,
			&o.Transaction.Currency,
			&o.Transaction.LastFour,
			&o.Transaction.ExpiryMonth,
			&o.Transaction.ExpiryYear,
			&o.Transaction.PaymentIntent,
			&o.Transaction.BankReturnCode,
			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		orders = append(orders, &o)
	}

	query = `
		SELECT COUNT(o.id)
			FROM orders o
			LEFT JOIN widgets w ON (o.widget_id = w.id)
		WHERE
			w.is_recurring = 0
	`

	var totalRecords int
	countRow := m.DB.QueryRowContext(ctx, query)
	err = countRow.Scan(&totalRecords)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := totalRecords / pageSize

	return orders, lastPage, totalRecords, nil
}

func (m *DBModel) GetAllSubscriptions() ([]*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var orders []*Order

	query := `SELECT
		o.id, o.widget_id, o.transaction_id, o.customer_id,
		o.status_id, o.quantity, o.amount, o.created_at, 
		o.updated_at, w.id, w.name, t.id, t.amount, t.currency,
		t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
		t.bank_return_code, c.id, c.first_name, c.last_name, c.email 
	FROM
		orders o
		LEFT JOIN widgets w ON (o.widget_id = w.id)
		LEFT JOIN transactions t ON (o.transaction_id = t.id)
		LEFT JOIN customers c ON (o.customer_id = c.id)
	WHERE 
		w.is_recurring = 1
	ORDER BY
		o.created_at DESC;`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.WidgetID,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.Widget.ID,
			&o.Widget.Name,
			&o.Transaction.ID,
			&o.Transaction.Amount,
			&o.Transaction.Currency,
			&o.Transaction.LastFour,
			&o.Transaction.ExpiryMonth,
			&o.Transaction.ExpiryYear,
			&o.Transaction.PaymentIntent,
			&o.Transaction.BankReturnCode,
			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	return orders, nil
}

func (m *DBModel) GetAllSubscriptionsPaginated(pageSize, currentPage int) ([]*Order, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var orders []*Order

	offset := (currentPage - 1) * pageSize

	query := `SELECT
		o.id, o.widget_id, o.transaction_id, o.customer_id,
		o.status_id, o.quantity, o.amount, o.created_at, 
		o.updated_at, w.id, w.name, t.id, t.amount, t.currency,
		t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
		t.bank_return_code, c.id, c.first_name, c.last_name, c.email 
	FROM
		orders o
		LEFT JOIN widgets w ON (o.widget_id = w.id)
		LEFT JOIN transactions t ON (o.transaction_id = t.id)
		LEFT JOIN customers c ON (o.customer_id = c.id)
	WHERE 
		w.is_recurring = 1
	LIMIT ? 
	OFFSET ?`

	rows, err := m.DB.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.WidgetID,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.Widget.ID,
			&o.Widget.Name,
			&o.Transaction.ID,
			&o.Transaction.Amount,
			&o.Transaction.Currency,
			&o.Transaction.LastFour,
			&o.Transaction.ExpiryMonth,
			&o.Transaction.ExpiryYear,
			&o.Transaction.PaymentIntent,
			&o.Transaction.BankReturnCode,
			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		orders = append(orders, &o)
	}

	query = `
		SELECT COUNT(o.id)
			FROM orders o
			LEFT JOIN widgets w ON (o.widget_id = w.id)
		WHERE
			w.is_recurring = 1
	`

	var totalRecords int
	countRow := m.DB.QueryRowContext(ctx, query)
	err = countRow.Scan(&totalRecords)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := totalRecords / pageSize

	return orders, lastPage, totalRecords, nil
}

func (m *DBModel) GetOrderByID(id int) (Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var o Order

	query := `SELECT
		o.id, o.widget_id, o.transaction_id, o.customer_id,
		o.status_id, o.quantity, o.amount, o.created_at, 
		o.updated_at, w.id, w.name, t.id, t.amount, t.currency,
		t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
		t.bank_return_code, c.id, c.first_name, c.last_name, c.email 
	FROM
		orders o
		LEFT JOIN widgets w ON (o.widget_id = w.id)
		LEFT JOIN transactions t ON (o.transaction_id = t.id)
		LEFT JOIN customers c ON (o.customer_id = c.id)
	WHERE 
		o.id = ?`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&o.ID,
		&o.WidgetID,
		&o.TransactionID,
		&o.CustomerID,
		&o.StatusID,
		&o.Quantity,
		&o.Amount,
		&o.CreatedAt,
		&o.UpdatedAt,
		&o.Widget.ID,
		&o.Widget.Name,
		&o.Transaction.ID,
		&o.Transaction.Amount,
		&o.Transaction.Currency,
		&o.Transaction.LastFour,
		&o.Transaction.ExpiryMonth,
		&o.Transaction.ExpiryYear,
		&o.Transaction.PaymentIntent,
		&o.Transaction.BankReturnCode,
		&o.Customer.ID,
		&o.Customer.FirstName,
		&o.Customer.LastName,
		&o.Customer.Email,
	)
	if err != nil {
		return o, err
	}

	return o, nil
}

func (m *DBModel) UpdateOrderStatus(id, statusID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE orders SET status_id = ? WHERE id = ?`

	_, err := m.DB.QueryContext(ctx, query, statusID, id)

	return err
}

func (m *DBModel) GetAllUsers() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var users []*User

	query := `
	SELECT 
		id, last_name, first_name, email, created_at, updated_at
	FROM
		users
	ORDER BY
		last_name, first_name					
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		err = rows.Scan(
			&u.ID,
			&u.LastName,
			&u.FirstName,
			&u.Email,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (m *DBModel) GetOneUser(id int) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var usr User

	query := `
	SELECT 
		id, last_name, first_name, email, created_at, updated_at
	FROM
		users
	WHERE id = ?				
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&usr.ID,
		&usr.LastName,
		&usr.FirstName,
		&usr.Email,
		&usr.CreatedAt,
		&usr.UpdatedAt,
	)
	if err != nil {
		return usr, err
	}

	return usr, nil
}

func (m *DBModel) EditUser(usr User) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	UPDATE users SET
		first_name = ?,
		last_name = ?,
		email = ?,
		updated_at = ?
	WHERE
		id = ?
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		usr.FirstName,
		usr.LastName,
		usr.Email,
		time.Now(),
		usr.ID,
	)
	if err != nil {
		return usr, err
	}
	return usr, nil

}

func (m *DBModel) AddUser(usr User, hash string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	INSERT INTO users (first_name, last_name, email, password, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		usr.FirstName,
		usr.LastName,
		usr.Email,
		hash,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return usr, err
	}
	return usr, nil

}

func (m *DBModel) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM users WHERE id = ?`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	return err
}
