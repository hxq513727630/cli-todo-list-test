package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"todoList_db/prisma/db" // 匯入生成的 Prisma Client
)

type Store struct {
	client *db.PrismaClient
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func NewStore() *Store {
	client := db.NewClient()
	if err := client.Connect(); err != nil {
		panic(err)
	}
	return &Store{client: client}
}

func (s *Store) Add(title string) error {
	_, err := s.client.Task.CreateOne(
		db.Task.Title.Set(title),
	).Exec(context.Background())
	return err
}

func (s *Store) List() ([]db.TaskModel, error) {
	return s.client.Task.FindMany().Exec(context.Background())
}

func (s *Store) Update(id int, title string) error {
	_, err := s.client.Task.FindUnique(
		db.Task.ID.Equals(id),
	).Update(
		db.Task.Title.Set(title),
	).Exec(context.Background())
	return err
}

func (s *Store) Delete(id int) error {
	_, err := s.client.Task.FindUnique(
		db.Task.ID.Equals(id),
	).Delete().Exec(context.Background())
	return err
}

func (s *Store) Toggle(id int) error {
	task, err := s.client.Task.FindUnique(
		db.Task.ID.Equals(id),
	).Exec(context.Background())
	if err != nil || task == nil {
		return fmt.Errorf("task %d not found", id)
	}
	_, err = s.client.Task.FindUnique(
		db.Task.ID.Equals(id),
	).Update(
		db.Task.Done.Set(!task.Done),
	).Exec(context.Background())
	return err
}

// CLI Helper
func printHelp() {
	fmt.Println("Commands:")
	fmt.Println("  add <title>           - add new task")
	fmt.Println("  list                  - list tasks")
	fmt.Println("  update <id> <title>   - update task title")
	fmt.Println("  delete <id>           - delete task")
	fmt.Println("  done <id>             - toggle task done/undone")
	fmt.Println("  help                  - show this help")
	fmt.Println("  exit                  - quit")
}

func main() {
	store := NewStore()
	defer store.client.Disconnect()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Todo CLI — type 'help' for commands")

	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		cmd := parts[0]

		switch cmd {
		case "help":
			printHelp()
		case "add":
			title := strings.TrimSpace(line[len("add"):])
			if title == "" {
				fmt.Println("usage: add <title>")
				continue
			}
			if err := store.Add(title); err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println("added:", title)
			}
		case "list":
			tasks, err := store.List()
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
			if len(tasks) == 0 {
				fmt.Println("no tasks")
				continue
			}
			for _, t := range tasks {
				status := "[ ]"
				if t.Done {
					status = "[x]"
				}
				fmt.Printf("%d. %s %s\n", t.ID, status, t.Title)
			}
		case "update":
			if len(parts) < 3 {
				fmt.Println("usage: update <id> <title>")
				continue
			}
			id, _ := strconv.Atoi(parts[1])
			title := strings.TrimSpace(line[len("update ")+len(parts[1]):])
			if err := store.Update(id, title); err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println("updated")
			}
		case "delete":
			if len(parts) != 2 {
				fmt.Println("usage: delete <id>")
				continue
			}
			id, _ := strconv.Atoi(parts[1])
			if err := store.Delete(id); err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println("deleted")
			}
		case "done":
			if len(parts) != 2 {
				fmt.Println("usage: done <id>")
				continue
			}
			id, _ := strconv.Atoi(parts[1])
			if err := store.Toggle(id); err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println("toggled")
			}
		case "exit":
			fmt.Println("bye")
			return
		default:
			fmt.Println("unknown command:", cmd)
		}
	}
}
