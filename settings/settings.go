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
	Setup bool `json:"setup"`
}

var ParentWindow fyne.Window
var config *Config
var dirLocation *widget.Entry

var ReturnToMenu func()
var ChangeVolume func(float64)

func Render() *fyne.Container {
	//Spacer
	spacer := layout.NewSpacer()
	
	//Back Button
	back := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {ReturnToMenu()})
	backContainer := container.New(layout.NewHBoxLayout(), spacer, back)
	
	//Music Dir
	dirLabel := widget.NewLabel("Music Directory")
	dirLabel.Alignment = 1 //Center
	dirLocation = widget.NewEntry()
	dirLocation.Text = config.Dir
	button := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {selectMusicDir()})
	dirRow := container.NewBorder(nil, nil, nil, button, dirLocation)
	dirContainer := container.New(layout.NewVBoxLayout(), spacer, dirLabel, dirRow, spacer)
	
	//Master Volume
	volumeLabel := widget.NewLabel("Master Volume")
	volumeLabel.Alignment = 1 //Center
	volume := widget.NewSlider(0, 1)
	volume.Step = 0.1
	volume.OnChanged = func(v float64) {changeVolume(v)}
	volume.Value = config.Volume
	
	optionsContainer := container.New(layout.NewVBoxLayout(), dirContainer, volumeLabel, volume)
	
	saveButton := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {saveConfig()})
	
	settingsContainer := container.NewBorder(backContainer, saveButton, nil, nil, optionsContainer)
	
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
	
	if !config.Setup {
		config.Setup = true
	}
	
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
	
	ReturnToMenu()
}

//Actions
func selectMusicDir () {
	dialog.ShowFolderOpen(onSelectedDir, ParentWindow)
}

func onSelectedDir(folder fyne.ListableURI, err error) {
	if folder != nil {
		config.Dir = folder.Path()
		dirLocation.Text = config.Dir
		dirLocation.Refresh()
	}
}

func changeVolume(newVolume float64) {
	config.Volume = newVolume
	ChangeVolume(newVolume)
}
