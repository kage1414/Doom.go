package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// GameSettings holds the schema definition for the GameSettings entity.
type GameSettings struct {
	ent.Schema
}

// Fields of the GameSettings.
func (GameSettings) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Default("default").
			Comment("Settings ID - using 'default' for main settings"),
		field.Float("fire_rate").
			Default(0.08).
			Comment("Fire rate in seconds between shots"),
		field.Float("bullet_speed").
			Default(22.0).
			Comment("Bullet speed multiplier"),
		field.Int("level_count").
			Default(5).
			Comment("Number of levels to play"),
		field.Time("created_at").
			Optional().
			Comment("When these settings were created"),
		field.Time("updated_at").
			Optional().
			Comment("When these settings were last updated"),
	}
}

// Edges of the GameSettings.
func (GameSettings) Edges() []ent.Edge {
	return nil
}

// Indexes of the GameSettings.
func (GameSettings) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
	}
}
