package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/schema/mixins"
)

// Borrower holds the schema definition for the Borrower entity.
// Borrowers are individuals who can check out equipment. They are NOT system users.
type Borrower struct {
	ent.Schema
}

func (Borrower) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.BaseMixin{},
		GroupMixin{ref: "borrowers"},
	}
}

func (Borrower) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
		index.Fields("email"),
		index.Fields("is_active"),
	}
}

// Fields of the Borrower.
func (Borrower) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(255).
			Comment("Full name of the borrower"),
		field.String("email").
			NotEmpty().
			MaxLen(255).
			Comment("Contact email address"),
		field.String("phone").
			Optional().
			MaxLen(50).
			Comment("Contact phone number"),
		field.String("organization").
			Optional().
			MaxLen(255).
			Comment("Company or startup name"),
		field.String("student_id").
			Optional().
			MaxLen(100).
			Comment("University student ID"),
		field.Text("notes").
			Optional().
			MaxLen(1000).
			Comment("Additional notes about the borrower"),
		field.Bool("is_active").
			Default(true).
			Comment("Whether borrower can currently check out equipment"),
		field.Bool("self_registered").
			Default(false).
			Comment("Whether borrower registered themselves via kiosk self-service"),
	}
}

// Edges of the Borrower.
func (Borrower) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("loans", Loan.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
