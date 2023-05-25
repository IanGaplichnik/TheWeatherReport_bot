package telegram

const helpMsg = `
1. /getweather: Get the current weather in your city 🌦️

2. /help: Prints a help message

3. /checkrain: Get a short message, stating if it's going to rain

4. /currentcity: Prints your currently set city

Note:
- Make sure to set your city using /setcity before using /getweather or /checkrain.
- You can change your city at any time by using /setcity again.`

const helloMsg = `Hi! I'm your weather bot. I send notifications if it's going to rain in your city. Here's how you can use me:` + helpMsg + `That's it! Stay updated with the weather in your city. ☔️🌤️`

const (
	unknownCommandMsg = "Oh no! This command is not supported! 🤔"
	cityNotSetMsg     = "Uh oh! You need to set the city first :)"
	citySetMsg        = "You've succesfully set a city! 🏙️"
	msgNoCity         = "Couldn't find such city :("
)
