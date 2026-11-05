package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
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
	err = os.WriteFile(alvo, originalData, 0o700)
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

func makeTheme(name string, files []string) error {
	home, err := os.UserHomeDir()
	var conf Config
	conf.Itens = make([]Item, len(files))
	dir := filepath.Join(home, "themes", name)
	if err != nil {
		return err
	}
	for i, file := range files {
		if strings.HasPrefix(file, "~") {
			file = untruncate(file, home)
		}
		conf.Itens[i].Alvo = file // saves the path to put the target in
		fileContentB, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		// extrats info
		fileName := filepath.Base(file)
		// makes the directory
		err = os.MkdirAll(dir, 0o700)
		if err != nil {
			return err
		}
		fileStream, err := os.Create(filepath.Join(dir, fileName))
		if err != nil {
			return err
		}
		// writes the config
		_, err = fileStream.Write(fileContentB)
		if err != nil {
			return err
		}
		err = fileStream.Close()
		if err != nil {
			fmt.Printf("Warning error closing %s\n", err.Error())
		}
		// saves the name of the copy that was made
		conf.Itens[i].Original = fileName
	}
	bytes, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	stream, err := os.Create(filepath.Join(dir, "config.json"))
	if err != nil {
		return err
	}
	_, err = stream.Write(bytes)
	if err != nil {
		return err
	}
	err = stream.Close()
	if err != nil {
		fmt.Printf("Warning error closing file %s\n", err.Error())
	}

	return nil
}

func main() {
	if len(os.Args) == 1 {
		listThemes()
		return
	}
	if len(os.Args) < 2 {
		fmt.Printf("Coloque o tema ou make mode\n")
		return
	}
	if os.Args[1] == "-h" {
		fmt.Printf("Ajuda: Uso: themes <tema> para aplicar ou themes para listalos\nTambem -m <nome do novo tema> <quais arquivos para o thema>\n")
		return
	}
	if os.Args[1] == "-m" {
		// 0: programa 1: -m 2: nome do thema 3..: arquuivos
		err := makeTheme(os.Args[2], os.Args[3:])
		if err != nil {
			fmt.Printf("Error creating theme %s\n", err.Error())
		}
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
