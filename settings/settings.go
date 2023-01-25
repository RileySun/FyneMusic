package settings

import(
	"image/color"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/canvas"
	
	"github.com/RileySun/FyneMusic/utils"
	
	"os"
)

type Config struct {
	Dir string `json:"musicDir"`
	Volume float64 `json:"volume"`
	Setup bool `json:"setup"`
}

var ParentWindow fyne.Window
var dirLocation *widget.Entry

var ReturnToMenu func()
var ChangeVolume func(float64)

func Render() *fyne.Container {
	//Spacer
	spacer := layout.NewSpacer()
	
	//Back Button
	settingsLabel := widget.NewLabel("Settings")
	back := widget.NewButtonWithIcon("", utils.Icons.Cancel, func() {ReturnToMenu()})
	backContainer := container.New(layout.NewHBoxLayout(), settingsLabel, spacer, back)
	topBorder := canvas.NewLine(color.NRGBA{R: 155, G: 155, B: 155, A: 255})
	topBorder.StrokeWidth = 2
	topContainer := container.NewBorder(nil, topBorder, nil, nil, backContainer)
	
	//Music Dir
	dirLabel := widget.NewLabel("Music Directory")
	dirLabel.Alignment = 1 //Center
	dirLocation = widget.NewEntry()
	dirLocation.Text = config.Dir
	button := widget.NewButtonWithIcon("", utils.Icons.Folder, func() {selectMusicDir()})
	dirRow := container.NewBorder(nil, nil, nil, button, dirLocation)
	dirContainer := container.New(layout.NewVBoxLayout(), spacer, dirLabel, dirRow, spacer)
	
	//Master Volume
	volumeLabel := widget.NewLabel("Master Volume")
	volumeLabel.Alignment = 1 //Center
	volume := widget.NewSlider(0, 1)
	volume.Step = 0.1
	volume.OnChanged = func(v float64) {changeVolume(v)}
	volume.Value = config.Volume
	
	//Credit
	creditIMG := canvas.NewImageFromResource(utils.Credit())
	creditIMG.FillMode = canvas.ImageFillOriginal
	creditLabel := widget.NewLabel("Sunshine")
	creditLabel.Alignment = 1
	
	optionsContainer := container.New(layout.NewVBoxLayout(), dirContainer, volumeLabel, volume, spacer, creditIMG, creditLabel)
	paddedContainer := container.New(layout.NewPaddedLayout(), optionsContainer)
	
	saveButton := widget.NewButtonWithIcon("Save", utils.Icons.Save, func() {saveConfig()})
	
	settingsContainer := container.NewBorder(topContainer, saveButton, nil, nil, paddedContainer)
	
	return settingsContainer
}

var config *Config
var mainApp fyne.App

//Util
func LoadSettings(app fyne.App) {
	//App
	mainApp = app
	
	//Get working dir, in case music dir isnt set, this stops the song searcher from checking every file on the computer on first load
	defaultDir, _ := os.Getwd()
	
	//Config
	config = new(Config)
	config.Dir = mainApp.Preferences().StringWithFallback("Dir", defaultDir)
	config.Volume = mainApp.Preferences().FloatWithFallback("Volume", 0.5)
	config.Setup = mainApp.Preferences().BoolWithFallback("Setup", false)
	
}
func GetSettings() *Config {
	return config
}

func saveConfig() {
	if !config.Setup {
		mainApp.Preferences().SetBool("Setup", true)
	}
	
	mainApp.Preferences().SetString("Dir", config.Dir)
	mainApp.Preferences().SetFloat("Volume", config.Volume)
	
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
	mainApp.Preferences().SetFloat("Volume", newVolume) //Do it here also cause people might just click close without saving
	ChangeVolume(newVolume)
}
