// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// SongReview is an object representing the database table.
type SongReview struct {
	SongReviewId       int64       `boil:"songReviewId" json:"songReviewId" toml:"songReviewId" yaml:"songReviewId"`
	SongId             int64       `boil:"songId" json:"songId" toml:"songId" yaml:"songId"`
	MemberId           int64       `boil:"memberId" json:"memberId" toml:"memberId" yaml:"memberId"`
	Gender             null.String `boil:"gender" json:"gender,omitempty" toml:"gender" yaml:"gender,omitempty"`
	Birthyear          null.Int    `boil:"birthyear" json:"birthyear,omitempty" toml:"birthyear" yaml:"birthyear,omitempty"`
	SongReviewOptionId int64       `boil:"songReviewOptionId" json:"songReviewOptionId" toml:"songReviewOptionId" yaml:"songReviewOptionId"`

	R *songReviewR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L songReviewL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SongReviewColumns = struct {
	SongReviewId       string
	SongId             string
	MemberId           string
	Gender             string
	Birthyear          string
	SongReviewOptionId string
}{
	SongReviewId:       "songReviewId",
	SongId:             "songId",
	MemberId:           "memberId",
	Gender:             "gender",
	Birthyear:          "birthyear",
	SongReviewOptionId: "songReviewOptionId",
}

var SongReviewTableColumns = struct {
	SongReviewId       string
	SongId             string
	MemberId           string
	Gender             string
	Birthyear          string
	SongReviewOptionId string
}{
	SongReviewId:       "songReview.songReviewId",
	SongId:             "songReview.songId",
	MemberId:           "songReview.memberId",
	Gender:             "songReview.gender",
	Birthyear:          "songReview.birthyear",
	SongReviewOptionId: "songReview.songReviewOptionId",
}

// Generated where

var SongReviewWhere = struct {
	SongReviewId       whereHelperint64
	SongId             whereHelperint64
	MemberId           whereHelperint64
	Gender             whereHelpernull_String
	Birthyear          whereHelpernull_Int
	SongReviewOptionId whereHelperint64
}{
	SongReviewId:       whereHelperint64{field: "`songReview`.`songReviewId`"},
	SongId:             whereHelperint64{field: "`songReview`.`songId`"},
	MemberId:           whereHelperint64{field: "`songReview`.`memberId`"},
	Gender:             whereHelpernull_String{field: "`songReview`.`gender`"},
	Birthyear:          whereHelpernull_Int{field: "`songReview`.`birthyear`"},
	SongReviewOptionId: whereHelperint64{field: "`songReview`.`songReviewOptionId`"},
}

// SongReviewRels is where relationship names are stored.
var SongReviewRels = struct {
}{}

// songReviewR is where relationships are stored.
type songReviewR struct {
}

// NewStruct creates a new relationship struct
func (*songReviewR) NewStruct() *songReviewR {
	return &songReviewR{}
}

// songReviewL is where Load methods for each relationship are stored.
type songReviewL struct{}

var (
	songReviewAllColumns            = []string{"songReviewId", "songId", "memberId", "gender", "birthyear", "songReviewOptionId"}
	songReviewColumnsWithoutDefault = []string{"songId", "memberId", "gender", "birthyear", "songReviewOptionId"}
	songReviewColumnsWithDefault    = []string{"songReviewId"}
	songReviewPrimaryKeyColumns     = []string{"songReviewId"}
	songReviewGeneratedColumns      = []string{}
)

type (
	// SongReviewSlice is an alias for a slice of pointers to SongReview.
	// This should almost always be used instead of []SongReview.
	SongReviewSlice []*SongReview
	// SongReviewHook is the signature for custom SongReview hook methods
	SongReviewHook func(context.Context, boil.ContextExecutor, *SongReview) error

	songReviewQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	songReviewType                 = reflect.TypeOf(&SongReview{})
	songReviewMapping              = queries.MakeStructMapping(songReviewType)
	songReviewPrimaryKeyMapping, _ = queries.BindMapping(songReviewType, songReviewMapping, songReviewPrimaryKeyColumns)
	songReviewInsertCacheMut       sync.RWMutex
	songReviewInsertCache          = make(map[string]insertCache)
	songReviewUpdateCacheMut       sync.RWMutex
	songReviewUpdateCache          = make(map[string]updateCache)
	songReviewUpsertCacheMut       sync.RWMutex
	songReviewUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var songReviewAfterSelectMu sync.Mutex
var songReviewAfterSelectHooks []SongReviewHook

var songReviewBeforeInsertMu sync.Mutex
var songReviewBeforeInsertHooks []SongReviewHook
var songReviewAfterInsertMu sync.Mutex
var songReviewAfterInsertHooks []SongReviewHook

var songReviewBeforeUpdateMu sync.Mutex
var songReviewBeforeUpdateHooks []SongReviewHook
var songReviewAfterUpdateMu sync.Mutex
var songReviewAfterUpdateHooks []SongReviewHook

var songReviewBeforeDeleteMu sync.Mutex
var songReviewBeforeDeleteHooks []SongReviewHook
var songReviewAfterDeleteMu sync.Mutex
var songReviewAfterDeleteHooks []SongReviewHook

var songReviewBeforeUpsertMu sync.Mutex
var songReviewBeforeUpsertHooks []SongReviewHook
var songReviewAfterUpsertMu sync.Mutex
var songReviewAfterUpsertHooks []SongReviewHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *SongReview) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *SongReview) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *SongReview) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *SongReview) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *SongReview) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *SongReview) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *SongReview) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *SongReview) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *SongReview) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range songReviewAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSongReviewHook registers your hook function for all future operations.
func AddSongReviewHook(hookPoint boil.HookPoint, songReviewHook SongReviewHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		songReviewAfterSelectMu.Lock()
		songReviewAfterSelectHooks = append(songReviewAfterSelectHooks, songReviewHook)
		songReviewAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		songReviewBeforeInsertMu.Lock()
		songReviewBeforeInsertHooks = append(songReviewBeforeInsertHooks, songReviewHook)
		songReviewBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		songReviewAfterInsertMu.Lock()
		songReviewAfterInsertHooks = append(songReviewAfterInsertHooks, songReviewHook)
		songReviewAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		songReviewBeforeUpdateMu.Lock()
		songReviewBeforeUpdateHooks = append(songReviewBeforeUpdateHooks, songReviewHook)
		songReviewBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		songReviewAfterUpdateMu.Lock()
		songReviewAfterUpdateHooks = append(songReviewAfterUpdateHooks, songReviewHook)
		songReviewAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		songReviewBeforeDeleteMu.Lock()
		songReviewBeforeDeleteHooks = append(songReviewBeforeDeleteHooks, songReviewHook)
		songReviewBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		songReviewAfterDeleteMu.Lock()
		songReviewAfterDeleteHooks = append(songReviewAfterDeleteHooks, songReviewHook)
		songReviewAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		songReviewBeforeUpsertMu.Lock()
		songReviewBeforeUpsertHooks = append(songReviewBeforeUpsertHooks, songReviewHook)
		songReviewBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		songReviewAfterUpsertMu.Lock()
		songReviewAfterUpsertHooks = append(songReviewAfterUpsertHooks, songReviewHook)
		songReviewAfterUpsertMu.Unlock()
	}
}

