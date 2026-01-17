package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sinjitonayo/task-cli-go/internal/model"
	"github.com/sinjitonayo/task-cli-go/internal/storage"
)

type Handler struct {
	store *storage.JSONStore
}

func NewHandler(store *storage.JSONStore) *Handler {
	return &Handler{store: store}
}

// Run = entry to run CLI handler
// args = os.Args[1:]
func (h *Handler) Run(args []string) {
	if len(args) == 0 {
		h.PrintHelp()
		return
	}

	cmd := args[0]

	switch cmd {
	case "add":
		h.Add(args[1:])
	case "list":
		h.List(args[1:])
	case "update":
		h.Update(args[1:])
	case "delete":
		h.Delete(args[1:])
	case "mark-in-progress":
		h.MarkStatus(args[1:], model.StatusInProgress)
	case "mark-done":
		h.MarkStatus(args[1:], model.StatusDone)
	case "help":
		h.PrintHelp()
	default:
		fmt.Println("❌ Unknown command:", cmd)
		h.PrintHelp()
	}
}

func (h *Handler) Add(args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Usage: task-cli add \"Task description\"")
		return
	}

	desc := strings.TrimSpace(strings.Join(args, " "))
	if desc == "" {
		fmt.Println("❌ Description cannot be empty")
		return
	}

	tasks, err := h.store.LoadTasks()
	if err != nil {
		fmt.Println("❌ Error loading tasks:", err)
		return
	}

	newID := getNextID(tasks)

	now := time.Now()
	newTask := model.Task{
		ID:          newID,
		Description: desc,
		Status:      model.StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, newTask)

	if err := h.store.SaveTasks(tasks); err != nil {
		fmt.Println("❌ Error saving tasks:", err)
		return
	}

	fmt.Printf("✅ Task added successfully (ID: %d)\n", newID)
}

func (h *Handler) List(args []string) {
	tasks, err := h.store.LoadTasks()
	if err != nil {
		fmt.Println("❌ Error loading tasks:", err)
		return
	}

	var filter *model.TaskStatus = nil

	// if there is a status filter argument
	if len(args) >= 1 {
		statusStr := strings.TrimSpace(args[0])
		status, ok := parseStatus(statusStr)
		if !ok {
			fmt.Println("❌ Invalid status. Allowed: todo, in-progress, done")
			return
		}
		filter = &status
	}

	if len(tasks) == 0 {
		fmt.Println("(empty) No tasks found.")
		return
	}

	// print tasks
	for _, t := range tasks {
		if filter != nil && t.Status != *filter {
			continue
		}

		fmt.Printf(
			"[%d] (%s) %s | updated: %s\n",
			t.ID,
			t.Status,
			t.Description,
			t.UpdatedAt.Format("2006-01-02 15:04"),
		)
	}
}

func (h *Handler) Update(args []string) {
	if len(args) < 2 {
		fmt.Println("❌ Usage: task-cli update <id> \"New description\"")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("❌ ID must be a number")
		return
	}

	newDesc := strings.TrimSpace(strings.Join(args[1:], " "))
	if newDesc == "" {
		fmt.Println("❌ Description cannot be empty")
		return
	}

	tasks, err := h.store.LoadTasks()
	if err != nil {
		fmt.Println("❌ Error loading tasks:", err)
		return
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Description = newDesc
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		fmt.Println("❌ Task not found")
		return
	}

	if err := h.store.SaveTasks(tasks); err != nil {
		fmt.Println("❌ Error saving tasks:", err)
		return
	}

	fmt.Println("✅ Task updated successfully")
}

func (h *Handler) Delete(args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Usage: task-cli delete <id>")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("❌ ID must be a number")
		return
	}

	tasks, err := h.store.LoadTasks()
	if err != nil {
		fmt.Println("❌ Error loading tasks:", err)
		return
	}

	newTasks := make([]model.Task, 0, len(tasks))
	found := false

	for _, t := range tasks {
		if t.ID == id {
			found = true
			continue
		}
		newTasks = append(newTasks, t)
	}

	if !found {
		fmt.Println("❌ Task not found")
		return
	}

	if err := h.store.SaveTasks(newTasks); err != nil {
		fmt.Println("❌ Error saving tasks:", err)
		return
	}

	fmt.Println("✅ Task deleted successfully")
}

func (h *Handler) MarkStatus(args []string, status model.TaskStatus) {
	if len(args) < 1 {
		fmt.Printf("❌ Usage: task-cli mark-%s <id>\n", status)
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("❌ ID must be a number")
		return
	}

	tasks, err := h.store.LoadTasks()
	if err != nil {
		fmt.Println("❌ Error loading tasks:", err)
		return
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		fmt.Println("❌ Task not found")
		return
	}

	if err := h.store.SaveTasks(tasks); err != nil {
		fmt.Println("❌ Error saving tasks:", err)
		return
	}

	fmt.Printf("✅ Task marked as %s\n", status)
}

func (h *Handler) PrintHelp() {
	fmt.Println("Task CLI - Commands:")
	fmt.Println("  task-cli add \"Task description\"")
	fmt.Println("  task-cli update <id> \"New description\"")
	fmt.Println("  task-cli delete <id>")
	fmt.Println("  task-cli mark-in-progress <id>")
	fmt.Println("  task-cli mark-done <id>")
	fmt.Println("  task-cli list")
	fmt.Println("  task-cli list todo|in-progress|done")
	fmt.Println("  task-cli help")
}

// helper: generate ID
func getNextID(tasks []model.Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

// helper: parse status string
func parseStatus(s string) (model.TaskStatus, bool) {
	switch s {
	case "todo":
		return model.StatusTodo, true
	case "in-progress":
		return model.StatusInProgress, true
	case "done":
		return model.StatusDone, true
	default:
		return "", false
	}
}
