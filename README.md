# Fyne Music

A pure golang music player app built off the [Fyne Framework](https://fyne.io/ "Fyne Framework").
Currently only supports MP3 format audio files.

------------
### Install Instructions  

FyneMusic has pre-built application releases for the following platforms  

MacOS: [Download](https://github.com/RileySun/FyneMusic/releases/download/v1.0/FyneMusic.macOS.dmg)  
Windows: [Download](https://github.com/RileySun/FyneMusic/releases/download/v1.0/FyneMusic.Windows.zip)  
Linux: Please build using instructions below  

------------

### Build Instructions
  
You must have golang installed, see [here](https://go.dev/doc/install).  
  
1. Install Fyne  
Go Version > 1.16:  
`go install fyne.io/fyne/v2/cmd/fyne@latest`  
Go Version < 1.16:  
`go get fyne.io/fyne/v2/cmd/fyne`  


2. Build App (for cross compiling see [here](https://developer.fyne.io/started/cross-compiling))
macOS:   
`fyne package -os darwin -icon App_Icon.png`  
windows:  
`fyne package -os windows -icon App_Icon.png`  
linux:  
`fyne package -os linux -icon App_Icon.png`  
  
3. Enjoy  
:3  
