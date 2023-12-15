package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

type Model struct {
	DB DBModel
}

// NewModels returns a model type with database connection pool
func NewModel(db *sql.DB) Model {
	return Model{
		DB: DBModel{DB: db},
	}
}

// Collection is the type for collection
type Collection struct {
	ID          string   `json:"id,omitempty"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Objects     []Object `json:"objects"`
	Alias       string   `json:"alias"`
	Can_read    bool     `json:"can_read"`
	Can_write   bool     `json:"can_write"`
	Media_types string   `json:"media_types"`
}

// Collections is the type for collections
type Collections struct {
	Collections []Collection
}

// type ManifestRecord is the type for Manifest-Record
type ManifestRecord struct {
	ID         string    `json:"id"`
	Date_added time.Time `json:"date_added"`
	Version    string    `json:"version"`
	Media_type string    `json:"media_type"`
}

// Manifest is the type for manifests
type Manifest struct {
	More            bool             `json:"more"`
	ManifestRecords []ManifestRecord `json:"objects"`
}

// Object is the type for objects
type Object struct {
	ID           int       `json:"id"`
	CollectionID int       `json:"collection_id"`
	Type         string    `json:"type"`
	Content      string    `json:"content"`
	Timestamp    time.Time `json:"timestamp"`
}

// User is the type for users
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Status is the type for statuses
type Status struct {
	ID                string          `json:"id"`
	Status            string          `json:"status"`
	Request_timestamp time.Time       `json:"request_timestamp"`
	Total_count       int             `json:"total_count"`
	Success_count     int             `json:"success_count"`
	Successes         []StatusDetails `json:"successes"`
	Failure_count     int             `json:"failure_count"`
	Failures          []StatusDetails `json:"failures"`
	Pending_count     int             `json:"pending_count"`
	Pendings          []StatusDetails `json:"pendings"`
}

// StatusDetails is the type for status details
type StatusDetails struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Message string `json:"message"`
}

// Envelope is the type for envelopes
type Envelope struct {
	More    bool     `json:"more"`
	Next    string   `json:"next"`
	Objects []Object `json:"objects"` //list of type <STIX Object>
}

// Discovery is the type for discovery
type Discovery struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Contact     string   `json:"contact"`
	Default     string   `json:"default"`
	API_ROOTS   []string `json:"api_roots"`
}

// API_ROOT is the type for API Roots
type API_ROOTS struct {
	Title              string `json:"title"`
	Description        string `json:"description"`
	Versions           string `json:"versions"`
	Max_content_length int    `json:"max_content_length"`
}

// GetApiRoot returns one API_ROOT by name (api-root)
func (m *DBModel) GetApiRoot(apiroot string) (API_ROOTS, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// "ar" is short for "api-root"
	var ar API_ROOTS

	row := m.DB.QueryRowContext(ctx,
		`select 
			title, description, versions, max_content_length
		from 
			apiroot
		where name = ?`, apiroot)

	err := row.Scan(
		&ar.Title,
		&ar.Description,
		&ar.Versions,
		&ar.Max_content_length,
	)
	if err != nil {
		return ar, err
	}
	return ar, nil

}

// GetAllApiRoots returns a slice of all API_ROOTS
func (m *DBModel) GetAllApiRoots() (Discovery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var apiroots []*API_ROOTS
	query_api_roots := `
		select
			name
		from
			apiroot
	`
	row := m.DB.QueryRowContext(ctx, query_api_roots)
	fmt.Println(row)

	var discovery Discovery
	query := `
		select 
			* 
		from 
			discovery
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return Discovery{}, err
	}
	defer rows.Close()

	for rows.Next() {
		// ar short for api root
		var ar API_ROOTS
		err = rows.Scan(
			&discovery.Title,
			&discovery.Description,
			&discovery.Contact,
			&discovery.Default,
			&discovery.API_ROOTS,
		)
		if err != nil {
			return Discovery{}, err
		}
		apiroots = append(apiroots, &ar)
	}
	return discovery, nil
}

// GetAllCollections returns all the collections from the api-root
func (m *DBModel) GetAllCollections() ([]*Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var collections []*Collection

	query := `
		select
			title, description, alias, can_read, can_write, media_types
		from
			collection
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return collections, nil
	}
	defer rows.Close()

	for rows.Next() {
		var c Collection
		err = rows.Scan(
			&c.ID,
			&c.Title,
			&c.Description,
			&c.Alias,
			&c.Can_read,
			&c.Can_write,
			&c.Media_types,
		)
		if err != nil {
			return nil, err
		}
		collections = append(collections, &c)
	}
	return collections, nil

}

// GetCollectionID returns the collection based on passed ID
func (m *DBModel) GetCollectionID(id int) (Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var collection Collection
	query := `
		select
			title, description, alias, can_read, can_write, media_types
		from
			collection
		where id = ?
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&collection.ID,
		&collection.Title,
		&collection.Description,
		&collection.Can_read,
		&collection.Can_write,
		&collection.Media_types,
	)
	if err != nil {
		return collection, err
	}
	return collection, nil
}

// GetManifestRecord returns the manifest-record (metadata) about the collection by ID
func (m *DBModel) GetManifestRecord(id int) (ManifestRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var mr ManifestRecord

	query := `select id, date_added, version, media_type 
			from manifest_records
			where id = ?
			`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&mr.ID,
		&mr.Date_added,
		&mr.Version,
		&mr.Media_type,
	)

	if err != nil {
		return mr, err
	}
	return mr, nil
}

// SaveObject saves an object into db
func (m *DBModel) SaveObject(object Object) error {
	// Save the object to MySQL database
	// ...

	// Notify subscribers about the new project
	NewEventBus().Publish(Event{
		CollectionID: object.CollectionID,
		Object:       object,
	})
	return nil
}
