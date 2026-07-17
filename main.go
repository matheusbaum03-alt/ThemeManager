package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

type Item struct {
	Original string `json:"original"`
	Alvo     string `json:"alvo"`
}
type Config struct {
	Itens []Item `json:"itens"`
}

func untruncate(truncatedPath string, home string) string {
	if !strings.HasPrefix(truncatedPath, "~") {
		fmt.Printf("Sua string não começa com '~' string: %s\n", truncatedPath)
		return truncatedPath
	} else {
		incompletePath, _ := strings.CutPrefix(truncatedPath, "~")
		path := path.Join(home, incompletePath)
		return path

	}
}

func copiar(original string, alvo string) {
	originalData, err := os.ReadFile(original)
	if err != nil {
		fmt.Printf("Erro ao ler o original: %s\n", err)
		return
	}
	err = os.WriteFile(alvo, originalData, 0o755)
	if err != nil {
		fmt.Printf("Erro ao escrever a pasta de destino existe?: %s\n", err)
	}
}

func listThemes() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Erro no home dir: %s\n", err)
	}
	pathThemes := path.Join(home, "themes")
	themeDir, err := os.ReadDir(pathThemes)
	if err != nil {
		fmt.Printf("Erro ao abrir diretorio de temas: %s\n", err)
		return
	}
	for i, theme := range themeDir {
		fmt.Printf("%d. %s\n", i+1, theme.Name())
	}
}

func main() {
	if len(os.Args) == 1 {
		listThemes()
		return
	}
	if len(os.Args) > 2 {
		fmt.Printf("Uso: themes <tema> para aplicar ou themes para listalos\n")
		return
	}
	if os.Args[1] == "-h" {
		fmt.Printf("Uso: themes <tema> para aplicar ou themes para listalos\n")
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Erro no home dir: %s\n", err)
		return
	}

	pathDir := path.Join(home, "themes")
	themeDir, err := os.ReadDir(pathDir)
	if err != nil {
		fmt.Printf("erro ao ler themes: %s\n", err)
	}
	var selected string
	for _, theme := range themeDir {
		if !theme.IsDir() {
			fmt.Printf("Entrada %s não é diretorio remova-o\n", theme.Name())
			continue
		}
		if theme.Name() != os.Args[1] {
			continue
		}
		selected = theme.Name()
	}
	if selected == "" {
		fmt.Printf("Erro thema inexesistente ou não achado\nTemas disponiveis:\n")
		listThemes()
		return
	}
	SelThemePath := path.Join(pathDir, selected)
	themeDir, err = os.ReadDir(SelThemePath)
	if err != nil {
		fmt.Printf("Erro ao abir tema selecionado: %s\n", err)
		return
	}
	var configPath string
	for _, theme := range themeDir {
		if theme.Name() == "config.json" {
			configPath = path.Join(SelThemePath, "config.json")
		} else {
			continue
		}
	}
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Erro ao ler .json: %s\n", err)
		return
	}
	var listaItens Config
	err = json.Unmarshal(configFile, &listaItens)
	if err != nil {
		fmt.Printf("Erro ao desfazer .json esta escrito corretamente?: %s\n", err)
		return
	}
	for _, item := range listaItens.Itens {
		originalPath := path.Join(SelThemePath, item.Original)
		pathAlvo := untruncate(item.Alvo, home)
		fmt.Printf("Copiando: %s -> %s\n ", originalPath, pathAlvo)
		copiar(originalPath, pathAlvo)
	}
}
