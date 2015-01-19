# bgmysword
Program to generate MySword Bible modules by looking up non-copyrighted Bibles on biblegateway.com.  
(Written in Go/Golang)  
Emphasis on **NON**-copyrighted  
Please do not distribute any modules created by this program.

I'll show you how to install the program. But first, you need to have `Go` (Google's programming language) installed and a couple other `Go` modules installed (it's easy; don't worry).  
##Instructions to Install Go  
###Windows  
1. Go to the `Go` download page: [https://golang.org/dl/](https://golang.org/dl/)  
2. Download the `.msi` file (e.g. `go1.4.1.windows-amd64.msi`).
3. Double-click the `.msi` to install Go.  

###Ubuntu
1. Run `sudo apt-get install golang`  

###Arch Linux  
1. Run `sudo pacman -S go`

##Instructions to Install Dependencies  
###Windows  
1. Open `cmd.exe`. (Either by typing 'cmd.exe' on the Start menu's searchbar or by following the instructions [here](http://windows.microsoft.com/en-us/windows-vista/open-a-command-prompt-window)).  
2. Run these commands:  
    `go get github.com/PuerkitoBio/goquery`  
    `go get github.com/fatih/color`  
    `go get github.com/mattn/go-sqlite3`  

3. Keep the black `cmd.exe` window open and go to the next set of instructions.  

###Linux (Ubuntu, Arch, etc.)
1. Open `Terminal` (ctrl+alt+t).
2. Run these commands:  
    `go get github.com/PuerkitoBio/goquery`  
    `go get github.com/fatih/color`  
    `go get github.com/mattn/go-sqlite3`  

3. Keep Terminal open.  

##Instructions to Install bgmysword  
###Windows  
1. With the black `cmd.exe` window open, run this command:  
    `go get github.com/GreenRaccoon23/bgmysword`  

2. This *should* install it automatically (I haven't been able to test it on another computer yet). Run it like this (replacing KJV with the translation you want):  
    `bgmysword KJV`  

###Linux (Ubuntu, Arch, etc.)  
1. With Terminal open, run this command:  
    `go get github.com/GreenRaccoon23/bgmysword`  

2. This *should* install it automatically (I haven't been able to test it on another computer yet). Run it like this (replacing KJV with the translation you want):  
    `bgmysword KJV`

