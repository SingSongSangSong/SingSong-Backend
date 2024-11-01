// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// PostLike is an object representing the database table.
type PostLike struct {
	PostLikeID int64     `boil:"post_like_id" json:"post_like_id" toml:"post_like_id" yaml:"post_like_id"`
	MemberID   int64     `boil:"member_id" json:"member_id" toml:"member_id" yaml:"member_id"`
	PostID     int64     `boil:"post_id" json:"post_id" toml:"post_id" yaml:"post_id"`
	CreatedAt  null.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt  null.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
	DeletedAt  null.Time `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`

	R *postLikeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L postLikeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PostLikeColumns = struct {
	PostLikeID string
	MemberID   string
	PostID     string
	CreatedAt  string
	UpdatedAt  string
	DeletedAt  string
}{
	PostLikeID: "post_like_id",
	MemberID:   "member_id",
	PostID:     "post_id",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

var PostLikeTableColumns = struct {
	PostLikeID string
	MemberID   string
	PostID     string
	CreatedAt  string
	UpdatedAt  string
	DeletedAt  string
}{
	PostLikeID: "post_like.post_like_id",
	MemberID:   "post_like.member_id",
	PostID:     "post_like.post_id",
	CreatedAt:  "post_like.created_at",
	UpdatedAt:  "post_like.updated_at",
	DeletedAt:  "post_like.deleted_at",
}

// Generated where

var PostLikeWhere = struct {
	PostLikeID whereHelperint64
	MemberID   whereHelperint64
	PostID     whereHelperint64
	CreatedAt  whereHelpernull_Time
	UpdatedAt  whereHelpernull_Time
	DeletedAt  whereHelpernull_Time
}{
	PostLikeID: whereHelperint64{field: "`post_like`.`post_like_id`"},
	MemberID:   whereHelperint64{field: "`post_like`.`member_id`"},
	PostID:     whereHelperint64{field: "`post_like`.`post_id`"},
	CreatedAt:  whereHelpernull_Time{field: "`post_like`.`created_at`"},
	UpdatedAt:  whereHelpernull_Time{field: "`post_like`.`updated_at`"},
	DeletedAt:  whereHelpernull_Time{field: "`post_like`.`deleted_at`"},
}

// PostLikeRels is where relationship names are stored.
var PostLikeRels = struct {
}{}

// postLikeR is where relationships are stored.
type postLikeR struct {
}

// NewStruct creates a new relationship struct
func (*postLikeR) NewStruct() *postLikeR {
	return &postLikeR{}
}

// postLikeL is where Load methods for each relationship are stored.
type postLikeL struct{}

var (
	postLikeAllColumns            = []string{"post_like_id", "member_id", "post_id", "created_at", "updated_at", "deleted_at"}
	postLikeColumnsWithoutDefault = []string{"member_id", "post_id", "deleted_at"}
	postLikeColumnsWithDefault    = []string{"post_like_id", "created_at", "updated_at"}
	postLikePrimaryKeyColumns     = []string{"post_like_id"}
	postLikeGeneratedColumns      = []string{}
)

type (
	// PostLikeSlice is an alias for a slice of pointers to PostLike.
	// This should almost always be used instead of []PostLike.
	PostLikeSlice []*PostLike
	// PostLikeHook is the signature for custom PostLike hook methods
	PostLikeHook func(context.Context, boil.ContextExecutor, *PostLike) error

	postLikeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	postLikeType                 = reflect.TypeOf(&PostLike{})
	postLikeMapping              = queries.MakeStructMapping(postLikeType)
	postLikePrimaryKeyMapping, _ = queries.BindMapping(postLikeType, postLikeMapping, postLikePrimaryKeyColumns)
	postLikeInsertCacheMut       sync.RWMutex
	postLikeInsertCache          = make(map[string]insertCache)
	postLikeUpdateCacheMut       sync.RWMutex
	postLikeUpdateCache          = make(map[string]updateCache)
	postLikeUpsertCacheMut       sync.RWMutex
	postLikeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var postLikeAfterSelectHooks []PostLikeHook

var postLikeBeforeInsertHooks []PostLikeHook
var postLikeAfterInsertHooks []PostLikeHook

var postLikeBeforeUpdateHooks []PostLikeHook
var postLikeAfterUpdateHooks []PostLikeHook

var postLikeBeforeDeleteHooks []PostLikeHook
var postLikeAfterDeleteHooks []PostLikeHook

var postLikeBeforeUpsertHooks []PostLikeHook
var postLikeAfterUpsertHooks []PostLikeHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *PostLike) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range postLikeAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *PostLike) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range postLikeBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *PostLike) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range postLikeAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *PostLike) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range postLikeBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *PostLike) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range postLikeAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *PostLike) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range postLikeBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *PostLike) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range postLikeAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *PostLike) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range postLikeBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *PostLike) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range postLikeAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddPostLikeHook registers your hook function for all future operations.
func AddPostLikeHook(hookPoint boil.HookPoint, postLikeHook PostLikeHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		postLikeAfterSelectHooks = append(postLikeAfterSelectHooks, postLikeHook)
	case boil.BeforeInsertHook:
		postLikeBeforeInsertHooks = append(postLikeBeforeInsertHooks, postLikeHook)
	case boil.AfterInsertHook:
		postLikeAfterInsertHooks = append(postLikeAfterInsertHooks, postLikeHook)
	case boil.BeforeUpdateHook:
		postLikeBeforeUpdateHooks = append(postLikeBeforeUpdateHooks, postLikeHook)
	case boil.AfterUpdateHook:
		postLikeAfterUpdateHooks = append(postLikeAfterUpdateHooks, postLikeHook)
	case boil.BeforeDeleteHook:
		postLikeBeforeDeleteHooks = append(postLikeBeforeDeleteHooks, postLikeHook)
	case boil.AfterDeleteHook:
		postLikeAfterDeleteHooks = append(postLikeAfterDeleteHooks, postLikeHook)
	case boil.BeforeUpsertHook:
		postLikeBeforeUpsertHooks = append(postLikeBeforeUpsertHooks, postLikeHook)
	case boil.AfterUpsertHook:
		postLikeAfterUpsertHooks = append(postLikeAfterUpsertHooks, postLikeHook)
	}
}

