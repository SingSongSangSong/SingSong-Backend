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

// MemberAction is an object representing the database table.
type MemberAction struct {
	MemberActionID int64       `boil:"member_action_id" json:"member_action_id" toml:"member_action_id" yaml:"member_action_id"`
	MemberID       int64       `boil:"member_id" json:"member_id" toml:"member_id" yaml:"member_id"`
	Gender         null.String `boil:"gender" json:"gender,omitempty" toml:"gender" yaml:"gender,omitempty"`
	Birthyear      null.Int    `boil:"birthyear" json:"birthyear,omitempty" toml:"birthyear" yaml:"birthyear,omitempty"`
	SongInfoID     int64       `boil:"song_info_id" json:"song_info_id" toml:"song_info_id" yaml:"song_info_id"`
	ActionType     string      `boil:"action_type" json:"action_type" toml:"action_type" yaml:"action_type"`
	ActionScore    float32     `boil:"action_score" json:"action_score" toml:"action_score" yaml:"action_score"`
	CreatedAt      null.Time   `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt      null.Time   `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
	DeletedAt      null.Time   `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`

	R *memberActionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L memberActionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var MemberActionColumns = struct {
	MemberActionID string
	MemberID       string
	Gender         string
	Birthyear      string
	SongInfoID     string
	ActionType     string
	ActionScore    string
	CreatedAt      string
	UpdatedAt      string
	DeletedAt      string
}{
	MemberActionID: "member_action_id",
	MemberID:       "member_id",
	Gender:         "gender",
	Birthyear:      "birthyear",
	SongInfoID:     "song_info_id",
	ActionType:     "action_type",
	ActionScore:    "action_score",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
}

var MemberActionTableColumns = struct {
	MemberActionID string
	MemberID       string
	Gender         string
	Birthyear      string
	SongInfoID     string
	ActionType     string
	ActionScore    string
	CreatedAt      string
	UpdatedAt      string
	DeletedAt      string
}{
	MemberActionID: "member_action.member_action_id",
	MemberID:       "member_action.member_id",
	Gender:         "member_action.gender",
	Birthyear:      "member_action.birthyear",
	SongInfoID:     "member_action.song_info_id",
	ActionType:     "member_action.action_type",
	ActionScore:    "member_action.action_score",
	CreatedAt:      "member_action.created_at",
	UpdatedAt:      "member_action.updated_at",
	DeletedAt:      "member_action.deleted_at",
}

// Generated where

type whereHelperfloat32 struct{ field string }

func (w whereHelperfloat32) EQ(x float32) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperfloat32) NEQ(x float32) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelperfloat32) LT(x float32) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperfloat32) LTE(x float32) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelperfloat32) GT(x float32) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperfloat32) GTE(x float32) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelperfloat32) IN(slice []float32) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperfloat32) NIN(slice []float32) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var MemberActionWhere = struct {
	MemberActionID whereHelperint64
	MemberID       whereHelperint64
	Gender         whereHelpernull_String
	Birthyear      whereHelpernull_Int
	SongInfoID     whereHelperint64
	ActionType     whereHelperstring
	ActionScore    whereHelperfloat32
	CreatedAt      whereHelpernull_Time
	UpdatedAt      whereHelpernull_Time
	DeletedAt      whereHelpernull_Time
}{
	MemberActionID: whereHelperint64{field: "`member_action`.`member_action_id`"},
	MemberID:       whereHelperint64{field: "`member_action`.`member_id`"},
	Gender:         whereHelpernull_String{field: "`member_action`.`gender`"},
	Birthyear:      whereHelpernull_Int{field: "`member_action`.`birthyear`"},
	SongInfoID:     whereHelperint64{field: "`member_action`.`song_info_id`"},
	ActionType:     whereHelperstring{field: "`member_action`.`action_type`"},
	ActionScore:    whereHelperfloat32{field: "`member_action`.`action_score`"},
	CreatedAt:      whereHelpernull_Time{field: "`member_action`.`created_at`"},
	UpdatedAt:      whereHelpernull_Time{field: "`member_action`.`updated_at`"},
	DeletedAt:      whereHelpernull_Time{field: "`member_action`.`deleted_at`"},
}

// MemberActionRels is where relationship names are stored.
var MemberActionRels = struct {
}{}

// memberActionR is where relationships are stored.
type memberActionR struct {
}

// NewStruct creates a new relationship struct
func (*memberActionR) NewStruct() *memberActionR {
	return &memberActionR{}
}

// memberActionL is where Load methods for each relationship are stored.
type memberActionL struct{}

