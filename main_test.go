package main

import "testing"
import "encoding/json"

// by this command: "go test *.go"
func TestThrow(t *testing.T){

	text := []byte(`
	{
		"_links": { "self": { "href": "https://MY_URL" } },
		"arena": { 
			"dims": [4,3],
			"state": {
				"https://MY_URL": { "x": 2, "y": 2, "wasHit": false, "direction":  "N"},
				"https://A_PLAYERS_URL": { "x": 2, "y": 1}
			}
		}
	}
	`)
	var arena StateUpdate
	err := json.Unmarshal(text, &arena)
    if err != nil {
        panic(err)
    }
	result := decisionTree(arena)
	expected := "T"

    if result != expected {
        t.Errorf("got %q but expecting %q", result, expected)
    }
}

func TestRight(t *testing.T){

	text := []byte(`
	{
		"_links": { "self": { "href": "https://MY_URL" } },
		"arena": { 
			"dims": [3,2], 
			"state": {
				"https://MY_URL": { "x": 2, "y": 1, "wasHit": false, "direction":  "S"}
			}
		}
	}
	`)
	var arena StateUpdate
	err := json.Unmarshal(text, &arena)
    if err != nil {
        panic(err)
    }
	result := decisionTree(arena)
	expected := "R"

    if result != expected {
        t.Errorf("got %q but expecting %q", result, expected)
    }
}

func TestLeft(t *testing.T){

	text := []byte(`
	{
		"_links": { "self": { "href": "https://MY_URL" } },
		"arena": { 
			"dims": [3,2], 
			"state": {
				"https://MY_URL": { "x": 2, "y": 1, "wasHit": false, "direction":  "E"}
			}
		}
	}
	`)
	var arena StateUpdate
	err := json.Unmarshal(text, &arena)
    if err != nil {
        panic(err)
    }
	result := decisionTree(arena)
	expected := "L"

    if result != expected {
        t.Errorf("got %q but expecting %q", result, expected)
    }
}

func TestNotThrow(t *testing.T){

	text := []byte(`
	{
		"test": "dont throw north",
		"_links": { "self": { "href": "https://MY_URL" } },
		"arena": { 
			"dims": [4,3],
			"state": {
				"https://MY_URL": { "x": 1, "y": 1, "wasHit": false, "direction":  "N"},
				"https://A_PLAYERS_URL": { "x": 2, "y": 1, "wasHit": false, "direction":  "N"}
			}
		}
	}
	`)
	var arena StateUpdate
	err := json.Unmarshal(text, &arena)
    if err != nil {
        panic(err)
    }
	result := decisionTree(arena)
	expected := "T"

    if result == expected {
        t.Errorf("got %q but it is not expected", result)
    }
}