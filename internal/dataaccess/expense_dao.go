package dataaccess

import (
	"context"
	"fmt"

	"github.com/Yokithaii/CLI_Expense_Tracker/internal/entities"
	"github.com/jackc/pgx/v5"
)

type ExpenseDAO struct {
	conn *pgx.Conn
}

func (e *ExpenseDAO) CloseConnect(ctx context.Context) {
	e.conn.Close(ctx)
}

// Returns an object that manages expenses data in a database.
// postgres://postgres:Asdewq232113@localhost:5433/postgres
func NewExpenseDAO(ctx context.Context, connectionString string) (*ExpenseDAO, error) {
	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	return &ExpenseDAO{conn: conn}, err
}

// Checks that expense if provided id exists.
func (e *ExpenseDAO) ExpenseExists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := e.conn.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM public.Expenses WHERE id = $1)`, id).Scan(&exists)
	return exists, err
}

// Adds a promoted expense (need to validate).
func (e *ExpenseDAO) CreateExpense(ctx context.Context, exp entities.Expense) error {
	query := `INSERT INTO public.Expenses (Date, Description, Amount) VALUES (@Date, @Description, @Amount)`

	args := pgx.NamedArgs{
		"Date":        exp.Date,
		"Description": exp.Desc,
		"Amount":      exp.Amount,
	}

	_, err := e.conn.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error inserting expense: %w", err)
	}
	return nil
}

// Returns an array of all expenses available.
func (e *ExpenseDAO) GetAllExpensesFromDb(ctx context.Context) ([]entities.Expense, error) {
	query := `SELECT Id, Date, Description, Amount FROM public.Expenses`

	rows, err := e.conn.Query(ctx, query)
	if err != nil {
		fmt.Println("Error querying the table")
		return nil, err
	}

	defer rows.Close()

	var expenses []entities.Expense

	for rows.Next() {
		var exp entities.Expense
		err := rows.Scan(&exp.ID, &exp.Date, &exp.Desc, &exp.Amount)
		if err != nil {
			fmt.Println("Error fetching expense details")
			return expenses, err
		}
		expenses = append(expenses, exp)
	}
	return expenses, nil
}

// Replaces expense with promoted id with new given expense.
func (e *ExpenseDAO) UpdateExpense(ctx context.Context, id int, exp entities.Expense) error {
	query := `
		UPDATE public.Expenses
		SET Date = @Date, Description = @Description, Amount = @Amount 
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":          id,
		"Date":        exp.Date,
		"Description": exp.Desc,
		"Amount":      exp.Amount,
	}

	_, err := e.conn.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

// Deletes an expense with promoted id.
func (e *ExpenseDAO) DeleteExpense(ctx context.Context, id int) error {
	query := `
	DELETE FROM public.Expenses WHERE id = @id`

	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := e.conn.Exec(ctx, query, args)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Returns an amount of all expenses.
func (e *ExpenseDAO) Summary(ctx context.Context) (int, error) {
	query := `
	SELECT SUM(e.Amount) FROM public.Expenses as e
	`

	rows, err := e.conn.Query(ctx, query)
	if err != nil {
		fmt.Println("Error getting summary")
		return 0, err
	}
	defer rows.Close()

	var amount int
	rows.Next()
	err = rows.Scan(&amount)
	if err != nil {
		return 0, err
	}
	return amount, nil
}

// Wipes out db (and resets id's too!).
func (e *ExpenseDAO) ResetExpenses(ctx context.Context) error {
	_, err := e.conn.Exec(ctx,
		`TRUNCATE TABLE public.Expenses RESTART IDENTITY`)
	return err
}
