package main

import (
	"ExpenseTracker/operator"
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

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

		fields := strings.Fields(scanner.Text())

		if len(fields) == 0 {
			fmt.Println("Вы ничего не ввели")
			continue
		}

		cmd := fields[0]

		switch cmd {
		case "exit":
			return

		case "add":
			err = addExpense(conn, fields)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Expense added succesfully")

		case "list":
			expenses, err := operator.GetAllExpensesFromDb(conn)
			if err != nil {
				fmt.Println(err)
				continue
			}
			printList(expenses)

		case "summary":
			summary, err := operator.Summary(conn)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(summary)

		case "help":
			operator.Help()

		case "delete":
			err = deleteExpense(conn, fields)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case "update":
			err = UpdateExpense(conn, fields)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case "reset":
			err = operator.ResetDatabase(conn)
			if err != nil {
				fmt.Println(err)
				continue
			}
		default:
			unknownCommand()

		}

	}

}

func unknownCommand() {
	fmt.Println("No such command!")
	operator.Help()
}

func addExpense(conn *pgx.Conn, fields []string) error {
	if len(fields) < 3 {
		return fmt.Errorf("not enough arguments provided for command \"add\"")
	}

	temp, err := strconv.Atoi(fields[2])
	if err != nil {
		return err
	}

	err = operator.ValidateData(fields[1], temp)
	if err != nil {
		return err
	}

	err = operator.AddExpensToDatabase(conn, operator.CreateExpense(fields[1], temp))
	if err != nil {
		return err
	} else {
		return nil
	}
}

func printList(expenses []operator.Expense) {
	fmt.Printf("%-4s %-12s %-20s %s\n", "#", "Date", "Description", "Amount")
	fmt.Println("--------------------------------------------------")
	for _, v := range expenses {
		fmt.Printf("%-4v %-12v %-20v $%d\n",
			v.Id,
			v.Date.Format("2006-01-02"),
			v.Desc,
			v.Amount)
	}
}

func deleteExpense(conn *pgx.Conn, fields []string) error {

	if len(fields) < 2 {
		return fmt.Errorf("you provided no id to delete")
	}
	id, err := strconv.Atoi(fields[1])
	if err != nil {
		return err
	} else {
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
		} else {
			fmt.Printf("Expense ID: %v deleted succesfully\n", id)
			return nil
		}

	}

}

func UpdateExpense(conn *pgx.Conn, fields []string) error {
	if len(fields) < 4 {
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

	updatedExpense := operator.CreateExpense(fields[2], amount)

	err = operator.UpdateExpense(conn, id, updatedExpense)
	if err != nil {
		return err
	}
	fmt.Printf("Id %v updated succesfully!\n", id)
	return nil
}
