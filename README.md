# GolangSizer

Open source image resizer coded by kasuraSH

A high quality image resizer built with NASA safety standards that preserves every pixel and prevents stretching

## What it does

Takes your images and resizes them to any dimension you want while keeping the quality intact

No blurry results no stretched images just clean professional resizing using bicubic interpolation

## Quick start

Build it first
go mod tidy go build -o bin/golangresizer.exe ./cmd/golangresizer


Then use it
bin/golangresizer.exe -i input.jpg -o output.png -w 1920 -h 1080


Thats it

## What you can do

Resize photos to any size from 1 pixel to 65535 pixels

Convert between formats like JPEG PNG BMP and TIFF

Upscale images up to 16 times their original size

Downscale images down to one sixteenth of their size

Create thumbnails batch process multiple images whatever you need

## Supported formats

Input works with JPEG PNG BMP TIFF and WebP

Output saves as JPEG PNG BMP or TIFF

Handles both 8 bit and 16 bit color depths

Works with color images and grayscale

## How to use it

Basic resize
bin/golangresizer.exe -i photo.jpg -o resized.png -w 800 -h 600


Make a thumbnail
bin/golangresizer.exe -i large.jpg -o thumb.jpg -w 150 -h 150


Convert formats while resizing
bin/golangresizer.exe -i input.png -o output.jpg -w 1024 -h 768


Get help
bin/golangresizer.exe -help


Check version
bin/golangresizer.exe -version


## Why its different

Built following NASA Power of 10 rules for safety critical code

Uses Mitchell Netravali bicubic interpolation for maximum quality

No pixel loss no stretching no artifacts

Every function is tested every error is handled

Memory usage is fixed and predictable

Code is clean readable and maintainable

## Technical stuff

Written in Go 1.23

Uses bicubic interpolation with a 4x4 pixel kernel

Processes each color channel independently

Clamps values to prevent overflow

Validates everything before processing

## Limits

Minimum size is 1 by 1 pixel

Maximum size is 65535 by 65535 pixels

Maximum file size is 1 gigabyte

Scale factor between one sixteenth and 16 times

## Project structure

Source code is in cmd internal and pkg folders

Main program is in cmd/golangresizer

Core logic is in internal/resizer

Math functions are in internal/interpolation

Input validation is in internal/validator

File operations are in pkg/imageio

## Building from source

Clone the repo
git clone https://github.com/kasurarykerion/GolangSizer.git cd GolangSizer


Get dependencies
go mod tidy


Build it
go build -o bin/golangresizer.exe ./cmd/golangresizer


Run it
bin/golangresizer.exe -help


## License

MIT License

Free to use modify and distribute

See LICENSE file for details

## Author

Created by kasuraSH

Built with NASA Power of 10 safety standards

No compromises on quality or reliability

## Contributing

Pull requests welcome

Follow the existing code style

Keep functions short and focused

Add comments for complex logic

Test your changes before submitting

## Why bicubic interpolation

Bicubic gives better results than bilinear or nearest neighbor

Preserves edges and fine details

Prevents the blocky look you get with simple methods

Industry standard for professional image processing

## Performance

Processes images in linear time based on output size

Memory usage is fixed no dynamic growth during processing

Typical speed is 10 to 50 megapixels per second depending on your CPU

## Safety features

All array access is bounds checked

Integer overflow is prevented

Nil pointers are caught before use

File operations are validated

Resources are properly cleaned up

## Examples

Resize for web
bin/golangresizer.exe -i photo.jpg -o web.jpg -w 1200 -h 800


Create social media image
bin/golangresizer.exe -i original.png -o instagram.jpg -w 1080 -h 1080


Make a desktop wallpaper
bin/golangresizer.exe -i image.jpg -o wallpaper.png -w 2560 -h 1440


Generate a favicon
bin/golangresizer.exe -i logo.png -o favicon.png -w 32 -h 32


## Get started now

Download or clone this repo

Build the executable

Point it at an image

Watch it resize with perfect quality

Thats all there is to it

Simple powerful and reliable
