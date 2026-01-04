package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Yokithaii/CLI_Expense_Tracker/internal/dataaccess"
	"github.com/Yokithaii/CLI_Expense_Tracker/internal/entities"
	"github.com/Yokithaii/CLI_Expense_Tracker/internal/services"
)

func main() {
	fmt.Print("Введите строку для подключения к бд: ")

	ctx := context.Background()

	scanner := bufio.NewScanner(os.Stdin)
	if ok := scanner.Scan(); !ok {
		fmt.Println("Ошибка ввода!")
		return
	}

	name := scanner.Text()
	name = strings.TrimSpace(string(name))

	expDAO, err := dataaccess.NewExpenseDAO(ctx, name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer expDAO.CloseConnect(ctx)

	service := services.NewExpenseService(expDAO)
	fmt.Println('\n')

	Help()

	for {
		fmt.Print("Введите команду: ")

		if ok := scanner.Scan(); !ok {
			fmt.Println("Ошибка ввода!")
			continue
		}
		text := scanner.Text()
		fields := strings.Fields(text)

		err = processCommand(ctx, fields, service)
		if err != nil {
			if err.Error() == "exit" {
				break
			}
			fmt.Println(err)
			continue
		}
	}
}

func Help() {
	fmt.Println("Доступные команды: 1. add *Description* *amount* -- adds your expense")
	fmt.Println("Доступные команды: 2. list -- shows list of your expenses")
	fmt.Println("Доступные команды: 3. summary -- shows amount you spent")
	fmt.Println("Доступные команды: 4. help -- shows a list of commands")
	fmt.Println("Доступные команды: 5. delete *id* -- deletes expense with given id(rest of id's stays the same)")
	fmt.Println("Доступные команды: 6. update *id* *description* *amount* -- updates id if exists")
	fmt.Println("Доступные команды: 7. exit -- ends a program")
	fmt.Println("Доступные команды: 8. reset -- resets database")
}

func processCommand(ctx context.Context, fields []string, s *services.ExpenseService) error {
	if len(fields) == 0 {
		return fmt.Errorf("вы ничего не ввели")
	}
	cmd := fields[0]
	switch cmd {
	case "exit":
		return fmt.Errorf("exit")
	case "add":
		return addExpense(ctx, s, fields)
	case "list":
		return listExpenses(ctx, s)
	case "summary":
		err := summaryOfExpenses(ctx, s)
		if err != nil {
			return err
		}
	case "help":
		Help()
	case "delete":
		err := deleteExpense(ctx, s, fields)
		if err != nil {
			return err
		}
	case "update":
		err := updateExpense(ctx, s, fields)
		if err != nil {
			return err
		}
	case "reset":
		err := s.ResetExpenses(ctx)
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
	Help()
}

func listExpenses(ctx context.Context, s *services.ExpenseService) error {
	expenses, err := s.GetAllExpenses(ctx)
	if err != nil {
		return err
	}
	printList(expenses)
	return nil
}
func summaryOfExpenses(ctx context.Context, s *services.ExpenseService) error {
	summ, err := s.Summary(ctx)
	if err != nil {
		return err
	}
	fmt.Println(summ)
	return nil
}

func addExpense(ctx context.Context, s *services.ExpenseService, fields []string) error {
	if len(fields) != 3 {
		return fmt.Errorf("wrong arguments provided for command \"add\"")
	}

	amount, err := strconv.Atoi(fields[2])
	if err != nil {
		return err
	}

	description := fields[1]
	newExpense, err := entities.NewExpense(description, amount)
	if err != nil {
		return err
	}

	err = s.CreateExpense(ctx, *newExpense)
	if err != nil {
		return err
	}
	fmt.Println("Expense added succesfully")
	return nil
}

func printList(expenses []entities.Expense) {
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

func deleteExpense(ctx context.Context, s *services.ExpenseService, fields []string) error {
	if len(fields) != 2 {
		return fmt.Errorf("you provided wrong arguments to \"delete\"")
	}
	id, err := strconv.Atoi(fields[1])
	if err != nil {
		return err
	}

	err = s.DeleteExpense(ctx, id)
	if err != nil {
		return err
	}
	fmt.Printf("Expense ID: %v deleted succesfully\n", id)
	return nil
}

func updateExpense(ctx context.Context, s *services.ExpenseService, fields []string) error {
	if len(fields) != 4 {
		return fmt.Errorf("wrong arguments provided for command \"update\"")
	}

	id, err := strconv.Atoi(fields[1])
	if err != nil {
		return err
	}

	amount, err := strconv.Atoi(fields[3])
	if err != nil {
		return err
	}

	updatedExpense, err := entities.NewExpense(fields[2], amount)
	if err != nil {
		return err
	}

	err = s.UpdateExpense(ctx, id, *updatedExpense)
	if err != nil {
		return err
	}
	fmt.Printf("Id %v updated succesfully!\n", id)
	return nil
}
