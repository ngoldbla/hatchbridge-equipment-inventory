package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/schema/mixins"
)

// Loan holds the schema definition for the Loan entity.
// A Loan tracks a single checkout of equipment to a borrower.
type Loan struct {
	ent.Schema
}

func (Loan) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.BaseMixin{},
		GroupMixin{ref: "loans"},
	}
}

func (Loan) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("checked_out_at"),
		index.Fields("due_at"),
		index.Fields("returned_at"),
	}
}

// Fields of the Loan.
func (Loan) Fields() []ent.Field {
	return []ent.Field{
		field.Time("checked_out_at").
			Default(time.Now).
			Comment("When the item was checked out"),
		field.Time("due_at").
			Comment("When the item is expected to be returned"),
		field.Time("returned_at").
			Optional().
			Nillable().
			Comment("When the item was actually returned (null = still on loan)"),
		field.Text("notes").
			Optional().
			MaxLen(1000).
			Comment("Notes made at checkout"),
		field.Text("return_notes").
			Optional().
			MaxLen(1000).
			Comment("Notes made at return (e.g., condition)"),
		field.Int("quantity").
			Default(1).
			Positive().
			Comment("Number of items borrowed (for items with quantity > 1)"),
		field.Bool("kiosk_action").
			Default(false).
			Comment("Whether this loan was created/returned via kiosk self-service"),
	}
}

// Edges of the Loan.
func (Loan) Edges() []ent.Edge {
	return []ent.Edge{
		// Required: which item is on loan
		edge.From("item", Item.Type).
			Ref("loans").
			Unique().
			Required(),
		// Required: who has the item
		edge.From("borrower", Borrower.Type).
			Ref("loans").
			Unique().
			Required(),
		// Optional: which admin processed the checkout
		edge.From("checked_out_by", User.Type).
			Ref("checkouts").
			Unique(),
		// Optional: which admin processed the return
		edge.From("returned_by", User.Type).
			Ref("returns").
			Unique(),
	}
}
