package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Chess struct {
	ent.Schema
}

func (Chess) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("主键"),
		field.String("before").Comment("初始棋面"),
		field.String("after").Comment("结束棋面"),
		field.Int64("count").Comment("频次"),
	}

}

func (Chess) Edges() []ent.Edge {
	return nil
}

func (Chess) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "chess",
		},
	}
}
