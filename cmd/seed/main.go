// cmd/seed/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const (
	defaultDBHost      = "host.docker.internal"
	defaultPostsDBPort = 5432
	defaultBotsDBPort  = 5433
	defaultProfDBPort  = 5434

	nProfiles = 25
	nBots     = 25
	nPosts    = 25
)

type dbCfg struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env: %s", k)
	}
	return v
}

func mustIntEnv(k string, def int) int {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid int env %s=%q: %v", k, v, err)
	}
	return n
}

func dsn(c dbCfg) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   c.Name,
	}
	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()
	return u.String()
}

func connect(ctx context.Context, name string, cfg dbCfg) *pgxpool.Pool {
	dsn := dsn(cfg)

	var lastErr error
	for i := 0; i < 30; i++ {
		pool, err := pgxpool.New(ctx, dsn)
		if err == nil {
			pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
			pingErr := pool.Ping(pingCtx)
			cancel()

			if pingErr == nil {
				return pool
			}

			lastErr = pingErr
			pool.Close()
		} else {
			lastErr = err
		}

		time.Sleep(1 * time.Second)
	}

	log.Fatalf("db connect failed (%s): %v", name, lastErr)
	return nil
}

type seededProfile struct {
	ID     uuid.UUID
	Name   string
	Email  string
	Prompt string
}

type seededBot struct {
	ID                 uuid.UUID
	Name               string
	SystemPrompt       string
	ModerationRequired bool
	ProfileIDs         []uuid.UUID
}

type seededPost struct {
	ID          uuid.UUID
	GroupID     uuid.UUID
	BotID       uuid.UUID
	BotName     string
	ProfileID   uuid.UUID
	ProfileName string
	UserPrompt  string
	Platform    string
	Type        *string
	Title       string
	Text        string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	userIDRaw := mustEnv("USER_ID")
	userID, err := uuid.Parse(userIDRaw)
	if err != nil {
		log.Fatalf("USER_ID is not uuid: %q: %v", userIDRaw, err)
	}

	postsCfg := dbCfg{
		Host:     getenv("SEED_DB_HOST", defaultDBHost),
		Port:     mustIntEnv("SEED_POSTS_DB_PORT", defaultPostsDBPort),
		User:     mustEnv("LOCAL_POSTS_DATABASE_USER"),
		Password: mustEnv("LOCAL_POSTS_DATABASE_PASSWORD"),
		Name:     mustEnv("LOCAL_POSTS_DATABASE_NAME"),
	}
	botsCfg := dbCfg{
		Host:     getenv("SEED_DB_HOST", defaultDBHost),
		Port:     mustIntEnv("SEED_BOTS_DB_PORT", defaultBotsDBPort),
		User:     mustEnv("LOCAL_BOTS_DATABASE_USER"),
		Password: mustEnv("LOCAL_BOTS_DATABASE_PASSWORD"),
		Name:     mustEnv("LOCAL_BOTS_DATABASE_NAME"),
	}
	profilesCfg := dbCfg{
		Host:     getenv("SEED_DB_HOST", defaultDBHost),
		Port:     mustIntEnv("SEED_PROFILES_DB_PORT", defaultProfDBPort),
		User:     mustEnv("LOCAL_PROFILES_DATABASE_EUSER"),
		Password: mustEnv("LOCAL_PROFILES_DATABASE_PASSWORD"),
		Name:     mustEnv("LOCAL_PROFILES_DATABASE_NAME"),
	}

	postsDB := connect(ctx, "posts", postsCfg)
	defer postsDB.Close()

	botsDB := connect(ctx, "bots", botsCfg)
	defer botsDB.Close()

	profilesDB := connect(ctx, "profiles", profilesCfg)
	defer profilesDB.Close()

	profiles := make([]seededProfile, 0, nProfiles)
	for i := 0; i < nProfiles; i++ {
		id := uuid.New()
		profiles = append(profiles, seededProfile{
			ID:     id,
			Name:   fmt.Sprintf("Seed Profile %02d", i+1),
			Email:  fmt.Sprintf("seed-profile-%02d@example.com", i+1),
			Prompt: fmt.Sprintf("System prompt for profile %02d", i+1),
		})
	}

	bots := make([]seededBot, 0, nBots)
	for i := 0; i < nBots; i++ {
		id := uuid.New()
		mr := i%3 == 0

		k := 1 + rng.Intn(4) // 1..4 profiles per bot
		used := map[uuid.UUID]struct{}{}
		pids := make([]uuid.UUID, 0, k)
		for len(pids) < k {
			p := profiles[rng.Intn(len(profiles))].ID
			if _, ok := used[p]; ok {
				continue
			}
			used[p] = struct{}{}
			pids = append(pids, p)
		}

		bots = append(bots, seededBot{
			ID:                 id,
			Name:               fmt.Sprintf("Seed Bot %02d", i+1),
			SystemPrompt:       fmt.Sprintf("System prompt for bot %02d", i+1),
			ModerationRequired: mr,
			ProfileIDs:         pids,
		})
	}

	postTypes := []string{"opinion", "knowledge", "history"}
	posts := make([]seededPost, 0, nPosts)
	for i := 0; i < nPosts; i++ {
		id := uuid.New()
		groupID := uuid.New()

		b := bots[rng.Intn(len(bots))]
		p := profiles[rng.Intn(len(profiles))]

		var pt *string
		t := postTypes[i%len(postTypes)]
		pt = &t

		posts = append(posts, seededPost{
			ID:          id,
			GroupID:     groupID,
			BotID:       b.ID,
			BotName:     b.Name,
			ProfileID:   p.ID,
			ProfileName: p.Name,
			UserPrompt:  fmt.Sprintf("Seed user prompt %02d", i+1),
			Platform:    "otveti",
			Type:        pt,
			Title:       fmt.Sprintf("Seed Post %02d", i+1),
			Text:        fmt.Sprintf("Seed post text %02d. Generated at %s.", i+1, time.Now().UTC().Format(time.RFC3339)),
		})
	}

	if err := seedProfiles(ctx, profilesDB, userID, profiles); err != nil {
		log.Fatalf("seed profiles failed: %v", err)
	}
	if err := seedBots(ctx, botsDB, userID, bots); err != nil {
		log.Fatalf("seed bots failed: %v", err)
	}
	if err := seedPosts(ctx, postsDB, userID, posts); err != nil {
		log.Fatalf("seed posts failed: %v", err)
	}

	log.Printf("seed ok: profiles=%d bots=%d posts=%d user_id=%s", len(profiles), len(bots), len(posts), userID.String())
}

