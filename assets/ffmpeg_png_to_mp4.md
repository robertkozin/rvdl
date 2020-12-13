https://trac.ffmpeg.org/wiki/Slideshow#Singleimage

    ffmpeg -loop 1 -i input.png -c:v libx264 -t 1 -pix_fmt yuv420p output.mp4