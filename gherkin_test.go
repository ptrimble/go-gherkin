package gherkin_test

import (
	"fmt"
	"github.com/muhqu/go-gherkin"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func parse(t *testing.T, logPrefix, text string) (gherkin.GherkinDOMParser, error) {
	gp := gherkin.NewGherkinDOMParser(text)
	if logPrefix != "" {
		depth := 0
		gp.WithLogFn(func(msg string, args ...interface{}) {
			isBegin := msg[0:5] == "Begin"
			isEnd := msg[0:3] == "End"
			if isEnd {
				depth = depth - 1
			}
			line := fmt.Sprintf(msg, args...)
			depthPrefix := strings.Repeat("  ", depth)
			fmt.Printf("%s%s%s\n", logPrefix, depthPrefix, line)
			if isBegin {
				depth = depth + 1
			}
		})
	}
	gp.Init()
	err := gp.Parse()
	if err == nil {
		gp.Execute()
	}
	return gp, err
}

func mustDomParse(t *testing.T, logPrefix, text string) gherkin.GherkinDOMParser {
	gp, err := parse(t, logPrefix, text)
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}
	return gp
}

func verifyDeadSimpleCalculator(t *testing.T, logPrefix, text string) {
	gp := mustDomParse(t, logPrefix, text)

	feature := gp.Feature()
	assert.NotNil(t, feature)
	assert.Equal(t, "Dead Simple Calculator", feature.Title())
	assert.Equal(t, "Bla Bla\nBla", feature.Description())
	assert.Equal(t, 3, len(feature.Scenarios()), "Number of scenarios")
	assert.Equal(t, 2, len(feature.Tags()), "Number of tags")
	assert.Equal(t, []string{"dead", "simple"}, feature.Tags(), "Feature Tags")

	assert.NotNil(t, feature.Background())
	assert.Equal(t, 1, len(feature.Background().Steps()), "Number of background steps")
	assert.Equal(t, "Given", feature.Background().Steps()[0].StepType())
	assert.Equal(t, "a Simple Calculator", feature.Background().Steps()[0].Text())

	scenario1 := feature.Scenarios()[0]
	assert.NotNil(t, scenario1)
	assert.Equal(t, gherkin.ScenarioNodeType, scenario1.NodeType())
	assert.Equal(t, 1, len(scenario1.Tags()), "Number of tags on Scenario 1")
	assert.Equal(t, []string{"wip"}, scenario1.Tags(), "Tags on Senario 1")
	assert.Equal(t, 5, len(scenario1.Steps()), "Number of steps in Scenario 1")
	assert.Equal(t, "When", scenario1.Steps()[0].StepType())
	assert.Equal(t, "I press the key \"2\"", scenario1.Steps()[0].Text())

	scenario2 := feature.Scenarios()[1]
	assert.NotNil(t, scenario2)
	assert.Equal(t, gherkin.OutlineNodeType, scenario2.NodeType())
	scenario2o, ok := scenario2.(gherkin.OutlineNode)
	assert.True(t, ok)
	assert.Equal(t, 2, len(scenario2.Tags()), "Number of tags on Scenario 2")
	assert.Equal(t, []string{"wip", "expensive"}, scenario2.Tags(), "Tags on Senario 2")
	assert.Equal(t, 5, len(scenario2.Steps()), "Number of steps in Scenario 2")
	assert.Equal(t, "When", scenario2.Steps()[0].StepType())
	assert.Equal(t, "I press the key \"<left>\"", scenario2.Steps()[0].Text())
	assert.Equal(t, [][]string{
		{"left", "operator", "right", "result"},
		{"2", "+", "2", "4"},
		{"3", "+", "4", "7"},
	}, scenario2o.Examples().Table().Rows())

	scenario3 := feature.Scenarios()[2]
	assert.NotNil(t, scenario3)
	assert.Equal(t, gherkin.ScenarioNodeType, scenario3.NodeType())
	assert.Equal(t, 0, len(scenario3.Tags()), "Number of tags on Scenario 3")
	assert.Equal(t, 2, len(scenario3.Steps()), "Number of steps in Scenario 3")
	assert.Equal(t, "When", scenario3.Steps()[0].StepType())
	assert.Equal(t, "I press the following keys:", scenario3.Steps()[0].Text())
	assert.NotNil(t, scenario3.Steps()[0].PyString())
	assert.Equal(t, "  2\n+ 2\n+ 5\n  =\n", scenario3.Steps()[0].PyString().String())
}

func TestParsingRegular(t *testing.T) {
	verifyDeadSimpleCalculator(t, "", `
@dead @simple
Feature: Dead Simple Calculator
  Bla Bla
  Bla

  Background: 
    Given a Simple Calculator
  
  @wip
  Scenario: Adding 2 numbers
     When I press the key "2"
      And I press the key "+"
	  And I press the key "2"
      And I press the key "="
     Then the result should be 4

  @wip @expensive
  Scenario Outline: Simple Math
     When I press the key "<left>"
      And I press the key "<operator>"
	  And I press the key "<right>"
      And I press the key "="
     Then the result should be "<result>"

    Examples:
     | left   | operator | right   | result |
     | 2      | +        | 2       | 4      |
     | 3      | +        | 4       | 7      |

  Scenario: Adding 3 numbers
     When I press the following keys:
     """
       2
     + 2
     + 5
       =
     """
     Then the result should be 9

`)
}

