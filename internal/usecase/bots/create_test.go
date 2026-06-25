package bots

import (
	"context"
	"testing"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

// --- test doubles for the usecase collaborators ---

type fakeRepo struct {
	createCalls int
	createErr   error
	lastCreated model.Bot
}

func (f *fakeRepo) Create(_ context.Context, b model.Bot) error {
	f.createCalls++
	f.lastCreated = b
	return f.createErr
}
func (f *fakeRepo) Get(context.Context, uuid.UUID) (model.Bot, error) {
	return model.Bot{}, nil
}
func (f *fakeRepo) List(context.Context) ([]model.Bot, error)   { return nil, nil }
func (f *fakeRepo) Update(context.Context, model.Bot) error     { return nil }
func (f *fakeRepo) Delete(context.Context, uuid.UUID) error     { return nil }
func (f *fakeRepo) AddProfileID(context.Context, uuid.UUID, uuid.UUID) error    { return nil }
func (f *fakeRepo) RemoveProfileID(context.Context, uuid.UUID, uuid.UUID) error { return nil }
func (f *fakeRepo) Search(context.Context, string) ([]model.Bot, error) {
	return nil, nil
}
func (f *fakeRepo) GetSummary(context.Context) (model.BotsSummary, error) {
	return model.BotsSummary{}, nil
}

type fakeValidator struct {
	err       error
	gotIDs    []uuid.UUID
	callCount int
}

func (f *fakeValidator) ValidateProfilesExist(_ context.Context, ids []uuid.UUID) error {
	f.callCount++
	f.gotIDs = ids
	return f.err
}

type fakeGateway struct {
	profiles  []model.Profile
	err       error
	callCount int
}

func (f *fakeGateway) GetProfilesByIDs(_ context.Context, _ []uuid.UUID) ([]model.Profile, error) {
	f.callCount++
	return f.profiles, f.err
}

func newService(r *fakeRepo, v *fakeValidator, g *fakeGateway) *Service {
	return NewService(r, v, g)
}

func TestCreate_TrimsNameAndRejectsEmpty(t *testing.T) {
	repo := &fakeRepo{}
	s := newService(repo, &fakeValidator{}, &fakeGateway{})

	_, err := s.Create(context.Background(), model.Bot{Name: "   "})
	if err == nil {
		t.Fatal("expected validation error for blank name, got nil")
	}
	if !errors.Is(err, constants.ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
	if repo.createCalls != 0 {
		t.Fatalf("repo.Create must not be called on validation failure, called %d times", repo.createCalls)
	}
}

func TestCreate_Success(t *testing.T) {
	repo := &fakeRepo{}
	validator := &fakeValidator{}
	s := newService(repo, validator, &fakeGateway{})

	ids := []uuid.UUID{uuid.New(), uuid.New()}
	got, err := s.Create(context.Background(), model.Bot{Name: "  My Bot  ", ProfileIDs: ids})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "My Bot" {
		t.Fatalf("expected trimmed name %q, got %q", "My Bot", got.Name)
	}
	if got.ProfilesCount != len(ids) {
		t.Fatalf("expected ProfilesCount %d, got %d", len(ids), got.ProfilesCount)
	}
	if got.ID == uuid.Nil {
		t.Fatal("expected a generated ID, got uuid.Nil")
	}
	if repo.createCalls != 1 {
		t.Fatalf("expected repo.Create to be called once, got %d", repo.createCalls)
	}
	if validator.callCount != 1 {
		t.Fatalf("expected profile validation to run once, got %d", validator.callCount)
	}
}

func TestCreate_PropagatesValidatorError(t *testing.T) {
	repo := &fakeRepo{}
	validator := &fakeValidator{err: constants.ErrNotFound}
	s := newService(repo, validator, &fakeGateway{})

	_, err := s.Create(context.Background(), model.Bot{Name: "Bot", ProfileIDs: []uuid.UUID{uuid.New()}})
	if !errors.Is(err, constants.ErrNotFound) {
		t.Fatalf("expected ErrNotFound from validator, got %v", err)
	}
	if repo.createCalls != 0 {
		t.Fatalf("repo.Create must not be called when validation fails, called %d", repo.createCalls)
	}
}

func TestCreateWithProfiles_SkipsGatewayWhenNoProfiles(t *testing.T) {
	repo := &fakeRepo{}
	gw := &fakeGateway{}
	s := newService(repo, &fakeValidator{}, gw)

	bot, profiles, err := s.CreateWithProfiles(context.Background(), model.Bot{Name: "Bot"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bot.Name != "Bot" {
		t.Fatalf("expected name %q, got %q", "Bot", bot.Name)
	}
	if len(profiles) != 0 {
		t.Fatalf("expected no profiles, got %d", len(profiles))
	}
	if gw.callCount != 0 {
		t.Fatalf("gateway must not be queried when bot has no profiles, called %d", gw.callCount)
	}
}

func TestCreateWithProfiles_ResolvesProfiles(t *testing.T) {
	repo := &fakeRepo{}
	gw := &fakeGateway{profiles: []model.Profile{{ID: uuid.New()}}}
	s := newService(repo, &fakeValidator{}, gw)

	_, profiles, err := s.CreateWithProfiles(context.Background(), model.Bot{
		Name:       "Bot",
		ProfileIDs: []uuid.UUID{uuid.New()},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 resolved profile, got %d", len(profiles))
	}
	if gw.callCount != 1 {
		t.Fatalf("expected gateway to be queried once, got %d", gw.callCount)
	}
	if repo.createCalls != 1 {
		t.Fatalf("expected exactly one repo.Create (no redundant re-read), got %d", repo.createCalls)
	}
}
