package main

import "testing"

func TestParseConfiguration(t *testing.T) {
	config, err := parseConfiguration([]byte(`
Log = true
GeneratePlatformUserId = false
Authentication = 'cached'

[Games]
Enabled = ['age2']

[Games.age2]
Hosts = ['192.168.1.10']

[Announcement]
Enabled = true
Multicast = true
Port = 32000
MulticastGroup = '239.31.97.8'
`))
	if err != nil {
		t.Fatalf("parseConfiguration returned an error: %v", err)
	}
	if !config.Log || config.Authentication != "cached" || config.Announcement.Port != 32000 {
		t.Fatalf("unexpected parsed configuration: %+v", config)
	}
	if len(config.Games.Enabled) != 1 || config.Games.Enabled[0] != "age2" {
		t.Fatalf("unexpected enabled games: %#v", config.Games.Enabled)
	}
	if len(config.Games.Age2.Hosts) != 1 || config.Games.Age2.Hosts[0] != "192.168.1.10" {
		t.Fatalf("unexpected hosts: %#v", config.Games.Age2.Hosts)
	}
}

func TestParseConfigurationRejectsInvalidTOML(t *testing.T) {
	if _, err := parseConfiguration([]byte("[Announcement\nPort = 32000")); err == nil {
		t.Fatal("parseConfiguration accepted invalid TOML")
	}
}

func TestParseConfigurationRejectsInvalidType(t *testing.T) {
	if _, err := parseConfiguration([]byte("[Announcement]\nPort = 'not-a-port'")); err == nil {
		t.Fatal("parseConfiguration accepted an invalid port type")
	}
}

func TestValidateConfiguration(t *testing.T) {
	cases := []struct {
		name      string
		config    Configuration
		field     string
		hasErrors bool
	}{
		{
			name: "valid",
			config: Configuration{
				Authentication: "disabled",
				Games:          GamesConfiguration{Enabled: []string{"age2"}},
				Announcement:   AnnouncementConfig{Port: 31978},
			},
		},
		{
			name: "unknown authentication",
			config: Configuration{
				Authentication: "unknown",
				Games:          GamesConfiguration{Enabled: []string{"age2"}},
				Announcement:   AnnouncementConfig{Port: 31978},
			},
			field:     "Authentication",
			hasErrors: true,
		},
		{
			name: "empty games and invalid port",
			config: Configuration{
				Authentication: "disabled",
				Announcement:   AnnouncementConfig{Port: 65536},
			},
			field:     "Games.Enabled",
			hasErrors: true,
		},
		{
			name: "unknown game",
			config: Configuration{
				Authentication: "disabled",
				Games:          GamesConfiguration{Enabled: []string{"age5"}},
				Announcement:   AnnouncementConfig{Port: 31978},
			},
			field:     "Games.Enabled.0",
			hasErrors: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			errors := validateConfiguration(&testCase.config)
			if !testCase.hasErrors && len(errors) != 0 {
				t.Fatalf("validateConfiguration returned unexpected errors: %#v", errors)
			}
			if testCase.hasErrors && len(errors) == 0 {
				t.Fatal("validateConfiguration accepted invalid configuration")
			}
			if testCase.field != "" && errors[0].Field != testCase.field {
				t.Fatalf("first error field = %q, want %q", errors[0].Field, testCase.field)
			}
		})
	}
}
