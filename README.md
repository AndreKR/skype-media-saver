# Skype Media Saver

## What is this?

There is a Skype setting called "When I receive a file..." where you can select a folder and Skype used to put all
received files into that directory.
 
At some point Skype stopped doing that for images and videos. Those are now **only in the cloud** (and I don't know
for how long) unless you explicitly save each picture manually.

So this little tool will watch the Skype media cache folder and whenever a new file is received, it will copy this file
to its output folder.

## How to use?

Download the [latest release](). Put it in a folder and make sure it can write to this folder. Put a link in the Start >
All Programs > Startup folder. Run it.

It will write every received picture to that folder.

Currently it only works on Windows because I don't know where Linux Skype puts its files. PRs welcome.

## How to exit?

Run it again, it will ask if you want it to exit.

## How to build?

```
go get
go build -ldflags "-s -w -H=windowsgui"
```

(`-s -w` for smaller executables, `-H=windowsgui` for no console window) 

### What is `rsrc.syso`?

It will be linked by the go compiler to [tell windows to use modern windows buttons](https://msdn.microsoft.com/en-us/library/windows/desktop/bb773175(v=vs.85).aspx)
in the dialogs. It can be generated from `beautiful_buttons.manifest.xml` with
`rsrc -manifest=beautiful_dialogs.manifest.xml` (by [rsrc]()). 
