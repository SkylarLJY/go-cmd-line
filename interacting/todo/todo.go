package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

func (l *List) Add(task string) {
	t := item{task, false, time.Now(), time.Time{}}
	*l = append(*l, t)
}

func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}
	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()
	return nil
}

func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}
	*l = append(ls[:i-1], ls[i:]...)
	return nil
}

func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, js, 0644)
}

func (l *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, l)
}

// Implement fmt.Stringer interface
func (l *List) String() string {
	formatted := ""
	for k, v := range *l {
		prefix := " "
		if v.Done {
			prefix = "X "
		}
		formatted += fmt.Sprintf("%s%d: %s\n", prefix, k+1, v.Task)
	}
	return formatted
}

func (l *List) PrintVerbose() {
	str := ""
	for k, v := range *l {
		prefix := " "
		if v.Done {
			prefix = "X "
		}
		str += fmt.Sprintf("%s%d: %s - created at %v\n", prefix, k+1, v.Task, v.CreatedAt)
	}
	fmt.Print(str)
}