var (
	memberActionAllColumns            = []string{"member_action_id", "member_id", "gender", "birthyear", "song_info_id", "action_type", "action_score", "created_at", "updated_at", "deleted_at"}
	memberActionColumnsWithoutDefault = []string{"member_id", "gender", "birthyear", "song_info_id", "action_type", "action_score", "deleted_at"}
	memberActionColumnsWithDefault    = []string{"member_action_id", "created_at", "updated_at"}
	memberActionPrimaryKeyColumns     = []string{"member_action_id"}
	memberActionGeneratedColumns      = []string{}
)

type (
	// MemberActionSlice is an alias for a slice of pointers to MemberAction.
	// This should almost always be used instead of []MemberAction.
	MemberActionSlice []*MemberAction
	// MemberActionHook is the signature for custom MemberAction hook methods
	MemberActionHook func(context.Context, boil.ContextExecutor, *MemberAction) error

	memberActionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	memberActionType                 = reflect.TypeOf(&MemberAction{})
	memberActionMapping              = queries.MakeStructMapping(memberActionType)
	memberActionPrimaryKeyMapping, _ = queries.BindMapping(memberActionType, memberActionMapping, memberActionPrimaryKeyColumns)
	memberActionInsertCacheMut       sync.RWMutex
	memberActionInsertCache          = make(map[string]insertCache)
	memberActionUpdateCacheMut       sync.RWMutex
	memberActionUpdateCache          = make(map[string]updateCache)
	memberActionUpsertCacheMut       sync.RWMutex
	memberActionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var memberActionAfterSelectHooks []MemberActionHook

var memberActionBeforeInsertHooks []MemberActionHook
var memberActionAfterInsertHooks []MemberActionHook

var memberActionBeforeUpdateHooks []MemberActionHook
var memberActionAfterUpdateHooks []MemberActionHook

var memberActionBeforeDeleteHooks []MemberActionHook
var memberActionAfterDeleteHooks []MemberActionHook

var memberActionBeforeUpsertHooks []MemberActionHook
var memberActionAfterUpsertHooks []MemberActionHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *MemberAction) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range memberActionAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *MemberAction) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range memberActionBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *MemberAction) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range memberActionAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *MemberAction) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range memberActionBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *MemberAction) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range memberActionAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *MemberAction) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range memberActionBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *MemberAction) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range memberActionAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *MemberAction) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range memberActionBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *MemberAction) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range memberActionAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddMemberActionHook registers your hook function for all future operations.
func AddMemberActionHook(hookPoint boil.HookPoint, memberActionHook MemberActionHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		memberActionAfterSelectHooks = append(memberActionAfterSelectHooks, memberActionHook)
	case boil.BeforeInsertHook:
		memberActionBeforeInsertHooks = append(memberActionBeforeInsertHooks, memberActionHook)
	case boil.AfterInsertHook:
		memberActionAfterInsertHooks = append(memberActionAfterInsertHooks, memberActionHook)
	case boil.BeforeUpdateHook:
		memberActionBeforeUpdateHooks = append(memberActionBeforeUpdateHooks, memberActionHook)
	case boil.AfterUpdateHook:
		memberActionAfterUpdateHooks = append(memberActionAfterUpdateHooks, memberActionHook)
	case boil.BeforeDeleteHook:
		memberActionBeforeDeleteHooks = append(memberActionBeforeDeleteHooks, memberActionHook)
	case boil.AfterDeleteHook:
		memberActionAfterDeleteHooks = append(memberActionAfterDeleteHooks, memberActionHook)
	case boil.BeforeUpsertHook:
		memberActionBeforeUpsertHooks = append(memberActionBeforeUpsertHooks, memberActionHook)
	case boil.AfterUpsertHook:
		memberActionAfterUpsertHooks = append(memberActionAfterUpsertHooks, memberActionHook)
	}
}

// One returns a single memberAction record from the query.
func (q memberActionQuery) One(ctx context.Context, exec boil.ContextExecutor) (*MemberAction, error) {
	o := &MemberAction{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysql: failed to execute a one query for member_action")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all MemberAction records from the query.
func (q memberActionQuery) All(ctx context.Context, exec boil.ContextExecutor) (MemberActionSlice, error) {
	var o []*MemberAction

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "mysql: failed to assign all query results to MemberAction slice")
	}

	if len(memberActionAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all MemberAction records in the query.
func (q memberActionQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to count member_action rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q memberActionQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "mysql: failed to check if member_action exists")
	}

	return count > 0, nil
}