// One returns a single postLike record from the query.
func (q postLikeQuery) One(ctx context.Context, exec boil.ContextExecutor) (*PostLike, error) {
	o := &PostLike{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysql: failed to execute a one query for post_like")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all PostLike records from the query.
func (q postLikeQuery) All(ctx context.Context, exec boil.ContextExecutor) (PostLikeSlice, error) {
	var o []*PostLike

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "mysql: failed to assign all query results to PostLike slice")
	}

	if len(postLikeAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all PostLike records in the query.
func (q postLikeQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to count post_like rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q postLikeQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "mysql: failed to check if post_like exists")
	}

	return count > 0, nil
}

// PostLikes retrieves all the records using an executor.
func PostLikes(mods ...qm.QueryMod) postLikeQuery {
	mods = append(mods, qm.From("`post_like`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`post_like`.*"})
	}

	return postLikeQuery{q}
}

// FindPostLike retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPostLike(ctx context.Context, exec boil.ContextExecutor, postLikeID int64, selectCols ...string) (*PostLike, error) {
	postLikeObj := &PostLike{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `post_like` where `post_like_id`=?", sel,
	)

	q := queries.Raw(query, postLikeID)

	err := q.Bind(ctx, exec, postLikeObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysql: unable to select from post_like")
	}

	if err = postLikeObj.doAfterSelectHooks(ctx, exec); err != nil {
		return postLikeObj, err
	}

	return postLikeObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *PostLike) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("mysql: no post_like provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(postLikeColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	postLikeInsertCacheMut.RLock()
	cache, cached := postLikeInsertCache[key]
	postLikeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			postLikeAllColumns,
			postLikeColumnsWithDefault,
			postLikeColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(postLikeType, postLikeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(postLikeType, postLikeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `post_like` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `post_like` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `post_like` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, postLikePrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "mysql: unable to insert into post_like")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.PostLikeID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == postLikeMapping["post_like_id"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.PostLikeID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to populate default values for post_like")
	}

CacheNoHooks:
	if !cached {
		postLikeInsertCacheMut.Lock()
		postLikeInsertCache[key] = cache
		postLikeInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the PostLike.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *PostLike) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	postLikeUpdateCacheMut.RLock()
	cache, cached := postLikeUpdateCache[key]
	postLikeUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			postLikeAllColumns,
			postLikePrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("mysql: unable to update post_like, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `post_like` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, postLikePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(postLikeType, postLikeMapping, append(wl, postLikePrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update post_like row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by update for post_like")
	}

	if !cached {
		postLikeUpdateCacheMut.Lock()
		postLikeUpdateCache[key] = cache
		postLikeUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q postLikeQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update all for post_like")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to retrieve rows affected for post_like")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PostLikeSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("mysql: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), postLikePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `post_like` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, postLikePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update all in postLike slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to retrieve rows affected all in update all postLike")
	}
	return rowsAff, nil
}

var mySQLPostLikeUniqueColumns = []string{
	"post_like_id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *PostLike) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("mysql: no post_like provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(postLikeColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLPostLikeUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	postLikeUpsertCacheMut.RLock()
	cache, cached := postLikeUpsertCache[key]
	postLikeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			postLikeAllColumns,
			postLikeColumnsWithDefault,
			postLikeColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			postLikeAllColumns,
			postLikePrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("mysql: unable to upsert post_like, could not build update column list")
		}

		ret := strmangle.SetComplement(postLikeAllColumns, strmangle.SetIntersect(insert, update))

		cache.query = buildUpsertQueryMySQL(dialect, "`post_like`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `post_like` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(postLikeType, postLikeMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(postLikeType, postLikeMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "mysql: unable to upsert for post_like")
	}

	var lastID int64
	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.PostLikeID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == postLikeMapping["post_like_id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(postLikeType, postLikeMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to retrieve unique values for post_like")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to populate default values for post_like")
	}

CacheNoHooks:
	if !cached {
		postLikeUpsertCacheMut.Lock()
		postLikeUpsertCache[key] = cache
		postLikeUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single PostLike record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *PostLike) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("mysql: no PostLike provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), postLikePrimaryKeyMapping)
	sql := "DELETE FROM `post_like` WHERE `post_like_id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete from post_like")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by delete for post_like")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q postLikeQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("mysql: no postLikeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete all from post_like")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by deleteall for post_like")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PostLikeSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(postLikeBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), postLikePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `post_like` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, postLikePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete all from postLike slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by deleteall for post_like")
	}

	if len(postLikeAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *PostLike) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindPostLike(ctx, exec, o.PostLikeID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PostLikeSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PostLikeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), postLikePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `post_like`.* FROM `post_like` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, postLikePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to reload all in PostLikeSlice")
	}

	*o = slice

	return nil
}

// PostLikeExists checks if the PostLike row exists.
func PostLikeExists(ctx context.Context, exec boil.ContextExecutor, postLikeID int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `post_like` where `post_like_id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, postLikeID)
	}
	row := exec.QueryRowContext(ctx, sql, postLikeID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "mysql: unable to check if post_like exists")
	}

	return exists, nil
}

// Exists checks if the PostLike row exists.
func (o *PostLike) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return PostLikeExists(ctx, exec, o.PostLikeID)
}
