package main

import (
	"io/fs"
	"io/ioutil"
	"os/exec"
	"path"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func getBookList(path string) []fs.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	return files
}

func tableGenerate() (*tview.Table, []string) {

	table := tview.NewTable()
	filesInfo := getBookList("/home/nguyen/Books/")

	//Read file name
	var fileNames []string
	for _, file := range filesInfo {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		} else {
			tempFilesInfo := getBookList("/home/nguyen/Books/" +file.Name())
			for _, nestedFile := range tempFilesInfo {
                //filesName contain the directory name also
				fileNames = append(fileNames,file.Name()+"/"+nestedFile.Name())
			}
		}
	}

    //Sort the fileNames
    sort.Strings(fileNames)

	// Generate the table
	cols, rows := 1, len(fileNames)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
            cellName := path.Base(fileNames[r])
			table.SetCell(r, c, tview.NewTableCell(cellName).SetTextColor(color))
		}
	}

	return table, fileNames
}

func selectItem(app *tview.Application, table *tview.Table,fileNames []string) {
	table.Select(0, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {  
        fileName := fileNames[row]
        cmd := exec.Command("zathura","/home/nguyen/Books/" + fileName)    
        cmd.Run()
		table.SetSelectable(false, false)
	})
}

func main() {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	app := tview.NewApplication()
	// main := newPrimitive("Main content")
	bookList,fileNames := tableGenerate()
	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0).
		SetBorders(true).
        AddItem(newPrimitive("Kiwi Library"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("Created by Shiny Kiwis - 2022"), 2, 0, 1, 3, 0, 0, false)

    // Layout for screen less than 100 cells
	grid.AddItem(bookList, 1, 0, 1, 3, 0, 0, false)
    // Layout for screen wider than 100 cells
	grid.AddItem(bookList, 1, 0, 1, 3, 0, 100, false)

    // Function for selecting books
	selectItem(app, bookList,fileNames)

    //Todo
    //Allow to choose between books and document(latex files)
	if err := app.SetRoot(grid, true).SetFocus(bookList).Run(); err != nil {
		panic(err)
	}
}
