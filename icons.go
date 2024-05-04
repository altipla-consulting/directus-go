package directus

import "encoding/json"

type Icon string

func (icon *Icon) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*icon = Icon(str)
	return nil
}

func (icon *Icon) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*icon))
}
