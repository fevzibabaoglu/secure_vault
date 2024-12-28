package ui

import (
	"path/filepath"
	"secure_vault/vault"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ShowVaultDashboard(app fyne.App, window fyne.Window, v *vault.Vault, key []byte, vaultPath string) {
	selectedFileIndex := int64(-1)
	vaultFiles := v.FilesMetadata

	vaultNameLabel := widget.NewLabel("Vault Name: " + filepath.Base(vaultPath))
	vaultCreatedAtLabel := widget.NewLabel("Vault Created At: " + v.Metadata.CreatedAt.Format("2006-01-02 15:04"))

	filesList := widget.NewList(
		func() int {
			return len(vaultFiles)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.ErrorIcon()),
				widget.NewLabel(""),
				widget.NewLabel(""),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			fileName := vaultFiles[id].Name
			addedAt := vaultFiles[id].AddedAt.Format("2006-01-02 15:04")
			fileIntegrity, _ := vault.CheckFileIntegrity(v, int64(id))

			// Set the icon for integrity status
			if fileIntegrity {
				obj.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.ConfirmIcon())
			} else {
				obj.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.CancelIcon())
			}

			// Set the AddedAt value
			obj.(*fyne.Container).Objects[1].(*widget.Label).SetText(addedAt)

			// Set the file name
			obj.(*fyne.Container).Objects[2].(*widget.Label).SetText(fileName)
		},
	)

	filesList.HideSeparators = true

	filesList.OnSelected = func(id widget.ListItemID) {
		selectedFileIndex = vaultFiles[id].Index
	}

	filesList.OnUnselected = func(id widget.ListItemID) {
		selectedFileIndex = -1
	}

	addFileButton := widget.NewButton("Add File", func() {
		dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
			if uri != nil {
				filePath := uri.URI().Path()
				uri.Close()

				// Ask user for file deletion
				dialog.ShowConfirm("Question", "Do you want to delete the original file?",
					func(confirmed bool) {
						var err error
						if confirmed {
							err = vault.AddFileToVault(v, key, filePath, true)
						} else {
							err = vault.AddFileToVault(v, key, filePath, false)
						}

						if err != nil {
							dialog.NewError(err, window).Show()
						}

						vaultFiles = v.FilesMetadata
						filesList.Refresh()
					}, window)
			}
			filesList.UnselectAll()
		}, window).Show()
	})

	extractFileButton := widget.NewButton("Extract File", func() {
		if selectedFileIndex != -1 {
			dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
				if uri != nil {
					folderPath := uri.Path()

					err := vault.ExtractFileFromVault(v, key, selectedFileIndex, folderPath)
					if err != nil {
						dialog.NewError(err, window).Show()
					}

					vaultFiles = v.FilesMetadata
					filesList.Refresh()
				}
				filesList.UnselectAll()
			}, window).Show()
		} else {
			dialog.NewInformation("Error", "Please select a file first.", window).Show()
		}
	})

	removeFileButton := widget.NewButton("Remove File", func() {
		if selectedFileIndex != -1 {
			err := vault.RemoveFileFromVault(v, selectedFileIndex)
			if err != nil {
				dialog.NewError(err, window).Show()
			}
			vaultFiles = v.FilesMetadata
			filesList.UnselectAll()
			filesList.Refresh()
		} else {
			dialog.NewInformation("Error", "Please select a file first.", window).Show()
		}
	})

	saveVaultButton := widget.NewButton("Save Vault", func() {
		err := vault.SaveVault(v, key, vaultPath)
		if err != nil {
			dialog.NewError(err, window).Show()
		}
	})

	backButton := widget.NewButton("Close Vault", func() {
		ShowSelectVaultPage(app, window, filepath.Dir(vaultPath))
	})

	// Content for the top section
	topContent := container.NewVBox(
		vaultNameLabel,
		vaultCreatedAtLabel,
	)

	// Content for the bottom section
	bottomContent := container.NewVBox(
		addFileButton,
		extractFileButton,
		removeFileButton,
		saveVaultButton,
		backButton,
	)

	// Use Border layout to position elements
	content := container.NewBorder(
		topContent,
		bottomContent,
		nil,
		nil,
		filesList,
	)

	window.SetContent(content)
}
