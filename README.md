# Getting it running on Linux 


## OS: Linux Mint 21 Vanessa, based on Ubuntu 22.04

Install libopencv-dev

		$ sudo apt install -y libopencv-dev

Check the version

		$ apt-cache info libopencv-dev | grep -i Version

for me it was 4.5.4

## OS: Fedora 35

Install libopencv-dev

		$ sudo dnf install opencv-devel

## Version management of opencv and gocv

Check https://github.com/hybridgroup/gocv/releases for a compatible version, for me it was 0.29.0.
Modify go.mod to use that version
		
		require gocv.io/x/gocv v0.29.0

Ensure this version is downloaded and checksums are in place in go.sum

		$ go mod download gocv.io/x/gocv

You are now (hopefully) able to compile and run all examples.

# Getting it running using Docker
Unsolved: I am still unable to properly display a video window, general processing works

Make sure that the docker image is there

		$ sudo docker pull gocv/opencv:4.7.0
		car-count-workshop/$ sudo docker run -it -v $PWD/:/car gocv/opencv:4.7.0 bash
		(in docker)$ cd /car/
		
This passed the first video device into docker to allow use of it

		$ sudo docker run -it --device=/dev/video0 -v $PWD/:/car gocv/opencv:4.7.0 bash