func seedProfiles(ctx context.Context, db *pgxpool.Pool, userID uuid.UUID, profiles []seededProfile) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const q = `
INSERT INTO profiles (id, name, email, system_prompt, user_id)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO NOTHING;
`
	for _, p := range profiles {
		if _, err := tx.Exec(ctx, q, p.ID, p.Name, p.Email, p.Prompt, userID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func seedBots(ctx context.Context, db *pgxpool.Pool, userID uuid.UUID, bots []seededBot) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const q = `
INSERT INTO bots (id, bot_name, system_prompt, moderation_required, is_deleted, profile_ids, profiles_count, user_id)
VALUES ($1, $2, $3, $4, false, $5, $6, $7)
ON CONFLICT (id) DO NOTHING;
`

	for _, b := range bots {
		if _, err := tx.Exec(ctx, q,
			b.ID,
			b.Name,
			b.SystemPrompt,
			b.ModerationRequired,
			b.ProfileIDs,
			len(b.ProfileIDs),
			userID,
		); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func seedPosts(ctx context.Context, db *pgxpool.Pool, userID uuid.UUID, posts []seededPost) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const q = `
INSERT INTO posts (
	id,
	bot_id, bot_name,
	profile_id, profile_name,
	user_prompt,
	group_id,
	platform_type,
	post_type,
	post_title,
	post_text,
	user_id,
	is_seen
) VALUES (
	$1,
	$2, $3,
	$4, $5,
	$6,
	$7,
	$8::platform_type_enum,
	$9::post_type_enum,
	$10,
	$11,
	$12,
	false
)
ON CONFLICT (id) DO NOTHING;
`

	for _, p := range posts {
		postType := ""
		if p.Type != nil {
			postType = *p.Type
		}
		if _, err := tx.Exec(ctx, q,
			p.ID,
			p.BotID, p.BotName,
			p.ProfileID, p.ProfileName,
			p.UserPrompt,
			p.GroupID,
			p.Platform,
			postType,
			p.Title,
			p.Text,
			userID,
		); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
