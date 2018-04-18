# Mono Image Resizer & Image Upload Cli

```
A tool to resize images and upload to a site

Usage:
  mono-go-image-cli [command]

Available Commands:
  help        Help about any command
  resize      Resize file or directory
  upload      Upload a image or a directory of images

Flags:
  -h, --help   help for mono-go-image-cli
```

# Resize

To resize a single file to a width of 200px aand save the file as test-1.png
```
$ mono-image resize -f test.png -w 200 -o test-1.png
```
To resize a folder with images to a width of 200px and save them in a folder called test-1
```
$ mono-image resize -d ./test -w 200 -o ./test-1
```

# Upload

To Upload a single file to a website
```
$ ./mono-go-image-cli.exe upload -t TOKEN -s SITEID -u ./test.png
```

To Upload a folder of images to a website
```
$ ./mono-go-image-cli.exe upload -t TOKEN -s SITEID -u ./test/
```