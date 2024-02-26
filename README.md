# APOD Wallpaper
<strong>A simple Go program that uses Nasa's Astronomy Picture of the Day to set the desktop wallpaper to todays picture.</strong>

It works by contacting the API and finding the picture link in the response. It then takes that link and downloads the picture into the project directory. After that is done it moves any old pictures into the archived folder and sets todays picture as the wallpaper.
<br><br>
## How to use
Go to the link and generate your API key<br>
https://api.nasa.gov/<br><br>
Create a file named ``.env`` in the project directory<br><br>
Write this in the file and enter your API key instead of ``yourapikeyhere``<br>
```
API_KEY=yourapikeyhere
```
Now you can either run the ``main.go`` file from the terminal or build the executable and be on your merry way!

<br>
<strong>Currently only works on Windows because of the way the wallpaper is set.</strong>
