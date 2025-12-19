package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/schema/mixins"
)

// KioskSession holds the schema definition for the KioskSession entity.
// A KioskSession tracks when an admin has activated kiosk mode.
type KioskSession struct {
	ent.Schema
}

func (KioskSession) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.BaseMixin{},
	}
}

func (KioskSession) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("is_active"),
	}
}

// Fields of the KioskSession.
func (KioskSession) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("is_active").
			Default(true).
			Comment("Whether kiosk mode is currently active"),
		field.Time("unlocked_until").
			Optional().
			Nillable().
			Comment("When the temporary admin unlock expires (null = locked)"),
	}
}

// Edges of the KioskSession.
func (KioskSession) Edges() []ent.Edge {
	return []ent.Edge{
		// Required: which admin user owns this kiosk session
		edge.From("user", User.Type).
			Ref("kiosk_session").
			Unique().
			Required().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
