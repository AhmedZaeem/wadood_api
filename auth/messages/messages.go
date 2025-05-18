package messages

func Get(key, lang string) string {
    switch lang {
    case "ar":
        if msg, ok := ArMessages[key]; ok {
            return msg
        }
    default:
        if msg, ok := EnMessages[key]; ok {
            return msg
        }
    }
    return "Message not found."
}