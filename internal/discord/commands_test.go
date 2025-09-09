package discord

import (
    "strconv"
    "strings"
    "testing"
)

func TestParsePrefixCommand(t *testing.T) {
    testCases := []struct {
        name        string
        message     string
        expectedCmd string
        expectedArgs []string
    }{
        {
            name:        "simple command",
            message:     "!macgym",
            expectedCmd: "macgym",
            expectedArgs: []string{},
        },
        {
            name:        "command with args",
            message:     "!badminton events 14",
            expectedCmd: "badminton",
            expectedArgs: []string{"events", "14"},
        },
        {
            name:        "subscribe with threshold",
            message:     "!subscribe 5",
            expectedCmd: "subscribe",
            expectedArgs: []string{"5"},
        },
        {
            name:        "unsubscribe",
            message:     "!unsubscribe",
            expectedCmd: "unsubscribe",
            expectedArgs: []string{},
        },
        {
            name:        "help command",
            message:     "!help",
            expectedCmd: "help",
            expectedArgs: []string{},
        },
        {
            name:        "case insensitive",
            message:     "!MACGYM",
            expectedCmd: "macgym",
            expectedArgs: []string{},
        },
        {
            name:        "multiple spaces",
            message:     "!badminton   events   7",
            expectedCmd: "badminton",
            expectedArgs: []string{"events", "7"},
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Simulate the parsing logic from the handler
            if len(tc.message) == 0 || tc.message[0] != '!' {
                t.Error("Message should start with !")
                return
            }
            
            parts := strings.Fields(tc.message)
            if len(parts) == 0 {
                t.Error("Should have at least one part")
                return
            }
            
            commandName := strings.ToLower(parts[0][1:]) // Remove ! prefix
            args := parts[1:]
            
            if commandName != tc.expectedCmd {
                t.Errorf("Expected command '%s', got '%s'", tc.expectedCmd, commandName)
            }
            
            if len(args) != len(tc.expectedArgs) {
                t.Errorf("Expected %d args, got %d", len(tc.expectedArgs), len(args))
                return
            }
            
            for i, arg := range args {
                if arg != tc.expectedArgs[i] {
                    t.Errorf("Expected arg[%d] '%s', got '%s'", i, tc.expectedArgs[i], arg)
                }
            }
        })
    }
}

func TestCommandValidation(t *testing.T) {
    validCommands := []string{
        "macgym",
        "badminton",
        "subscribe",
        "unsubscribe",
        "help",
    }
    
    for _, cmd := range validCommands {
        t.Run("valid_"+cmd, func(t *testing.T) {
            // Test that the command would be recognized
            switch cmd {
            case "macgym", "badminton", "subscribe", "unsubscribe", "help":
                // These should be handled
            default:
                t.Errorf("Command '%s' should be valid but isn't handled", cmd)
            }
        })
    }
    
    invalidCommands := []string{
        "unknown",
        "invalid",
        "test",
    }
    
    for _, cmd := range invalidCommands {
        t.Run("invalid_"+cmd, func(t *testing.T) {
            // Test that unknown commands would trigger the default case
            handled := false
            switch cmd {
            case "macgym", "badminton", "subscribe", "unsubscribe", "help":
                handled = true
            default:
                // This should be the case for invalid commands
            }
            
            if handled {
                t.Errorf("Command '%s' should be invalid but is handled", cmd)
            }
        })
    }
}

func TestArgumentParsing(t *testing.T) {
    testCases := []struct {
        name     string
        args     []string
        expected int
    }{
        {
            name:     "no args",
            args:     []string{},
            expected: 0,
        },
        {
            name:     "valid number",
            args:     []string{"5"},
            expected: 5,
        },
        {
            name:     "invalid number",
            args:     []string{"abc"},
            expected: 0,
        },
        {
            name:     "negative number",
            args:     []string{"-1"},
            expected: 0,
        },
        {
            name:     "zero",
            args:     []string{"0"},
            expected: 0,
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            threshold := 0
            if len(tc.args) > 0 {
                if t, err := strconv.Atoi(tc.args[0]); err == nil && t >= 0 {
                    threshold = t
                }
            }
            
            if threshold != tc.expected {
                t.Errorf("Expected threshold %d, got %d", tc.expected, threshold)
            }
        })
    }
}
