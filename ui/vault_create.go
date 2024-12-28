package ui

import (
	"path/filepath"
	uiUtils "secure_vault/ui/utils"
	"secure_vault/vault"
	vaultUtils "secure_vault/vault/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowCreateVaultPage(app fyne.App, window fyne.Window, folderPath string) {
	vaultNameEntry := widget.NewEntry()
	vaultNameEntry.SetPlaceHolder("Enter vault name")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter password")

	confirmButton := widget.NewButton("Create Vault", func() {
		vaultName := vaultNameEntry.Text
		password := passwordEntry.Text

		// Validation
		if vaultName == "" || password == "" {
			dialog.NewInformation("Error", "Both vault name and password are required.", window).Show()
			return
		}

		// Check if a file with the same name already exists in the folder
		vaultFiles, _ := uiUtils.ReadFolderForVaults(folderPath)
		for _, file := range vaultFiles {
			if file.Name() == vaultName+".vault" {
				dialog.NewInformation("Error", "A vault with this name already exists.", window).Show()
				return
			}
		}

		// Full path for the new vault
		vaultPath := filepath.Join(folderPath, vaultName+".vault")

		// Create the vault
		v, err := vault.CreateVault(password)
		if err != nil {
			dialog.NewError(err, window).Show()
			return
		}

		// Derive encryption key
		key := vaultUtils.DeriveKey(password, v.Metadata.Salt)

		err = vault.SaveVault(v, key, vaultPath)
		if err != nil {
			dialog.NewError(err, window).Show()
			return
		}

		ShowSelectVaultPage(app, window, folderPath)
	})

	// Back button to navigate to the previous page
	backButton := widget.NewButton("Back", func() {
		ShowSelectVaultPage(app, window, folderPath)
	})

	// Content for the center section
	centerContent := container.NewVBox(
		vaultNameEntry,
		passwordEntry,
	)

	// Content for the bottom section
	bottomContent := container.NewVBox(
		confirmButton,
		backButton,
	)

	// Use Border layout to position elements
	content := container.NewBorder(
		container.NewCenter(widget.NewLabel("Create a New Vault")),
		bottomContent,
		nil,
		nil,
		centerContent,
	)

	window.SetContent(content)
}
