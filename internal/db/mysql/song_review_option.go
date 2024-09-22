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

// SongReviewOption is an object representing the database table.
type SongReviewOption struct {
	SongReviewOptionID int64       `boil:"song_review_option_id" json:"song_review_option_id" toml:"song_review_option_id" yaml:"song_review_option_id"`
	Title              null.String `boil:"title" json:"title,omitempty" toml:"title" yaml:"title,omitempty"`
	CreatedAt          null.Time   `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt          null.Time   `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
	DeletedAt          null.Time   `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`
	Enum               null.String `boil:"enum" json:"enum,omitempty" toml:"enum" yaml:"enum,omitempty"`

	R *songReviewOptionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L songReviewOptionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SongReviewOptionColumns = struct {
	SongReviewOptionID string
	Title              string
	CreatedAt          string
	UpdatedAt          string
	DeletedAt          string
	Enum               string
}{
	SongReviewOptionID: "song_review_option_id",
	Title:              "title",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
	DeletedAt:          "deleted_at",
	Enum:               "enum",
}

var SongReviewOptionTableColumns = struct {
	SongReviewOptionID string
	Title              string
	CreatedAt          string
	UpdatedAt          string
	DeletedAt          string
	Enum               string
}{
	SongReviewOptionID: "song_review_option.song_review_option_id",
	Title:              "song_review_option.title",
	CreatedAt:          "song_review_option.created_at",
	UpdatedAt:          "song_review_option.updated_at",
	DeletedAt:          "song_review_option.deleted_at",
	Enum:               "song_review_option.enum",
}

// Generated where

var SongReviewOptionWhere = struct {
	SongReviewOptionID whereHelperint64
	Title              whereHelpernull_String
	CreatedAt          whereHelpernull_Time
	UpdatedAt          whereHelpernull_Time
	DeletedAt          whereHelpernull_Time
	Enum               whereHelpernull_String
}{
	SongReviewOptionID: whereHelperint64{field: "`song_review_option`.`song_review_option_id`"},
	Title:              whereHelpernull_String{field: "`song_review_option`.`title`"},
	CreatedAt:          whereHelpernull_Time{field: "`song_review_option`.`created_at`"},
	UpdatedAt:          whereHelpernull_Time{field: "`song_review_option`.`updated_at`"},
	DeletedAt:          whereHelpernull_Time{field: "`song_review_option`.`deleted_at`"},
	Enum:               whereHelpernull_String{field: "`song_review_option`.`enum`"},
}

// SongReviewOptionRels is where relationship names are stored.
var SongReviewOptionRels = struct {
}{}

// songReviewOptionR is where relationships are stored.
type songReviewOptionR struct {
}

// NewStruct creates a new relationship struct
func (*songReviewOptionR) NewStruct() *songReviewOptionR {
	return &songReviewOptionR{}
}

// songReviewOptionL is where Load methods for each relationship are stored.
type songReviewOptionL struct{}

var (
	songReviewOptionAllColumns            = []string{"song_review_option_id", "title", "created_at", "updated_at", "deleted_at", "enum"}
	songReviewOptionColumnsWithoutDefault = []string{"title", "deleted_at", "enum"}
	songReviewOptionColumnsWithDefault    = []string{"song_review_option_id", "created_at", "updated_at"}
	songReviewOptionPrimaryKeyColumns     = []string{"song_review_option_id"}
	songReviewOptionGeneratedColumns      = []string{}
)

type (
	// SongReviewOptionSlice is an alias for a slice of pointers to SongReviewOption.
	// This should almost always be used instead of []SongReviewOption.
	SongReviewOptionSlice []*SongReviewOption
	// SongReviewOptionHook is the signature for custom SongReviewOption hook methods
	SongReviewOptionHook func(context.Context, boil.ContextExecutor, *SongReviewOption) error

	songReviewOptionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	songReviewOptionType                 = reflect.TypeOf(&SongReviewOption{})
	songReviewOptionMapping              = queries.MakeStructMapping(songReviewOptionType)
	songReviewOptionPrimaryKeyMapping, _ = queries.BindMapping(songReviewOptionType, songReviewOptionMapping, songReviewOptionPrimaryKeyColumns)
	songReviewOptionInsertCacheMut       sync.RWMutex
	songReviewOptionInsertCache          = make(map[string]insertCache)
	songReviewOptionUpdateCacheMut       sync.RWMutex
	songReviewOptionUpdateCache          = make(map[string]updateCache)
	songReviewOptionUpsertCacheMut       sync.RWMutex
	songReviewOptionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var songReviewOptionAfterSelectHooks []SongReviewOptionHook

var songReviewOptionBeforeInsertHooks []SongReviewOptionHook
var songReviewOptionAfterInsertHooks []SongReviewOptionHook

var songReviewOptionBeforeUpdateHooks []SongReviewOptionHook
var songReviewOptionAfterUpdateHooks []SongReviewOptionHook

var songReviewOptionBeforeDeleteHooks []SongReviewOptionHook
var songReviewOptionAfterDeleteHooks []SongReviewOptionHook

var songReviewOptionBeforeUpsertHooks []SongReviewOptionHook
var songReviewOptionAfterUpsertHooks []SongReviewOptionHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *SongReviewOption) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewOptionAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *SongReviewOption) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewOptionBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *SongReviewOption) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewOptionAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *SongReviewOption) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewOptionBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *SongReviewOption) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewOptionAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *SongReviewOption) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewOptionBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *SongReviewOption) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewOptionAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *SongReviewOption) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewOptionBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *SongReviewOption) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewOptionAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSongReviewOptionHook registers your hook function for all future operations.
func AddSongReviewOptionHook(hookPoint boil.HookPoint, songReviewOptionHook SongReviewOptionHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		songReviewOptionAfterSelectHooks = append(songReviewOptionAfterSelectHooks, songReviewOptionHook)
	case boil.BeforeInsertHook:
		songReviewOptionBeforeInsertHooks = append(songReviewOptionBeforeInsertHooks, songReviewOptionHook)
	case boil.AfterInsertHook:
		songReviewOptionAfterInsertHooks = append(songReviewOptionAfterInsertHooks, songReviewOptionHook)
	case boil.BeforeUpdateHook:
		songReviewOptionBeforeUpdateHooks = append(songReviewOptionBeforeUpdateHooks, songReviewOptionHook)
	case boil.AfterUpdateHook:
		songReviewOptionAfterUpdateHooks = append(songReviewOptionAfterUpdateHooks, songReviewOptionHook)
	case boil.BeforeDeleteHook:
		songReviewOptionBeforeDeleteHooks = append(songReviewOptionBeforeDeleteHooks, songReviewOptionHook)
	case boil.AfterDeleteHook:
		songReviewOptionAfterDeleteHooks = append(songReviewOptionAfterDeleteHooks, songReviewOptionHook)
	case boil.BeforeUpsertHook:
		songReviewOptionBeforeUpsertHooks = append(songReviewOptionBeforeUpsertHooks, songReviewOptionHook)
	case boil.AfterUpsertHook:
		songReviewOptionAfterUpsertHooks = append(songReviewOptionAfterUpsertHooks, songReviewOptionHook)
	}
}

