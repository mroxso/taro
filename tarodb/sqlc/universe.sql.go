// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: universe.sql

package sqlc

import (
	"context"
	"time"
)

const deleteUniverseServer = `-- name: DeleteUniverseServer :exec
DELETE FROM universe_servers
WHERE server_host = $1 OR id = $2
`

type DeleteUniverseServerParams struct {
	TargetServer string
	TargetID     int32
}

func (q *Queries) DeleteUniverseServer(ctx context.Context, arg DeleteUniverseServerParams) error {
	_, err := q.db.ExecContext(ctx, deleteUniverseServer, arg.TargetServer, arg.TargetID)
	return err
}

const fetchUniverseKeys = `-- name: FetchUniverseKeys :many
SELECT leaves.minting_point, leaves.script_key_bytes
FROM universe_leaves leaves
WHERE leaves.leaf_node_namespace = $1
`

type FetchUniverseKeysRow struct {
	MintingPoint   []byte
	ScriptKeyBytes []byte
}

func (q *Queries) FetchUniverseKeys(ctx context.Context, namespace string) ([]FetchUniverseKeysRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchUniverseKeys, namespace)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchUniverseKeysRow
	for rows.Next() {
		var i FetchUniverseKeysRow
		if err := rows.Scan(&i.MintingPoint, &i.ScriptKeyBytes); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchUniverseRoot = `-- name: FetchUniverseRoot :one
SELECT universe_roots.asset_id, group_key, mssmt_nodes.hash_key root_hash, 
       mssmt_nodes.sum root_sum, genesis_assets.asset_tag asset_name
FROM universe_roots
JOIN mssmt_roots 
    ON universe_roots.namespace_root = mssmt_roots.namespace
JOIN mssmt_nodes 
    ON mssmt_nodes.hash_key = mssmt_roots.root_hash AND
       mssmt_nodes.namespace = mssmt_roots.namespace
JOIN genesis_assets
     ON genesis_assets.asset_id = universe_roots.asset_id
WHERE mssmt_nodes.namespace = $1
`

type FetchUniverseRootRow struct {
	AssetID   []byte
	GroupKey  []byte
	RootHash  []byte
	RootSum   int64
	AssetName string
}

func (q *Queries) FetchUniverseRoot(ctx context.Context, namespace string) (FetchUniverseRootRow, error) {
	row := q.db.QueryRowContext(ctx, fetchUniverseRoot, namespace)
	var i FetchUniverseRootRow
	err := row.Scan(
		&i.AssetID,
		&i.GroupKey,
		&i.RootHash,
		&i.RootSum,
		&i.AssetName,
	)
	return i, err
}

const insertUniverseLeaf = `-- name: InsertUniverseLeaf :exec
INSERT INTO universe_leaves (
    asset_genesis_id, script_key_bytes, universe_root_id, leaf_node_key, 
    leaf_node_namespace, minting_point
) VALUES (
    $1, $2, $3, $4,
    $5, $6
)
`

type InsertUniverseLeafParams struct {
	AssetGenesisID    int32
	ScriptKeyBytes    []byte
	UniverseRootID    int32
	LeafNodeKey       []byte
	LeafNodeNamespace string
	MintingPoint      []byte
}

func (q *Queries) InsertUniverseLeaf(ctx context.Context, arg InsertUniverseLeafParams) error {
	_, err := q.db.ExecContext(ctx, insertUniverseLeaf,
		arg.AssetGenesisID,
		arg.ScriptKeyBytes,
		arg.UniverseRootID,
		arg.LeafNodeKey,
		arg.LeafNodeNamespace,
		arg.MintingPoint,
	)
	return err
}

const insertUniverseServer = `-- name: InsertUniverseServer :exec
INSERT INTO universe_servers(
    server_host, last_sync_time
) VALUES (
    $1, $2
)
`

type InsertUniverseServerParams struct {
	ServerHost   string
	LastSyncTime time.Time
}

func (q *Queries) InsertUniverseServer(ctx context.Context, arg InsertUniverseServerParams) error {
	_, err := q.db.ExecContext(ctx, insertUniverseServer, arg.ServerHost, arg.LastSyncTime)
	return err
}

const listUniverseServers = `-- name: ListUniverseServers :many
SELECT id, server_host, last_sync_time FROM universe_servers
`

func (q *Queries) ListUniverseServers(ctx context.Context) ([]UniverseServer, error) {
	rows, err := q.db.QueryContext(ctx, listUniverseServers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UniverseServer
	for rows.Next() {
		var i UniverseServer
		if err := rows.Scan(&i.ID, &i.ServerHost, &i.LastSyncTime); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const logServerSync = `-- name: LogServerSync :exec
UPDATE universe_servers
SET last_sync_time = $1
WHERE server_host = $2
`

type LogServerSyncParams struct {
	NewSyncTime  time.Time
	TargetServer string
}

func (q *Queries) LogServerSync(ctx context.Context, arg LogServerSyncParams) error {
	_, err := q.db.ExecContext(ctx, logServerSync, arg.NewSyncTime, arg.TargetServer)
	return err
}

const queryUniverseLeaves = `-- name: QueryUniverseLeaves :many
SELECT leaves.script_key_bytes, gen.gen_asset_id, nodes.value genesis_proof, 
       nodes.sum sum_amt
FROM universe_leaves leaves
JOIN mssmt_nodes nodes
    ON leaves.leaf_node_key = nodes.key AND
        leaves.leaf_node_namespace = nodes.namespace
JOIN genesis_info_view gen
    ON leaves.asset_genesis_id = gen.gen_asset_id
WHERE leaves.leaf_node_namespace = $1 
        AND 
    (leaves.minting_point = $2 OR 
        $2 IS NULL) 
        AND
    (leaves.script_key_bytes = $3 OR 
        $3 IS NULL)
`

type QueryUniverseLeavesParams struct {
	Namespace         string
	MintingPointBytes []byte
	ScriptKeyBytes    []byte
}

type QueryUniverseLeavesRow struct {
	ScriptKeyBytes []byte
	GenAssetID     int32
	GenesisProof   []byte
	SumAmt         int64
}

func (q *Queries) QueryUniverseLeaves(ctx context.Context, arg QueryUniverseLeavesParams) ([]QueryUniverseLeavesRow, error) {
	rows, err := q.db.QueryContext(ctx, queryUniverseLeaves, arg.Namespace, arg.MintingPointBytes, arg.ScriptKeyBytes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []QueryUniverseLeavesRow
	for rows.Next() {
		var i QueryUniverseLeavesRow
		if err := rows.Scan(
			&i.ScriptKeyBytes,
			&i.GenAssetID,
			&i.GenesisProof,
			&i.SumAmt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const universeLeaves = `-- name: UniverseLeaves :many
SELECT id, asset_genesis_id, minting_point, script_key_bytes, universe_root_id, leaf_node_key, leaf_node_namespace FROM universe_leaves
`

func (q *Queries) UniverseLeaves(ctx context.Context) ([]UniverseLeafe, error) {
	rows, err := q.db.QueryContext(ctx, universeLeaves)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UniverseLeafe
	for rows.Next() {
		var i UniverseLeafe
		if err := rows.Scan(
			&i.ID,
			&i.AssetGenesisID,
			&i.MintingPoint,
			&i.ScriptKeyBytes,
			&i.UniverseRootID,
			&i.LeafNodeKey,
			&i.LeafNodeNamespace,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const universeRoots = `-- name: UniverseRoots :many
SELECT universe_roots.asset_id, group_key, mssmt_roots.root_hash root_hash,
       mssmt_nodes.sum root_sum, genesis_assets.asset_tag asset_name
FROM universe_roots
JOIN mssmt_roots
    ON universe_roots.namespace_root = mssmt_roots.namespace
JOIN mssmt_nodes
    ON mssmt_nodes.hash_key = mssmt_roots.root_hash AND
       mssmt_nodes.namespace = mssmt_roots.namespace
JOIN genesis_assets
    ON genesis_assets.asset_id = universe_roots.asset_id
`

type UniverseRootsRow struct {
	AssetID   []byte
	GroupKey  []byte
	RootHash  []byte
	RootSum   int64
	AssetName string
}

func (q *Queries) UniverseRoots(ctx context.Context) ([]UniverseRootsRow, error) {
	rows, err := q.db.QueryContext(ctx, universeRoots)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UniverseRootsRow
	for rows.Next() {
		var i UniverseRootsRow
		if err := rows.Scan(
			&i.AssetID,
			&i.GroupKey,
			&i.RootHash,
			&i.RootSum,
			&i.AssetName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upsertUniverseRoot = `-- name: UpsertUniverseRoot :one
INSERT INTO universe_roots (
    namespace_root, asset_id, group_key
) VALUES (
    $1, $2, $3
) ON CONFLICT (namespace_root)
    DO UPDATE SET namespace_root = $1
RETURNING id
`

type UpsertUniverseRootParams struct {
	NamespaceRoot string
	AssetID       []byte
	GroupKey      []byte
}

func (q *Queries) UpsertUniverseRoot(ctx context.Context, arg UpsertUniverseRootParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, upsertUniverseRoot, arg.NamespaceRoot, arg.AssetID, arg.GroupKey)
	var id int32
	err := row.Scan(&id)
	return id, err
}
