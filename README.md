# APOD Wallpaper
<strong>A simple Go program that uses Nasa's Astronomy Picture of the Day to set the desktop wallpaper to todays picture.</strong>

It works by contacting the API and finding the picture link in the response. It then takes that link and downloads the picture into the project directory. After that is done it moves any old pictures into the archived folder and sets todays picture as the wallpaper.
<br>

<strong>Currently only works on Windows because of the way the wallpaper is set.</strong>