// MemberActions retrieves all the records using an executor.
func MemberActions(mods ...qm.QueryMod) memberActionQuery {
	mods = append(mods, qm.From("`member_action`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`member_action`.*"})
	}

	return memberActionQuery{q}
}

// FindMemberAction retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindMemberAction(ctx context.Context, exec boil.ContextExecutor, memberActionID int64, selectCols ...string) (*MemberAction, error) {
	memberActionObj := &MemberAction{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `member_action` where `member_action_id`=?", sel,
	)

	q := queries.Raw(query, memberActionID)

	err := q.Bind(ctx, exec, memberActionObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysql: unable to select from member_action")
	}

	if err = memberActionObj.doAfterSelectHooks(ctx, exec); err != nil {
		return memberActionObj, err
	}

	return memberActionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *MemberAction) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("mysql: no member_action provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(memberActionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	memberActionInsertCacheMut.RLock()
	cache, cached := memberActionInsertCache[key]
	memberActionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			memberActionAllColumns,
			memberActionColumnsWithDefault,
			memberActionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(memberActionType, memberActionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(memberActionType, memberActionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `member_action` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `member_action` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `member_action` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, memberActionPrimaryKeyColumns))
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
		return errors.Wrap(err, "mysql: unable to insert into member_action")
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

	o.MemberActionID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == memberActionMapping["member_action_id"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.MemberActionID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to populate default values for member_action")
	}

CacheNoHooks:
	if !cached {
		memberActionInsertCacheMut.Lock()
		memberActionInsertCache[key] = cache
		memberActionInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the MemberAction.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *MemberAction) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	memberActionUpdateCacheMut.RLock()
	cache, cached := memberActionUpdateCache[key]
	memberActionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			memberActionAllColumns,
			memberActionPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("mysql: unable to update member_action, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `member_action` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, memberActionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(memberActionType, memberActionMapping, append(wl, memberActionPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "mysql: unable to update member_action row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by update for member_action")
	}

	if !cached {
		memberActionUpdateCacheMut.Lock()
		memberActionUpdateCache[key] = cache
		memberActionUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q memberActionQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update all for member_action")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to retrieve rows affected for member_action")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o MemberActionSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), memberActionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `member_action` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, memberActionPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update all in memberAction slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to retrieve rows affected all in update all memberAction")
	}
	return rowsAff, nil
}

var mySQLMemberActionUniqueColumns = []string{
	"member_action_id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *MemberAction) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("mysql: no member_action provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(memberActionColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLMemberActionUniqueColumns, o)

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

	memberActionUpsertCacheMut.RLock()
	cache, cached := memberActionUpsertCache[key]
	memberActionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			memberActionAllColumns,
			memberActionColumnsWithDefault,
			memberActionColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			memberActionAllColumns,
			memberActionPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("mysql: unable to upsert member_action, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`member_action`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `member_action` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(memberActionType, memberActionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(memberActionType, memberActionMapping, ret)
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
		return errors.Wrap(err, "mysql: unable to upsert for member_action")
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

	o.MemberActionID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == memberActionMapping["member_action_id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(memberActionType, memberActionMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to retrieve unique values for member_action")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to populate default values for member_action")
	}

CacheNoHooks:
	if !cached {
		memberActionUpsertCacheMut.Lock()
		memberActionUpsertCache[key] = cache
		memberActionUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single MemberAction record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *MemberAction) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("mysql: no MemberAction provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), memberActionPrimaryKeyMapping)
	sql := "DELETE FROM `member_action` WHERE `member_action_id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete from member_action")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by delete for member_action")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q memberActionQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("mysql: no memberActionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete all from member_action")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by deleteall for member_action")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o MemberActionSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(memberActionBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), memberActionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `member_action` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, memberActionPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete all from memberAction slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by deleteall for member_action")
	}

	if len(memberActionAfterDeleteHooks) != 0 {
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
func (o *MemberAction) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindMemberAction(ctx, exec, o.MemberActionID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *MemberActionSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := MemberActionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), memberActionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `member_action`.* FROM `member_action` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, memberActionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to reload all in MemberActionSlice")
	}

	*o = slice

	return nil
}

// MemberActionExists checks if the MemberAction row exists.
func MemberActionExists(ctx context.Context, exec boil.ContextExecutor, memberActionID int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `member_action` where `member_action_id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, memberActionID)
	}
	row := exec.QueryRowContext(ctx, sql, memberActionID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "mysql: unable to check if member_action exists")
	}

	return exists, nil
}

// Exists checks if the MemberAction row exists.
func (o *MemberAction) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return MemberActionExists(ctx, exec, o.MemberActionID)
}