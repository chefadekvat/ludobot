package localization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateLocalizationWithEntryFromYaml(t *testing.T) {
	// given
	yamlData := "messages:\n" +
		"  something1:\n" +
		"    something2:\n" +
		"      entry:\n" +
		"        values:\n" +
		"          ru: \"some entry ru value\"\n" +
		"          en: \"some entry en value\"\n"

	expectedRuKeys := map[string]string{
		"something1.something2.entry": "some entry ru value",
	}
	expectedEnKeys := map[string]string{
		"something1.something2.entry": "some entry en value",
	}

	// when
	l := LocalizationWithEntries("some path", []byte(yamlData))

	// then
	assert.Equal(t, expectedRuKeys, l.ruKeys, "Invalid ru keys")
	assert.Equal(t, expectedEnKeys, l.enKeys, "Invalid en keys")
}

func TestShouldCreateLocalizationWithEntriesWithOptionalTranslationsFromYaml(t *testing.T) {
	// given
	yamlData := "messages:\n" +
		"  something1:\n" +
		"    something2:\n" +
		"      entry1:\n" +
		"        values:\n" +
		"          ru: \"some entry1 ru value\"\n" +
		"          en: \"some entry1 en value\"\n" +
		"      entry2:\n" +
		"        values:\n" +
		"          en: \"the only entry2 en value\"\n" +
		"    something3:\n" +
		"      entry3:\n" +
		"        values:\n" +
		"          en: \"some entry3 en value\"\n" +
		"          ru: \"some entry3 ru value\"\n" +
		"    entry4:\n" +
		"      values:\n" +
		"        ru: \"some entry4 ru value\"\n"

	expectedRuKeys := map[string]string{
		"something1.something2.entry1": "some entry1 ru value",
		"something1.something3.entry3": "some entry3 ru value",
		"something1.entry4":            "some entry4 ru value",
	}
	expectedEnKeys := map[string]string{
		"something1.something2.entry1": "some entry1 en value",
		"something1.something2.entry2": "the only entry2 en value",
		"something1.something3.entry3": "some entry3 en value",
	}

	// when
	l := LocalizationWithEntries("some path", []byte(yamlData))

	// then
	assert.Equal(t, expectedRuKeys, l.ruKeys, "Invalid ru keys")
	assert.Equal(t, expectedEnKeys, l.enKeys, "Invalid en keys")
}
