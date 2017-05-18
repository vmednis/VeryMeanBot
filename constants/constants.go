//
// The version of this file uploaded to github is heavily censored to preserve some common decency
// I only wanted to share my code not my insult crafting skills
//

package constants

//
// Version the version of bot
//
const Version = "dev-0.2"

//
// BotPath the directory where all the bots files are stored
//
const BotPath = "/home/valters/.config/meanbot/" //Beta
//const BotPath = "/home/meanbot/.meanbot/" //Production

//
// BotUserName bot username in lower-case
//
const BotUserName = "betaverymeanbot" //Beta
//const BotUserName = "verymeanbot" //Production

//
// BotAPIKey Telegram APIKey
//
const BotAPIKey = "censored" //Beta
//const BotAPIKey = "censored" //Production

//
// Strings bot uses in replies that aren't insults
//

// ReplyNotEnoughArguments reply for when not enough arugments are given for a command
const ReplyNotEnoughArguments = "censored"

// ReplyNotANumber reply for when a numeric argument %v[0] is not a number
const ReplyNotANumber = "censored %v"

// ReplyNotASetting reply for when setting %v[0] doesn't exist
const ReplyNotASetting = "censored %v %v"

// ReplySettingSet reply for when a setting is %v[0] is set to %v[1] successfuly
const ReplySettingSet = "censored %v %v"

// ReplyJoinChat message sent when bot gets added to a new chat where %v[0] is the type of the chat
const ReplyJoinChat = "censored %v"

// ReplyNonexistantCommand reply for when a user has sent a non existant command %v[0] targeted at the bot
const ReplyNonexistantCommand = "censored %v"

// ReplyAcceptSuggestion reply for when a suggestion is accepted
const ReplyAcceptSuggestion = "censored"

// ReplyMissingSuggestion reply for when /suggest command is missing argument "suggestion"
const ReplyMissingSuggestion = "censored"

// ReplyInsultItself reply for when a user tries to /insult the bot itself
const ReplyInsultItself = "censored"

// ReplyMissingTarget reply for when /insult command is missing argument "target"
const ReplyMissingTarget = "censored"
