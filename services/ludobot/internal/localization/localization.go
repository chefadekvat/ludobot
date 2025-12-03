package localization

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Locale string

type Localization struct {
	enKeys         map[string]string
	ruKeys         map[string]string
	pathToMessages string
}

const (
	LocaleRu Locale = "ru"
	LocaleEn Locale = "en"
)

func (l *Localization) addEntry(key string, value string, locale Locale) {
	switch locale {
	case LocaleRu:
		l.ruKeys[key] = value
	case LocaleEn:
		l.enKeys[key] = value
	}
}

func LocalizationWithEntries(pathToMessages string, data []byte) *Localization {
	l := Localization{
		make(map[string]string),
		make(map[string]string),
		pathToMessages,
	}

	var yamlData map[string]interface{}
	err := yaml.Unmarshal(data, &yamlData)
	if err != nil {
		panic(err)
	}

	addKey := func(prefix string, valuesMap map[string]interface{}) {
		parentKey := strings.TrimSuffix(prefix, ".")

		enValue, ok := valuesMap["en"].(string)
		if ok {
			l.addEntry(parentKey, enValue, LocaleEn)
		}

		ruValue, ok := valuesMap["ru"].(string)
		if ok {
			l.addEntry(parentKey, ruValue, LocaleRu)
		}
	}

	var traverse func(string, map[string]interface{})
	traverse = func(prefix string, node map[string]interface{}) {
		for key, value := range node {
			currentKey := prefix + key

			nestedMap, ok := value.(map[string]interface{})
			if key == "values" {
				addKey(prefix, nestedMap)
			} else if ok {
				traverse(currentKey+".", nestedMap)
			}
		}
	}

	messages, ok := yamlData["messages"].(map[string]interface{})
	if ok {
		traverse("", messages)
	}

	return &l
}

func (l *Localization) GetValue(key string, locale Locale) string {
	var value string
	var ok bool

	switch locale {
	case LocaleRu:
		value, ok = l.ruKeys[key]
	case LocaleEn:
		value, ok = l.enKeys[key]
	}

	if ok {
		return value
	}

	return key
}

func (l *Localization) UpdateMessages() error {
	data, err := os.ReadFile(l.pathToMessages)

	if err != nil {
		return err
	}

	newLocalization := LocalizationWithEntries(l.pathToMessages, data)
	l.ruKeys = newLocalization.ruKeys
	l.enKeys = newLocalization.enKeys

	return nil
}

func NewLocalization(pathToMessages string) *Localization {
	data, err := os.ReadFile(pathToMessages)

	if err != nil {
		panic(err)
	}

	return LocalizationWithEntries(pathToMessages, data)
}
