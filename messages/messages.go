package messages

func Get(key, lang string) string {
    switch lang {
    case "ar":
        if msg, ok := AR[key]; ok {
            return msg
        }
    default:
        if msg, ok := EN[key]; ok {
            return msg
        }
    }
    return key
}