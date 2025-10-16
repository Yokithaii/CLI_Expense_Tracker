package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Yokithaii/CLI_Expense_Tracker/internal/operator"

	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := operator.NewConnectionToDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(context.Background())

	scanner := bufio.NewScanner(os.Stdin)

	operator.Help()

	for {
		fmt.Print("Введите команду: ")

		if ok := scanner.Scan(); !ok {
			fmt.Println("Ошибка ввода!")
			continue
		}
		text := scanner.Text()
		fields := strings.Fields(text)

		err = processCommand(fields, conn)
		if err != nil {
			if err.Error() == "exit" {
				break
			}
			fmt.Println(err)
			continue
		}
	}
}

func processCommand(fields []string, conn *pgx.Conn) error {
	if len(fields) == 0 {
		return fmt.Errorf("вы ничего не ввели")
	}
	cmd := fields[0]
	switch cmd {
	case "exit":
		return fmt.Errorf("exit")
	case "add":
		return addExpense(conn, fields)
	case "list":
		return listExpenses(conn)
	case "summary":
		err := summaryOfExpenses(conn)
		if err != nil {
			return err
		}
	case "help":
		operator.Help()
	case "delete":
		err := deleteExpense(conn, fields)
		if err != nil {
			return err
		}
	case "update":
		err := updateExpense(conn, fields)
		if err != nil {
			return err
		}
	case "reset":
		err := operator.ResetDatabase(conn)
		if err != nil {
			return err
		}
	default:
		unknownCommand()
	}
	return nil
}
func unknownCommand() {
	fmt.Println("No such command!")
	operator.Help()
}

func listExpenses(conn *pgx.Conn) error {
	expenses, err := operator.GetAllExpensesFromDb(conn)
	if err != nil {
		return err
	}
	printList(expenses)
	return nil
}
func summaryOfExpenses(conn *pgx.Conn) error {
	summ, err := operator.Summary(conn)
	if err != nil {
		return err
	}
	fmt.Println(summ)
	return nil
}

func addExpense(conn *pgx.Conn, fields []string) error {
	if len(fields) != 3 {
		return fmt.Errorf("wrong arguments provided for command \"add\"")
	}

	amount, err := strconv.Atoi(fields[2])
	if err != nil {
		return err
	}

	description := fields[1]
	err = operator.ValidateData(description, amount)
	if err != nil {
		return err
	}

	err = operator.CreateExpense(conn, operator.NewExpense(fields[1], amount))
	if err != nil {
		return err
	}
	fmt.Println("Expense added succesfully")
	return nil
}

func printList(expenses []operator.Expense) {
	fmt.Printf("%-4s %-12s %-20s %s\n", "#", "Date", "Description", "Amount")
	fmt.Println("--------------------------------------------------")
	for _, v := range expenses {
		fmt.Printf("%-4v %-12v %-20v $%d\n",
			v.ID,
			v.Date.Format("2006-01-02"),
			v.Desc,
			v.Amount)
	}
}

func deleteExpense(conn *pgx.Conn, fields []string) error {
	if len(fields) != 2 {
		return fmt.Errorf("you provided wrong arguments to \"delete\"")
	}
	id, err := strconv.Atoi(fields[1])
	if err != nil {
		return err
	}
	exists, err := operator.ExpenseExists(conn, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("no expense with given id")
	}
	err = operator.DeleteExpense(conn, id)
	if err != nil {
		return err
	}
	fmt.Printf("Expense ID: %v deleted succesfully\n", id)
	return nil
}

func updateExpense(conn *pgx.Conn, fields []string) error {
	if len(fields) != 4 {
		return fmt.Errorf("not enough arguments provided for command \"update\"")
	}

	id, err := strconv.Atoi(fields[1])
	if err != nil {
		return err
	}

	exists, err := operator.ExpenseExists(conn, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("no expense with provided id")
	}

	amount, err := strconv.Atoi(fields[3])
	if err != nil {
		return err
	}

	err = operator.ValidateData(fields[2], amount)
	if err != nil {
		return err
	}

	updatedExpense := operator.NewExpense(fields[2], amount)

	err = operator.UpdateExpense(conn, id, updatedExpense)
	if err != nil {
		return err
	}
	fmt.Printf("Id %v updated succesfully!\n", id)
	return nil
}