// One returns a single songReviewOption record from the query.
func (q songReviewOptionQuery) One(ctx context.Context, exec boil.ContextExecutor) (*SongReviewOption, error) {
	o := &SongReviewOption{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysql: failed to execute a one query for song_review_option")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all SongReviewOption records from the query.
func (q songReviewOptionQuery) All(ctx context.Context, exec boil.ContextExecutor) (SongReviewOptionSlice, error) {
	var o []*SongReviewOption

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "mysql: failed to assign all query results to SongReviewOption slice")
	}

	if len(songReviewOptionAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all SongReviewOption records in the query.
func (q songReviewOptionQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to count song_review_option rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q songReviewOptionQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "mysql: failed to check if song_review_option exists")
	}

	return count > 0, nil
}

// SongReviewOptions retrieves all the records using an executor.
func SongReviewOptions(mods ...qm.QueryMod) songReviewOptionQuery {
	mods = append(mods, qm.From("`song_review_option`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`song_review_option`.*"})
	}

	return songReviewOptionQuery{q}
}

// FindSongReviewOption retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSongReviewOption(ctx context.Context, exec boil.ContextExecutor, songReviewOptionID int64, selectCols ...string) (*SongReviewOption, error) {
	songReviewOptionObj := &SongReviewOption{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `song_review_option` where `song_review_option_id`=?", sel,
	)

	q := queries.Raw(query, songReviewOptionID)

	err := q.Bind(ctx, exec, songReviewOptionObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysql: unable to select from song_review_option")
	}

	if err = songReviewOptionObj.doAfterSelectHooks(ctx, exec); err != nil {
		return songReviewOptionObj, err
	}

	return songReviewOptionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *SongReviewOption) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("mysql: no song_review_option provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(songReviewOptionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	songReviewOptionInsertCacheMut.RLock()
	cache, cached := songReviewOptionInsertCache[key]
	songReviewOptionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			songReviewOptionAllColumns,
			songReviewOptionColumnsWithDefault,
			songReviewOptionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(songReviewOptionType, songReviewOptionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(songReviewOptionType, songReviewOptionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `song_review_option` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `song_review_option` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `song_review_option` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, songReviewOptionPrimaryKeyColumns))
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
		return errors.Wrap(err, "mysql: unable to insert into song_review_option")
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

	o.SongReviewOptionID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == songReviewOptionMapping["song_review_option_id"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.SongReviewOptionID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to populate default values for song_review_option")
	}

CacheNoHooks:
	if !cached {
		songReviewOptionInsertCacheMut.Lock()
		songReviewOptionInsertCache[key] = cache
		songReviewOptionInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the SongReviewOption.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *SongReviewOption) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	songReviewOptionUpdateCacheMut.RLock()
	cache, cached := songReviewOptionUpdateCache[key]
	songReviewOptionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			songReviewOptionAllColumns,
			songReviewOptionPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("mysql: unable to update song_review_option, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `song_review_option` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, songReviewOptionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(songReviewOptionType, songReviewOptionMapping, append(wl, songReviewOptionPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "mysql: unable to update song_review_option row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by update for song_review_option")
	}

	if !cached {
		songReviewOptionUpdateCacheMut.Lock()
		songReviewOptionUpdateCache[key] = cache
		songReviewOptionUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q songReviewOptionQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update all for song_review_option")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to retrieve rows affected for song_review_option")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SongReviewOptionSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), songReviewOptionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `song_review_option` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, songReviewOptionPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update all in songReviewOption slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to retrieve rows affected all in update all songReviewOption")
	}
	return rowsAff, nil
}

