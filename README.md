Photosubset
===========

Intended for quickly and efficiently building subsets of photographs.

![Screenshot](screenshot.png?raw=true)

When I come home from a trip, I build several subsets of photos. The first subset consists of all the photos that I'll keep, the second is the subset I'll show to my family, the third is the subset I'll show to my friends... I used to create these subset painstakingly by copying all the photos and deleting the ones I didn't want.

Then I created photosubset, which makes creating these subsets easy as pie and doesn't take up any additional storage space for any of the subsets created, because it creates them via hardlinks to existing photos.

Simply browse through the photos with left/right arrows and put/delete them from a subset by pressing a number. This means you can create up to 10 subsets (numbered 0-9). They'll be located in folder `tags/tag-<NUM>`. You can append anything you wish to those folder names to make them easier to identify.

*Tested on OSX only, should work on Linux, shouldn't work on Windows since it uses hardlinks. Ping me if you need a Windows version, it should be easy enough to create a workaround.*

Build
-----

Photosubset is written in Go and QML.
First, follow the go-qml installation procedure at https://github.com/go-qml/qml

Then, get the genqrc tool:

    go get gopkg.in/qml.v1/cmd/genqrc

After that, run

    make


