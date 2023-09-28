### Megaminx 

This repository is a GUI and solver for the Megaminx, the 3x3x3 Rubik's Cube. 
This project was an assignment for CS463 at the University of Kentucky.

#### Build instructions
1. If you don't have it, install Go >=1.15
2. This project uses ebitengine as the interface between the graphics card and the code, you will need to 
do different things depending on your operating system. If you're on Windows, you don't need to do anything really.
Otherwise, follow this guide: https://ebitengine.org/en/documents/install.html
3. Run `go get` to install all dependencies for this project
4. `go build .` to build the executable
5. `./megaminx` to run it.
