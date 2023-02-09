package settings_test

import(
	"reflect"
	"testing"
	"github.com/stretchr/testify/assert"
	
	"fyne.io/fyne/v2/app"
	"github.com/RileySun/FyneMusic/settings"
)

func Test_GetSettings(t *testing.T) {
	//New App
	newApp := app.NewWithID("com.sunshine.fynemusic")
	//New Window
	window := newApp.NewWindow("FyneMusic")
	//Load Settings
	settings.LoadSettings(newApp, window)
	//Get Settings
	config := settings.GetSettings()
	//Test
	dirType := reflect.TypeOf(config.Dir).Kind()
	volumeType := reflect.TypeOf(config.Volume).Kind()
	setupType := reflect.TypeOf(config.Setup).Kind()
	logType := reflect.TypeOf(config.Log).Kind()
	
	//Assertions
	assert.True(t, dirType == reflect.String)
	assert.True(t, logType == reflect.String)
	assert.True(t, volumeType == reflect.Float64)
	assert.True(t, setupType == reflect.Bool)
	
	assert.True(t, config.Dir != "")
	assert.True(t, config.Log != "")
}