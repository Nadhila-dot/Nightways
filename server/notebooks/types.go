package notebook


type Optional struct {
    Tags        []string `json:"tags,omitempty"`
    Color       string   `json:"color,omitempty"`
    Description string   `json:"description,omitempty"`
}

type Notebook struct {
    ID          int               `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Optional    Optional          `json:"optional,omitempty"`
    Items       map[string]string `json:"items"`
}