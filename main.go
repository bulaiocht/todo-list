package main

import (
	"flag"
	"fmt"
	"log"
	"todo-list/todo"
)

const (
	todoFile = "default"
)

func main() {
	list, err := todo.LoadFromFile(todoFile)
	if err != nil {
		panic(err)
	}
	defer func(l *todo.List) {
		e := todo.SaveToFile(l)
		if e != nil {
			log.Fatal(e)
		}
	}(list)

	show := *flag.Bool("l", false, "show tasks from list")
	add := *flag.String("a", "", "add new task into list")
	del := *flag.Int("rm", -1, "delete task from the list")
	done := *flag.Int("D", -1, "mark task as Done")

	flag.Parse()

	switch {
	case show:
		fmt.Println(list)
	case len(add) > 0:
		list.Add(add)
	case del > 0:
		list.Remove(del)
	case done > 0:
		list.MarkAsDone(done)
	}
}
