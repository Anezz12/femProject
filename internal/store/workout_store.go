package store

import (
	"database/sql"
	"time"
)

type Workout struct {
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int       `json:"id"`
	WorkoutID       int       `json:"workout_id"`
	ExerciseName    string    `json:"exercise_name"`
	Sets            int       `json:"sets"`
	Reps            *int      `json:"reps,omitempty"`
	DurationSeconds *int      `json:"duration_seconds,omitempty"`
	Weight          *float64  `json:"weight,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	OrderIndex      int       `json:"order_index"`
	CreatedAt       time.Time `json:"created_at"`
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutByID(id int64) (*Workout, error)
}

func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert workout
	query := `
        INSERT INTO workouts (title, description, duration_minutes, calories_burned)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at
    `

	err = tx.QueryRow(
		query,
		workout.Title,
		workout.Description,
		workout.DurationMinutes,
		workout.CaloriesBurned,
	).Scan(&workout.ID, &workout.CreatedAt, &workout.UpdatedAt) // FIX: scan 3 field

	if err != nil {
		return nil, err
	}

	// Insert entries
	for i := range workout.Entries {
		entryQuery := `
            INSERT INTO workout_entries (
                workout_id, exercise_name, sets, reps, 
                duration_seconds, weight, notes, order_index
            )
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
            RETURNING id, created_at
        `

		err = tx.QueryRow(
			entryQuery,
			workout.ID,
			workout.Entries[i].ExerciseName,
			workout.Entries[i].Sets,
			workout.Entries[i].Reps,
			workout.Entries[i].DurationSeconds,
			workout.Entries[i].Weight,
			workout.Entries[i].Notes,
			workout.Entries[i].OrderIndex,
		).Scan(&workout.Entries[i].ID, &workout.Entries[i].CreatedAt)

		if err != nil {
			return nil, err
		}

		// Set WorkoutID
		workout.Entries[i].WorkoutID = workout.ID
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
	// Get workout
	query := `
        SELECT id, title, description, duration_minutes, 
               calories_burned, created_at, updated_at
        FROM workouts
        WHERE id = $1
    `

	var workout Workout
	err := pg.db.QueryRow(query, id).Scan(
		&workout.ID,
		&workout.Title,
		&workout.Description,
		&workout.DurationMinutes,
		&workout.CaloriesBurned,
		&workout.CreatedAt,
		&workout.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Get entries
	entriesQuery := `
        SELECT id, workout_id, exercise_name, sets, reps, 
               duration_seconds, weight, notes, order_index, created_at
        FROM workout_entries
        WHERE workout_id = $1
        ORDER BY order_index
    `

	rows, err := pg.db.Query(entriesQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []WorkoutEntry
	for rows.Next() {
		var entry WorkoutEntry
		err := rows.Scan(
			&entry.ID,
			&entry.WorkoutID,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.DurationSeconds,
			&entry.Weight,
			&entry.Notes,
			&entry.OrderIndex,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	workout.Entries = entries

	return &workout, nil
}
