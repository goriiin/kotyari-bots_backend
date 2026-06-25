package posts

import (
	"testing"

	"github.com/google/uuid"
)

// TestKafkaUpdatePostRequestToModel guards the IDOR fix: the UserID must be
// carried from the Kafka command into the model so the repository can scope the
// UPDATE by user_id (otherwise any user could edit any post).
func TestKafkaUpdatePostRequestToModel(t *testing.T) {
	postID := uuid.New()
	userID := uuid.New()

	req := KafkaUpdatePostRequest{
		PostID: postID,
		UserID: userID,
		Title:  "title",
		Text:   "text",
	}

	m := req.ToModel()
	if m.ID != postID {
		t.Fatalf("ID = %v, want %v", m.ID, postID)
	}
	if m.UserID != userID {
		t.Fatalf("UserID = %v, want %v (IDOR scoping would be lost)", m.UserID, userID)
	}
	if m.Title != "title" || m.Text != "text" {
		t.Fatalf("title/text not carried: %+v", m)
	}
}
