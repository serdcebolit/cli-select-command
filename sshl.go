package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/manifoldco/promptui"
)

type Cmd struct {
	Name        string   `json:"name"`
	ProgramName string   `json:"program"`
	Args        []string `json:"args"`
}

type Config struct {
	Cmds []Cmd `json:"commands"`
}

func main() {
	config := getDataFromJsonFile()

	var names []string
	for _, cmd := range config.Cmds {
		names = append(names, cmd.Name)
	}

	prompt := promptui.Select{
		Label: "Выберите сервер",
		Items: names,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Отмена\n")
		return
	}

	for _, cmd := range config.Cmds {
		if cmd.Name == result {

			cmd := exec.Command(cmd.ProgramName, cmd.Args...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println("Ошибка запуска команды: ", err)
				return
			}

			return
		}
	}
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return dir
}

func getDataFromJsonFile() Config {
	var cmds Config
	var configPath = filepath.Join(getCurrentDirectory(), "config.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("config.json не найден")
		return cmds
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Ошибка: ", err)
	}

	if err := json.NewDecoder(configFile).Decode(&cmds); err != nil {
		fmt.Println("Ошибка: ", err)
	}

	return cmds
}
