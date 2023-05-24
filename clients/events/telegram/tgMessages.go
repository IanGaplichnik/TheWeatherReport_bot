package telegram

const helpMsg = `
1. /setcity <city_name>: Set your city for weather updates 🏙️.
Example: /setcity London

2. /getweather: Get the current weather in your city 🌦️.
Example: /getweather

Note:
- Make sure to set your city using /setcity before using /getweather or /checkrain.
- You can change your city at any time by using /setcity again.`

const helloMsg = `Hi! I'm your weather bot. I send notifications if it's going to rain in your city. Here's how you can use me:` + helpMsg + `That's it! Stay updated with the weather in your city. ☔️🌤️`

const (
	unknownCommandMsg = "Oh no! This command is not supported! 🤔"
	cityNotSetMsg     = "Uh oh! You need to set the city first :)"
	citySetMsg        = "You've succesfully set a city! 🏙️"
)
