package main

import (
	"secure_vault/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	secureVaultApp := app.NewWithID("securevaultapp")
	mainWindow := secureVaultApp.NewWindow("Secure Vault Manager")

	// Set windows properities
	mainWindow.Resize(fyne.NewSize(1000, 600))
	mainWindow.CenterOnScreen()
	mainWindow.SetFixedSize(true)

	// Start with the main page
	ui.ShowMainPage(secureVaultApp, mainWindow, "")

	mainWindow.ShowAndRun()
}
