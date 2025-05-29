package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/rivo/tview"
)

type Item struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

var (
	inventory     = []Item{}
	inventoryFile = "inventory.json"
)

func loadInventory() {
	if _, err := os.Stat(inventoryFile); err == nil {
		data, err := os.ReadFile(inventoryFile)
		if err != nil {
			log.Fatal("Error loading inventory file! -", err)
		}
		json.Unmarshal(data, &inventory)
	}
}

func saveInventory() {
	data, err := json.MarshalIndent(inventory, "", " ")
	if err != nil {
		log.Fatal("Error saving inventory! -", err)
	}
	os.WriteFile(inventoryFile, data, 0644)
}

func deleteItem(index int) {
	if index < 0 || index >= len(inventory) {
		fmt.Println("Invalid item index")
		return
	}
	inventory = append(inventory[:index], inventory[index+1:]...)
	saveInventory()
}

func main() {
	app := tview.NewApplication()
	loadInventory()

	inventoryList := tview.NewTextView().SetDynamicColors(true).SetWordWrap(true)
	inventoryList.SetBorder(true).SetTitle("Inventory Items")

	refreshInventory := func() {
		inventoryList.Clear()
		if len(inventory) < 0 {
			fmt.Fprintf(inventoryList, "No items in inventory")
		} else {
			for i, item := range inventory {
				fmt.Fprintf(inventoryList, "[%d] %s (Quantity: %d)\n", i+1, item.Name, item.Quantity)
			}
		}
	}

	itemNameInput := tview.NewInputField().SetLabel("Item name: ")
	itemQuantityInput := tview.NewInputField().SetLabel("Quantity: ")
	itemIdInput := tview.NewInputField().SetLabel("Item to delete: ")

	form := tview.NewForm().AddFormItem(itemNameInput).AddFormItem(itemQuantityInput).AddFormItem(itemIdInput).AddButton("Add Item", func() {
		name := itemNameInput.GetText()
		quantity := itemQuantityInput.GetText()
		if name != "" && quantity != "" {
			quantity, err := strconv.Atoi(quantity)
			if err != nil {
				fmt.Fprintln(inventoryList, "Invalid item quantity")
				return
			}
			inventory = append(inventory, Item{Name: name, Quantity: quantity})
			saveInventory()
			refreshInventory()
			itemNameInput.SetText("")
			itemQuantityInput.SetText("")
		}
	}).AddButton("Delete Item", func() {
		idS := itemIdInput.GetText()
		if idS == "" {
			fmt.Fprintln(inventoryList, "Please enter a item Id to delete")
			return
		}
		id, err := strconv.Atoi(idS)
		if err != nil || id < 0 || id > len(inventory) {
			fmt.Fprintln(inventoryList, "Invalid item Id")
			return
		}
		deleteItem(id - 1)
		fmt.Fprintf(inventoryList, "Item [%d] deleted.\n", id)
		refreshInventory()
		itemIdInput.SetText("")
	}).AddButton("Exit", func() {
		app.Stop()
	})

	form.SetBorder(true).SetTitle("Manage Inventory").SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().AddItem(inventoryList, 0, 1, false).AddItem(form, 0, 1, true)

	refreshInventory()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
