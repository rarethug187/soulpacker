package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var version = "1.0.0"

func main() {
	a := app.New()
	w := a.NewWindow("SoulPacker by rare thug")
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(500, 163))
	w.CenterOnScreen()

	folderEntry := widget.NewEntry()
	folderEntry.SetPlaceHolder("Path to folder...")
	folderEntry.Resize(fyne.NewSize(451, 36))

	folderButton := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		newWindow := a.NewWindow("Select folder...")
		newWindow.Resize(fyne.NewSize(500, 380))
		newWindow.CenterOnScreen()
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				log.Println("Error opening folder:", err)
				return
			}
			folderPath := uri.Path()
			folderEntry.SetText(folderPath)
			newWindow.Close()
		}, newWindow)
		newWindow.SetFixedSize(true)
		newWindow.Show()
	})
	folderButton.Resize(fyne.NewSize(35, 35))
	folderButton.Move(fyne.NewPos(456, 1))

	var extension string
	extensionRadio := widget.NewRadioGroup([]string{".wav", ".mp3", ".flac"}, func(s string) {
		extension = s
	})
	extensionRadio.Horizontal = true

	filenameEntry := widget.NewEntry()
	filenameEntry.SetPlaceHolder("Enter a file name...")

	var num int
	renameButton := widget.NewButton("Rename Files", func() {
		folderPath := folderEntry.Text
		if folderPath == "" {
			err := errors.New("Please select a folder")
			dialog := dialog.NewError(err, w)
			dialog.Show()
			return
		}
		if filenameEntry.Text == "" {
			err := errors.New("Please enter a file name")
			dialog := dialog.NewError(err, w)
			dialog.Show()
			return
		}
		if extension == "" {
			err := errors.New("Please select a file extension")
			dialog := dialog.NewError(err, w)
			dialog.Show()
			return
		}
		err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == extension {
				newName := fmt.Sprintf("%s [%d]%s", filenameEntry.Text, num, filepath.Ext(path))
				newPath := filepath.Join(filepath.Dir(path), newName)
				err = os.Rename(path, newPath)
				if err != nil {
					return err
				}
				num++
			}
			return nil
		})
		if err != nil {
			dialog := dialog.NewError(err, w)
			dialog.Show()
			return
		}
		dialog := dialog.NewInformation("Success", "Files renamed successfully", w)
		dialog.Show()
	})

	content := container.NewVBox(
		container.NewWithoutLayout(
			folderEntry,
			folderButton,
		),
		container.NewHBox(
			widget.NewLabel("Select a file extension:"),
			extensionRadio,
		),
		filenameEntry,
		renameButton,
	)

	w.SetContent(content)
	w.ShowAndRun()
}

// rare thug industries //
// still 187 //
