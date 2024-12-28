package ui

import (
	"fmt"
	"path/filepath"
	"secure_vault/vault"
	"secure_vault/vault/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowPasswordPage(app fyne.App, window fyne.Window, vaultPath string) {
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter password")

	submitButton := widget.NewButton("Submit", func() {
		password := passwordEntry.Text

		integrity, err := vault.CheckVaultIntegrity(vaultPath)
		if err != nil {
			dialog.NewError(err, window)
			return
		}

		if integrity {
			v, err := vault.LoadVault(password, vaultPath)
			if err != nil {
				dialog.NewError(fmt.Errorf("%v: Typed password may be wrong", err), window).Show()
				return
			}

			// Derive encryption key
			key := utils.DeriveKey(password, v.Metadata.Salt)
			ShowVaultDashboard(app, window, v, key, vaultPath)

		} else {
			dialog.ShowConfirm("Error", "Vault hash does not match. Do you wish to continue?",
				func(confirmed bool) {
					if confirmed {
						// Proceed despite the error
						v, err := vault.LoadVault(password, vaultPath)
						if err != nil {
							dialog.NewError(fmt.Errorf("%v: Typed password may be wrong", err), window).Show()
							return
						}

						// Derive encryption key
						key := utils.DeriveKey(password, v.Metadata.Salt)
						ShowVaultDashboard(app, window, v, key, vaultPath)
					}
				}, window)
		}
	})

	backButton := widget.NewButton("Back", func() {
		ShowSelectVaultPage(app, window, filepath.Dir(vaultPath))
	})

	// Content for the center section
	centerContent := container.NewVBox(
		passwordEntry,
	)

	// Content for the bottom section
	bottomContent := container.NewVBox(
		submitButton,
		backButton,
	)

	// Use Border layout to position elements
	content := container.NewBorder(
		widget.NewLabel("Selected Vault File: "+filepath.Base(vaultPath)),
		bottomContent,
		nil,
		nil,
		centerContent,
	)

	window.SetContent(content)
}
