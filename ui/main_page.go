package ui

import (
	"io/fs"
	"secure_vault/ui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowMainPage(app fyne.App, window fyne.Window, selectedFolder string) {
	var vaultFiles []fs.DirEntry
	selectedFolderLabel := widget.NewLabel("No folder selected")

	if selectedFolder != "" {
		vaultFiles, _ = utils.ReadFolderForVaults(selectedFolder)
		selectedFolderLabel.SetText("Selected Folder: " + selectedFolder)
	}

	vaultList := widget.NewList(
		func() int {
			return len(vaultFiles)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Vault File")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(vaultFiles[id].Name())
		},
	)

	vaultList.HideSeparators = true

	vaultList.OnSelected = func(id widget.ListItemID) {
		vaultList.Unselect(id)
	}

	selectFolderButton := widget.NewButton("Select Vault Storage Folder", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				selectedFolder = uri.Path()
				selectedFolderLabel.SetText("Selected Folder: " + selectedFolder)
				vaultFiles, _ = utils.ReadFolderForVaults(selectedFolder)
				vaultList.Refresh()
			}
		}, window).Show()
	})

	nextButton := widget.NewButton("Next", func() {
		if selectedFolder != "" {
			ShowSelectVaultPage(app, window, selectedFolder)
		} else {
			dialog.NewInformation("Error", "Please select a folder first.", window).Show()
		}
	})

	// Content for the top section
	topContent := container.NewVBox(
		selectFolderButton,
		selectedFolderLabel,
	)

	// Use Border layout to position elements
	content := container.NewBorder(
		topContent,
		nextButton,
		nil,
		nil,
		vaultList,
	)

	window.SetContent(content)
}
