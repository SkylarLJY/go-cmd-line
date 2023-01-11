package main

import (
	"flag"
	"fmt"
	"go-cmd-line/interacting/todo"
	"os"
)

var todoFileName = ".todo.json"

// func init() {
// 	_, err := os.Open(todoFileName)
// 	if err != nil {
// 		f, _ := os.Create(todoFileName)
// 		f.Write([]byte("[]"))
// 	}
// }

func main() {
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}
	_, err := os.Open(todoFileName)
	if err != nil {
		f, _ := os.Create(todoFileName)
		f.Write([]byte("[]"))
	}
	task := flag.String("task", "", "Task to be added to the list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	flag.Parse()

	l := &todo.List{}
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	switch {
	case *list:
		fmt.Print(l)
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *task != "":
		l.Add(*task)

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stderr, "invalid option")
		os.Exit(1)
	}

}
