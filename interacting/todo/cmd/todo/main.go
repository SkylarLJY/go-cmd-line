package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-cmd-line/interacting/todo"
	"io"
	"os"
	"strings"
)

var todoFileName = ".todo.json"

func getTasks(r io.Reader, args ...string) ([]string, error) {

	if len(args) > 0 {
		return []string{strings.Join(args, " ")}, nil
	}

	res := []string{}
	s := bufio.NewScanner(r)
	for s.Scan() {
		task := s.Text()
		if len(task) == 0 {
			break
		}
		res = append(res, task)
	}
	return res, nil
}

func main() {
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}
	_, err := os.Open(todoFileName)
	if err != nil {
		f, _ := os.Create(todoFileName)
		f.Write([]byte("[]"))
	}

	add := flag.Bool("add", false, "Add a task to the list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	delete := flag.Int("delete", 0, "Delete a task")
	verbose := flag.Bool("v", false, "Verbose output with data time info")
	flag.Parse()

	l := &todo.List{}
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	switch {
	case *list:
		if *verbose {
			l.PrintVerbose()
		} else {
			fmt.Print(l)
		}
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:

		tasks, err := getTasks(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, task := range tasks {
			l.Add(task)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *delete > 0:
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stderr, "invalid option")
		os.Exit(1)
	}

}
