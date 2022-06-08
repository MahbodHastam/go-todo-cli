package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []item

func (t *Todos) Add(task string) {
	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Now(),
	}

	*t = append(*t, todo)
}

func (t *Todos) Complete(idx int) error {
	list := *t
	if idx <= 0 || idx > len(list) {
		return errors.New("invalid index")
	}

	if list[idx-1].Done {
		return errors.New("it's already marked as completed")
	}

	list[idx-1].CompletedAt = time.Now()
	list[idx-1].Done = true

	return nil
}

func (t *Todos) Delete(idx int) error {
	list := *t
	if idx <= 0 || idx > len(list) {
		return errors.New("invalid index")
	}

	*t = append(list[:idx-1], list[idx:]...)

	return nil
}

func (t *Todos) Load(fileName string) error {
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) Store(fileName string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, data, 0644)
}

func (t *Todos) Print() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Completed"},
			{Align: simpletable.AlignCenter, Text: "Created At"},
			{Align: simpletable.AlignCenter, Text: "Completed At"},
			{Align: simpletable.AlignCenter, Text: "Completed In"},
		},
	}

	var cells [][]*simpletable.Cell

	for idx, item := range *t {
		idx++
		task := blue(item.Task)
		diff := ""

		if item.Done {
			task = green(fmt.Sprintf("✅ %s", item.Task))

			diffOp := item.CreatedAt.Day() - item.CompletedAt.Day()

			if diffOp != 0 {
				diff = fmt.Sprintf("%d days", diffOp)
			} else {
				diffOp = item.CreatedAt.Hour() - item.CompletedAt.Hour()

				if diffOp != 0 {
					diff = fmt.Sprintf("%d hours", item.CompletedAt.Hour()-item.CreatedAt.Hour())
				} else {
					diff = fmt.Sprintf("%d minutes", item.CompletedAt.Minute()-item.CreatedAt.Minute())
				}
			}
		}

		completedAt := item.CompletedAt.Format(time.RFC822)
		if item.CreatedAt.Unix() == item.CompletedAt.Unix() {
			completedAt = "..."
		}

		cells = append(cells, []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%d", idx)},
			{Align: simpletable.AlignCenter, Text: task},
			{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%t", item.Done)},
			{Align: simpletable.AlignCenter, Text: item.CreatedAt.Format(time.RFC822)},
			{Align: simpletable.AlignCenter, Text: completedAt},
			{Align: simpletable.AlignCenter, Text: diff},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}
	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 6, Text: red(fmt.Sprintf("%d pending todos", t.CountPending()))},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func (t *Todos) CountPending() int {
	total := 0

	for _, item := range *t {
		if !item.Done {
			total++
		}
	}

	return total
}
