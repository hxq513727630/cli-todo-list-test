package main

import (
	"bufio"   //讀取stdin
	"fmt"     //輸出與格式化
	"os"      //作業系統介面
	"strconv" //字串與整數互轉換
	"strings" //字串處理
)

// Task 代表一個代辦事項
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
} //struct tag 'json' 方便未來用endoing/json序列化or反序列化

// Store 抽象出儲存層的介面（便於之後換成檔案或 DB）
type Store interface {
	Add(title string) Task
	List() []Task
	Update(id int, title string) error
	Delete(id int) error
	Toggle(id int) error
}

// MemoryStore 是記憶體的實作
type MemoryStore struct {
	tasks  []Task
	nextID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{tasks: []Task{}, nextID: 1}
}

func (s *MemoryStore) Add(title string) Task {
	t := Task{ID: s.nextID, Title: title, Done: false}
	s.nextID++
	s.tasks = append(s.tasks, t)
	return t
}

func (s *MemoryStore) List() []Task {
	return s.tasks
}

func (s *MemoryStore) findIndex(id int) int {
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			return i
		}
	}
	return -1
}

func (s *MemoryStore) Update(id int, title string) error {
	idx := s.findIndex(id)
	if idx == -1 {
		return fmt.Errorf("task %d not found", id)
	}
	s.tasks[idx].Title = title
	return nil
}

func (s *MemoryStore) Delete(id int) error {
	idx := s.findIndex(id)
	if idx == -1 {
		return fmt.Errorf("task %d not found", id)
	}
	s.tasks = append(s.tasks[:idx], s.tasks[idx+1:]...)
	return nil
}

func (s *MemoryStore) Toggle(id int) error {
	idx := s.findIndex(id)
	if idx == -1 {
		return fmt.Errorf("task %d not found", id)
	}
	s.tasks[idx].Done = !s.tasks[idx].Done
	return nil
}

// CLI helper
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
	store := NewMemoryStore()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Todo CLI — type 'help' for commands")
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "read error:", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 取出指令（第一個字）與其餘內容
		parts := strings.Fields(line)
		cmd := parts[0]

		switch cmd {
		case "help":
			printHelp()
		case "add":
			// 支援標題內含空白：取整行的剩餘部分
			title := strings.TrimSpace(line[len("add"):])
			if title == "" {
				fmt.Println("usage: add <title>")
				continue
			}
			t := store.Add(title)
			fmt.Printf("added: %d - %s\n", t.ID, t.Title)
		case "list":
			tasks := store.List()
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
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("invalid id")
				continue
			}
			// 取得 id 後面的標題（保留空白）
			title := strings.TrimSpace(line[len("update ")+len(parts[1]):])
			if title == "" {
				fmt.Println("empty title")
				continue
			}
			if err := store.Update(id, title); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("updated")
			}
		case "delete":
			if len(parts) != 2 {
				fmt.Println("usage: delete <id>")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("invalid id")
				continue
			}
			if err := store.Delete(id); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("deleted")
			}
		case "done":
			if len(parts) != 2 {
				fmt.Println("usage: done <id>")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("invalid id")
				continue
			}
			if err := store.Toggle(id); err != nil {
				fmt.Println(err)
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