var mySQLSongReviewOptionUniqueColumns = []string{
	"song_review_option_id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *SongReviewOption) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("mysql: no song_review_option provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(songReviewOptionColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLSongReviewOptionUniqueColumns, o)

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

	songReviewOptionUpsertCacheMut.RLock()
	cache, cached := songReviewOptionUpsertCache[key]
	songReviewOptionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			songReviewOptionAllColumns,
			songReviewOptionColumnsWithDefault,
			songReviewOptionColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			songReviewOptionAllColumns,
			songReviewOptionPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("mysql: unable to upsert song_review_option, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`song_review_option`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `song_review_option` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(songReviewOptionType, songReviewOptionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(songReviewOptionType, songReviewOptionMapping, ret)
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
		return errors.Wrap(err, "mysql: unable to upsert for song_review_option")
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

	o.SongReviewOptionID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == songReviewOptionMapping["song_review_option_id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(songReviewOptionType, songReviewOptionMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to retrieve unique values for song_review_option")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to populate default values for song_review_option")
	}

CacheNoHooks:
	if !cached {
		songReviewOptionUpsertCacheMut.Lock()
		songReviewOptionUpsertCache[key] = cache
		songReviewOptionUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single SongReviewOption record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SongReviewOption) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("mysql: no SongReviewOption provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), songReviewOptionPrimaryKeyMapping)
	sql := "DELETE FROM `song_review_option` WHERE `song_review_option_id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete from song_review_option")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by delete for song_review_option")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q songReviewOptionQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("mysql: no songReviewOptionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete all from song_review_option")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by deleteall for song_review_option")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SongReviewOptionSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(songReviewOptionBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), songReviewOptionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `song_review_option` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, songReviewOptionPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete all from songReviewOption slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by deleteall for song_review_option")
	}

	if len(songReviewOptionAfterDeleteHooks) != 0 {
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
func (o *SongReviewOption) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindSongReviewOption(ctx, exec, o.SongReviewOptionID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SongReviewOptionSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SongReviewOptionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), songReviewOptionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `song_review_option`.* FROM `song_review_option` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, songReviewOptionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to reload all in SongReviewOptionSlice")
	}

	*o = slice

	return nil
}

// SongReviewOptionExists checks if the SongReviewOption row exists.
func SongReviewOptionExists(ctx context.Context, exec boil.ContextExecutor, songReviewOptionID int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `song_review_option` where `song_review_option_id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, songReviewOptionID)
	}
	row := exec.QueryRowContext(ctx, sql, songReviewOptionID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "mysql: unable to check if song_review_option exists")
	}

	return exists, nil
}

// Exists checks if the SongReviewOption row exists.
func (o *SongReviewOption) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return SongReviewOptionExists(ctx, exec, o.SongReviewOptionID)
}
