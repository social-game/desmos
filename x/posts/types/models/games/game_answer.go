package games

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ---------------
// --- AnswerID
// ---------------

// AnswerID represents a unique answer id
type AnswerID uint64

// String implements fmt.Stringer
func (id AnswerID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

// MarshalJSON implements Marshaler
func (id AnswerID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

// UnmarshalJSON implements Unmarshaler
func (id *AnswerID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	postID, err := ParseAnswerID(s)
	if err != nil {
		return err
	}

	*id = postID
	return nil
}

// ParseAnswerID returns the AnswerID represented inside the provided
// value, or an error if no id could be parsed properly
func ParseAnswerID(value string) (AnswerID, error) {
	intVal, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return AnswerID(0), err
	}

	return AnswerID(intVal), err
}

// ---------------
// --- GameAnswer
// ---------------

// GameAnswer contains the data of a single Game answer inserted by the creator
type GameAnswer struct {
	ID   AnswerID `json:"id" yaml:"id"`     // Unique id inside the post, serialized as a string for Javascript compatibility
	Text string   `json:"text" yaml:"text"` // Text of the answer
}

// NewGameAnswer returns a new GameAnswer object
func NewGameAnswer(id AnswerID, text string) GameAnswer {
	return GameAnswer{
		ID:   id,
		Text: text,
	}
}

// String implements fmt.Stringer
func (ga GameAnswer) String() string {
	formattedID := strconv.FormatUint(uint64(ga.ID), 10)
	return fmt.Sprintf("Answer - ID: %s ; Text: %s", formattedID, ga.Text)
}

// Validate implements validator
func (ga GameAnswer) Validate() error {
	if strings.TrimSpace(ga.Text) == "" {
		return fmt.Errorf("answer text must be specified and cannot be empty")
	}

	return nil
}

// Equals allows to check whether the contents of p are the same of other
func (ga GameAnswer) Equals(other GameAnswer) bool {
	return ga.ID == other.ID && ga.Text == other.Text
}

// ---------------
// --- GameAnswers
// ---------------

// GameAnswers represents a slice of Game answers
type GameAnswers []GameAnswer

// NewGameAnswers builds a new GameAnswers object starting from the given answers
func NewGameAnswers(answers ...GameAnswer) GameAnswers {
	return answers
}

// Strings implements fmt.Stringer
func (answers GameAnswers) String() string {
	out := "Provided Answers:\n[ID] [Text]\n"
	for _, answer := range answers {
		out += fmt.Sprintf("[%s] [%s]\n",
			strconv.FormatUint(uint64(answer.ID), 10), answer.Text)
	}
	return strings.TrimSpace(out)
}

// Validate implements validator
func (answers GameAnswers) Validate() error {
	if len(answers) < 2 {
		return fmt.Errorf("game answers must be at least two")
	}

	for _, answer := range answers {
		if err := answer.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Equals returns true iff the answers slice contains the same
// data in the same order of the other slice
func (answers GameAnswers) Equals(other GameAnswers) bool {
	if len(answers) != len(other) {
		return false
	}

	for index, answer := range answers {
		if !answer.Equals(other[index]) {
			return false
		}
	}

	return true
}

// AppendIfMissing appends the given answer to the answers slice if it does not exist inside it yet.
// It returns a new slice of GameAnswers containing such GameAnswer.
func (answers GameAnswers) AppendIfMissing(newAnswer GameAnswer) GameAnswers {
	for _, answer := range answers {
		if answer.Equals(newAnswer) {
			return answers
		}
	}
	return append(answers, newAnswer)
}

// ExtractAnswersIDs appends every answer ID to a slice of IDs.
//It returns a slice of answers IDs.
func (answers GameAnswers) ExtractAnswersIDs() (answersIDs []AnswerID) {
	for _, answer := range answers {
		answersIDs = append(answersIDs, answer.ID)
	}
	return answersIDs
}
