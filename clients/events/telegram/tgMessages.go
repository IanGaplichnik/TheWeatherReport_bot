package telegram

const helpMsg = `
1. /help: Prints a help message

2. /checkweather: Get the current weather in your city 🌦️

3. /currentcity: Prints your currently set city

Note:
- Make sure to set your city by sending a location/city name in a message before using /getweather or /checkrain.
- You can change your city at any time by sending a location/city name in a message again.`

const helloMsg = `Hi! I'm your weather bot. I send notifications if it's going to rain in your city. Here's how you can use me:` + helpMsg + `That's it! Stay updated with the weather in your city. ☔️🌤️`

const (
	unknownCommandMsg = "Oh no! This command is not supported! 🤔"
	cityNotSetMsg     = "Uh oh! You need to set the city first ⚠️"
	citySetMsg        = "You've succesfully set a city! 🏙️"
	msgNoCity         = "Couldn't find such city :("
)