func TestParsingTabAligned(t *testing.T) {
	verifyDeadSimpleCalculator(t, "", `
@dead @simple
Feature: Dead Simple Calculator
	Bla Bla
	Bla

Background: 
	Given a Simple Calculator

@wip
Scenario: Adding 2 numbers
	When I press the key "2"
	And I press the key "+"
	And I press the key "2"
	And I press the key "="
	Then the result should be 4

@wip @expensive
Scenario Outline: Simple Math
	When I press the key "<left>"
	And I press the key "<operator>"
	And I press the key "<right>"
	And I press the key "="
	Then the result should be "<result>"

Examples:
	| left   | operator | right   | result |
	| 2      | +        | 2       | 4      |
	| 3      | +        | 4       | 7      |

Scenario: Adding 3 numbers
	When I press the following keys:
	"""
	  2
	+ 2
	+ 5
	  =
	"""
	Then the result should be 9
`)
}

func TestParsingCondensedAndTrailingWhitespace(t *testing.T) {
	verifyDeadSimpleCalculator(t, "", `@dead @simple Feature: Dead Simple Calculator         
Bla Bla                                 
Bla                                     
Background:                             
Given a Simple Calculator               
@wip Scenario: Adding 2 numbers              
When I press the key "2"                
And I press the key "+"                 
And I press the key "2"                 
And I press the key "="                 
Then the result should be 4             
@wip @expensive Scenario Outline: Simple Math           
When I press the key "<left>"           
And I press the key "<operator>"        
And I press the key "<right>"           
And I press the key "="                 
Then the result should be "<result>"    
Examples:                               
| left | operator | right | result |    
| 2 | + | 2 | 4 |                       
| 3 | + | 4 | 7 |                       
Scenario: Adding 3 numbers              
When I press the following keys:        
"""
  2
+ 2
+ 5
  =
"""
Then the result should be 9`)
}

func TestParsingMinimalNoScenarios(t *testing.T) {
	gp := mustDomParse(t, "", `Feature: Hello World`)
	feature := gp.Feature()
	assert.NotNil(t, feature)
	assert.Equal(t, "Hello World", feature.Title())
}

func TestParsingMinimalNoSteps(t *testing.T) {
	gp := mustDomParse(t, "", `Feature: Hello World
Scenario: Nice people`)
	feature := gp.Feature()
	assert.NotNil(t, feature)
	assert.Equal(t, "Hello World", feature.Title())

	assert.Equal(t, 1, len(feature.Scenarios()), "Number of Scenarios")

	scenario1 := feature.Scenarios()[0]
	assert.NotNil(t, scenario1)
	assert.Equal(t, gherkin.ScenarioNodeType, scenario1.NodeType())
	assert.Equal(t, 0, len(scenario1.Steps()), "Number of steps in Scenario 1")
}

func TestParsingMinimalWithSteps(t *testing.T) {
	gp := mustDomParse(t, "", `
Feature: Hello World

  Scenario: Nice people
    Given a nice person called "Bob"
      And a nice person called "Lisa"
     When "Bob" says to "Lisa": "Hello!"
     Then "Lisa" should reply to "Bob": "Hello!"

`)
	feature := gp.Feature()
	assert.NotNil(t, feature)
	assert.Equal(t, "Hello World", feature.Title())

	assert.Equal(t, 1, len(feature.Scenarios()), "Number of Scenarios")

	scenario1 := feature.Scenarios()[0]
	assert.NotNil(t, scenario1)
	assert.Equal(t, gherkin.ScenarioNodeType, scenario1.NodeType())
	assert.Equal(t, 4, len(scenario1.Steps()), "Number of steps in Scenario 1")
	i := 0
	assert.Equal(t, "Given", scenario1.Steps()[i].StepType())
	assert.Equal(t, `a nice person called "Bob"`, scenario1.Steps()[i].Text())
	i += 1
	assert.Equal(t, "And", scenario1.Steps()[i].StepType())
	assert.Equal(t, `a nice person called "Lisa"`, scenario1.Steps()[i].Text())
	i += 1
	assert.Equal(t, "When", scenario1.Steps()[i].StepType())
	assert.Equal(t, `"Bob" says to "Lisa": "Hello!"`, scenario1.Steps()[i].Text())
	i += 1
	assert.Equal(t, "Then", scenario1.Steps()[i].StepType())
	assert.Equal(t, `"Lisa" should reply to "Bob": "Hello!"`, scenario1.Steps()[i].Text())
}

func TestParsingFailure(t *testing.T) {
	_, err := parse(t, "", `
Feature: Dead Simple Calculator
  Scenario:
    Hurtz
`)
	assert.Error(t, err)
}