// One returns a single songReview record from the query.
func (q songReviewQuery) One(ctx context.Context, exec boil.ContextExecutor) (*SongReview, error) {
	o := &SongReview{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysql: failed to execute a one query for songReview")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all SongReview records from the query.
func (q songReviewQuery) All(ctx context.Context, exec boil.ContextExecutor) (SongReviewSlice, error) {
	var o []*SongReview

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "mysql: failed to assign all query results to SongReview slice")
	}

	if len(songReviewAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all SongReview records in the query.
func (q songReviewQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to count songReview rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q songReviewQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "mysql: failed to check if songReview exists")
	}

	return count > 0, nil
}

// SongReviews retrieves all the records using an executor.
func SongReviews(mods ...qm.QueryMod) songReviewQuery {
	mods = append(mods, qm.From("`songReview`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`songReview`.*"})
	}

	return songReviewQuery{q}
}

// FindSongReview retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSongReview(ctx context.Context, exec boil.ContextExecutor, songReviewId int64, selectCols ...string) (*SongReview, error) {
	songReviewObj := &SongReview{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `songReview` where `songReviewId`=?", sel,
	)

	q := queries.Raw(query, songReviewId)

	err := q.Bind(ctx, exec, songReviewObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysql: unable to select from songReview")
	}

	if err = songReviewObj.doAfterSelectHooks(ctx, exec); err != nil {
		return songReviewObj, err
	}

	return songReviewObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *SongReview) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("mysql: no songReview provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(songReviewColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	songReviewInsertCacheMut.RLock()
	cache, cached := songReviewInsertCache[key]
	songReviewInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			songReviewAllColumns,
			songReviewColumnsWithDefault,
			songReviewColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(songReviewType, songReviewMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(songReviewType, songReviewMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `songReview` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `songReview` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `songReview` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, songReviewPrimaryKeyColumns))
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
		return errors.Wrap(err, "mysql: unable to insert into songReview")
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

	o.SongReviewId = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == songReviewMapping["songReviewId"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.SongReviewId,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to populate default values for songReview")
	}

CacheNoHooks:
	if !cached {
		songReviewInsertCacheMut.Lock()
		songReviewInsertCache[key] = cache
		songReviewInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the SongReview.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *SongReview) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	songReviewUpdateCacheMut.RLock()
	cache, cached := songReviewUpdateCache[key]
	songReviewUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			songReviewAllColumns,
			songReviewPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("mysql: unable to update songReview, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `songReview` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, songReviewPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(songReviewType, songReviewMapping, append(wl, songReviewPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "mysql: unable to update songReview row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by update for songReview")
	}

	if !cached {
		songReviewUpdateCacheMut.Lock()
		songReviewUpdateCache[key] = cache
		songReviewUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q songReviewQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update all for songReview")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to retrieve rows affected for songReview")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SongReviewSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), songReviewPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `songReview` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, songReviewPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update all in songReview slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to retrieve rows affected all in update all songReview")
	}
	return rowsAff, nil
}

