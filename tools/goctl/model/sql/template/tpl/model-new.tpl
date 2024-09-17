func new{{.upperStartCamelObject}}Model(conn sqlx.SqlConn{{if .withCache}}, c cache.NodeConf, opts ...cache.Option{{end}}) *default{{.upperStartCamelObject}}Model {
	return &default{{.upperStartCamelObject}}Model{
		{{if .withCache}}CachedConn: sqlc.NewNodeConn(conn, c, opts...){{else}}conn:conn{{end}},
		table:      {{.table}},
	}
}
