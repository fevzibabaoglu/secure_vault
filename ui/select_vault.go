package ui

import (
	"path/filepath"
	"secure_vault/ui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowSelectVaultPage(app fyne.App, window fyne.Window, folderPath string) {
	var selectedFile string
	vaultFiles, _ := utils.ReadFolderForVaults(folderPath)

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

	vaultList.OnSelected = func(id widget.ListItemID) {
		selectedFile = vaultFiles[id].Name()
	}

	createVaultButton := widget.NewButton("Create a New Vault", func() {
		ShowCreateVaultPage(app, window, folderPath)
	})

	selectVaultButton := widget.NewButton("Select Vault", func() {
		if selectedFile != "" {
			ShowPasswordPage(app, window, filepath.Join(folderPath, selectedFile))
		} else {
			dialog.NewInformation("Error", "Please select a vault file first.", window).Show()
		}
	})

	backButton := widget.NewButton("Back", func() {
		ShowMainPage(app, window, folderPath)
	})

	// Content for the top section
	topContent := container.NewVBox(
		createVaultButton,
		widget.NewLabel("Selected Folder: "+folderPath),
	)

	// Content for the bottom section
	bottomContent := container.NewVBox(
		selectVaultButton,
		backButton,
	)

	// Use Border layout to position elements
	content := container.NewBorder(
		topContent,
		bottomContent,
		nil,
		nil,
		vaultList,
	)

	window.SetContent(content)
}
