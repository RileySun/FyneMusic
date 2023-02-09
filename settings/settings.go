package settings

import(
	"os"
	"runtime"
	"image/color"
	"path/filepath"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/canvas"
	
	"github.com/RileySun/FyneMusic/utils"
)

type Config struct {
	Dir string `json:"musicDir"`
	Volume float64 `json:"volume"`
	Setup bool `json:"setup"`
	Log string `json:"logDir"`
}

var ParentWindow fyne.Window
var dirLocation *widget.Entry

var ReturnToMenu func()
var ChangeVolume func(float64)

var config *Config
var mainApp fyne.App 		//Needed for Prefrences api, weird fyne design
var mainWindow fyne.Window	//Need for copying to clipboard

//Init
func LoadSettings(app fyne.App, window fyne.Window) {
	//App
	mainApp = app
	mainWindow = window
	
	//Platform specific log directories and default music dir
	var logDir string
	homeDir, _ := os.UserHomeDir()
	switch runtime.GOOS {
		case "windows":
			logDir = filepath.Join(filepath.Join(filepath.Join(homeDir, "AppData"), "Roaming"), "com.sunshine.fynemusic")
			_, err := os.Stat(logDir)
			if os.IsNotExist(err) {
				dirErr := os.Mkdir(logDir, 0755)
				if dirErr != nil {
					panic("settings: log file - " + dirErr.Error())
				}
			}
		case "darwin":			
			logDir = filepath.Join(filepath.Join(filepath.Join(homeDir, "Library"), "Application Support"), "com.sunshine.fynemusic")
			_, statErr := os.Stat(logDir)
			if os.IsNotExist(statErr) {
				dirErr := os.Mkdir(logDir, 0755)
				if dirErr != nil {
					panic("settings: " + dirErr.Error())
				}
			}
		case "linux":
			logDir = "/var/log/com.sunshine.fynemusic"
			_, statErr := os.Stat(logDir)
			if os.IsNotExist(statErr) {
				dirErr := os.Mkdir(logDir, 0755)
				if dirErr != nil {
					panic("settings: log file - " + dirErr.Error())
				}
			}
	}
	homeDir = filepath.Join(homeDir, "Music")
	
	//Config
	config = new(Config)
	config.Dir = mainApp.Preferences().StringWithFallback("Dir", homeDir)
	config.Volume = mainApp.Preferences().FloatWithFallback("Volume", 0.5)
	config.Setup = mainApp.Preferences().BoolWithFallback("Setup", false)
	config.Log = mainApp.Preferences().StringWithFallback("Log", logDir)
	
	utils.SetLogPath(config.Log)
	
}

//Render
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
	dirContainer := container.NewBorder(nil, nil, dirLabel, nil, dirRow)
	
	//Master Volume
	volumeLabel := widget.NewLabel("Master Volume")
	volumeLabel.Alignment = 1 //Center
	volume := widget.NewSlider(0, 1)
	volume.Step = 0.1
	volume.OnChanged = func(v float64) {changeVolume(v)}
	volume.Value = config.Volume
	
	//Log Dir
	logLabel := widget.NewLabel("Log Path: " + config.Log)
	logLabel.Alignment = 1 //Center
	logLabel.Wrapping = 1
	logButton := widget.NewButtonWithIcon("", utils.Icons.Copy, func() {copyLogDir()})
	logContainer := container.NewBorder(nil, nil, nil, logButton, logLabel)
	
	//Credit
	creditIMG := canvas.NewImageFromResource(utils.Credit())
	creditIMG.FillMode = canvas.ImageFillOriginal
	creditLabel := widget.NewLabel("Sunshine")
	creditLabel.Alignment = 1
	
	optionsContainer := container.New(layout.NewVBoxLayout(), dirContainer, volumeLabel, volume, spacer, logContainer, creditIMG, creditLabel)
	paddedContainer := container.New(layout.NewPaddedLayout(), optionsContainer)
	
	saveButton := widget.NewButtonWithIcon("Save", utils.Icons.Save, func() {saveConfig()})
	
	settingsContainer := container.NewBorder(topContainer, saveButton, nil, nil, paddedContainer)
	
	return settingsContainer
}

//Util
func GetSettings() *Config {
	return config
}

func saveConfig() {
	if !config.Setup {
		mainApp.Preferences().SetBool("Setup", true)
		mainApp.Preferences().SetString("Log", config.Log)
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

func copyLogDir() {
	clip := mainWindow.Clipboard()
	clip.SetContent(config.Log)
}