var mySQLSongReviewUniqueColumns = []string{
	"songReviewId",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *SongReview) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("mysql: no songReview provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(songReviewColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLSongReviewUniqueColumns, o)

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

	songReviewUpsertCacheMut.RLock()
	cache, cached := songReviewUpsertCache[key]
	songReviewUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			songReviewAllColumns,
			songReviewColumnsWithDefault,
			songReviewColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			songReviewAllColumns,
			songReviewPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("mysql: unable to upsert songReview, could not build update column list")
		}

		ret := strmangle.SetComplement(songReviewAllColumns, strmangle.SetIntersect(insert, update))

		cache.query = buildUpsertQueryMySQL(dialect, "`songReview`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `songReview` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(songReviewType, songReviewMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(songReviewType, songReviewMapping, ret)
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
		return errors.Wrap(err, "mysql: unable to upsert for songReview")
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

	o.SongReviewId = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == songReviewMapping["songReviewId"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(songReviewType, songReviewMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to retrieve unique values for songReview")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to populate default values for songReview")
	}

CacheNoHooks:
	if !cached {
		songReviewUpsertCacheMut.Lock()
		songReviewUpsertCache[key] = cache
		songReviewUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single SongReview record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SongReview) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("mysql: no SongReview provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), songReviewPrimaryKeyMapping)
	sql := "DELETE FROM `songReview` WHERE `songReviewId`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete from songReview")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by delete for songReview")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q songReviewQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("mysql: no songReviewQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete all from songReview")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by deleteall for songReview")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SongReviewSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(songReviewBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), songReviewPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `songReview` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, songReviewPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete all from songReview slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by deleteall for songReview")
	}

	if len(songReviewAfterDeleteHooks) != 0 {
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
func (o *SongReview) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindSongReview(ctx, exec, o.SongReviewId)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SongReviewSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SongReviewSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), songReviewPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `songReview`.* FROM `songReview` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, songReviewPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to reload all in SongReviewSlice")
	}

	*o = slice

	return nil
}

// SongReviewExists checks if the SongReview row exists.
func SongReviewExists(ctx context.Context, exec boil.ContextExecutor, songReviewId int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `songReview` where `songReviewId`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, songReviewId)
	}
	row := exec.QueryRowContext(ctx, sql, songReviewId)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "mysql: unable to check if songReview exists")
	}

	return exists, nil
}

// Exists checks if the SongReview row exists.
func (o *SongReview) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return SongReviewExists(ctx, exec, o.SongReviewId)
}
