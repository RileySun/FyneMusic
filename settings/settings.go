package settings

import(
	"os"
	"io/ioutil"
	"encoding/json"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
)

type Config struct {
	Dir string `json:"musicDir"`
	Volume float64 `json:"volume"`
}

var ParentWindow fyne.Window
var config *Config
var dirLabel *widget.Label

func Render() *fyne.Container {
	dirLabel = widget.NewLabel(config.Dir)
	button := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {selectMusicDir()})
	dirContainer := container.New(layout.NewHBoxLayout(), dirLabel, button)
	
	volume := widget.NewSlider(0, 100)
	volume.OnChanged = func(v float64) {changeVolume(v)}
	volume.Value = config.Volume
	
	saveButton := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {saveConfig()})
	
	settingsContainer := container.New(layout.NewVBoxLayout(), dirContainer, volume, saveButton)
	
	return settingsContainer
}

//Util
func GetSettings() *Config {
	data, fileErr := ioutil.ReadFile("./config.json")
	if fileErr != nil {
		panic("Config Error: " + fileErr.Error())
	}
	
	marshalErr := json.Unmarshal(data, &config)
	if marshalErr != nil {
		panic("Unmarshal Error: " + marshalErr.Error())
	}
	
	return config
}

func saveConfig() {
	data, jsonErr := json.MarshalIndent(config, "", "	")
	
	if jsonErr != nil {
		panic("Config Marshal Error: " + jsonErr.Error())
	}
	
	configFile, fileErr := os.Create("config.json") 
	
	if fileErr != nil {
		panic("Config File Overwrite Error: " + fileErr.Error())
	}
	
	_, saveErr := configFile.Write(data)
	
	if saveErr != nil {
		panic("Config Save Error: " + saveErr.Error())
	}
}

//Actions
func selectMusicDir () {
	dialog.ShowFolderOpen(onSelectedDir, ParentWindow)
}

func onSelectedDir(folder fyne.ListableURI, err error) {
	if folder != nil {
		config.Dir = folder.Path()
		dirLabel.SetText(config.Dir)
		dirLabel.Refresh()
	}
}

func changeVolume(newVolume float64) {
	config.Volume = newVolume
}
