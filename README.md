# Fyne Music

A pure golang music player app built off the [Fyne Framework](https://fyne.io/ "Fyne Framework").
Currently only supports MP3 format audio files.

------------

### Build Instructions
  
You must have golang installed, see [here](https://go.dev/doc/install).  
  
1. Install Fyne  
Go Version > 1.16:  
`go install fyne.io/fyne/v2/cmd/fyne@latest`  
Go Version < 1.16:  
`go get fyne.io/fyne/v2/cmd/fyne`  


2. Build App  
macOS:-  
`fyne package -os darwin -icon App_Icon.png`  
windows:  
`fyne package -os windows -icon App_Icon.png`  
linux:  
`fyne package -os linux -icon App_Icon.png`  
  
3. Enjoy  
:3  
