package todo_test

import (
	"go-cmd-line/interacting/todo"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	l := todo.List{}
	task := "new task"
	l.Add(task)
	if l[0].Task != task {
		t.Errorf("expected %q but got %q.", task, l[0].Task)
	}
}

func TestComplete(t *testing.T) {
	l := todo.List{}
	task := "new task"
	l.Add(task)
	if l[0].Task != task {
		t.Errorf("expected %q but got %q.", task, l[0].Task)
	}
	if l[0].Done {
		t.Errorf("new tasks should not be completed")
	}
	l.Complete(1)
	if !l[0].Done {
		t.Errorf("the task should be completed.")
	}

}

func TestDelete(t *testing.T) {
	l := todo.List{}
	tasks := []string{"t1", "t2", "t3"}
	for _, t := range tasks {
		l.Add(t)
	}
	if l[0].Task != tasks[0] {
		t.Errorf("expected %q but got %q.", tasks[0], l[0].Task)
	}
	l.Delete(2)
	if len(l) != 2 {
		t.Errorf("expected length %d but got %d", 2, len(l))
	}
	if l[1].Task != tasks[2] {
		t.Errorf("expected %q but got %q", tasks[2], l[1].Task)
	}

}

func TestSaveGet(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}
	task := "new task"
	l1.Add(task)
	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("error creating temp files: %s", err)
	}
	defer os.Remove(tf.Name())
	if err := l1.Save(tf.Name()); err != nil {
		t.Fatalf("error saving list to file: %s", err)
	}
	if err := l2.Get(tf.Name()); err != nil {
		t.Fatalf("error getting list from file %s", err)
	}

	if l1[0].Task != l2[0].Task {
		t.Errorf("task %q should match task %q", l1[0].Task, l2[0].Task)
	}
}
