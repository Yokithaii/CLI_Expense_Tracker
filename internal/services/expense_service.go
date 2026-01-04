package services

import (
	"context"
	"fmt"

	"github.com/Yokithaii/CLI_Expense_Tracker/internal/dataaccess"
	"github.com/Yokithaii/CLI_Expense_Tracker/internal/entities"
)

type ExpenseService struct {
	expenseDAO dataaccess.ExpenseDAO // Для работы с БД сервису нужен будет data access object для наших expense-ов
}

func NewExpenseService(expenseDAO *dataaccess.ExpenseDAO) *ExpenseService {
	return &ExpenseService{*expenseDAO}
}

// Checks if expense exists.
func (s *ExpenseService) validateExpenseExistence(ctx context.Context, id int) error {
	exists, err := s.expenseDAO.ExpenseExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("no expense with given id")
	}
	return nil
}

// Adds a promoted expense (need to validate).
func (s *ExpenseService) CreateExpense(ctx context.Context, expense entities.Expense) error {
	return s.expenseDAO.CreateExpense(ctx, expense)
}

// Deletes an expense with promoted id.
func (s *ExpenseService) DeleteExpense(ctx context.Context, id int) error {
	err := s.validateExpenseExistence(ctx, id)
	if err != nil {
		return err
	}

	return s.expenseDAO.DeleteExpense(ctx, id)
}

// Replaces expense with promoted id with new given expense.
func (s *ExpenseService) UpdateExpense(ctx context.Context, id int, expense entities.Expense) error {
	err := s.validateExpenseExistence(ctx, id)
	if err != nil {
		return err
	}

	return s.expenseDAO.UpdateExpense(ctx, id, expense)
}

// Returns an array of all expenses available.
func (s *ExpenseService) GetAllExpenses(ctx context.Context) ([]entities.Expense, error) {
	expenses, err := s.expenseDAO.GetAllExpensesFromDb(ctx)
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

// Returns an amount of all expenses.
func (s *ExpenseService) Summary(ctx context.Context) (int, error) {
	sum, err := s.expenseDAO.Summary(ctx)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

// Wipes out db (and resets id's too!).
func (s *ExpenseService) ResetExpenses(ctx context.Context) error {
	return s.expenseDAO.ResetExpenses(ctx)
}
