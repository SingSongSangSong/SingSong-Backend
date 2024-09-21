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

// AppVersion is an object representing the database table.
type AppVersion struct {
	AppVersionID int64     `boil:"app_version_id" json:"app_version_id" toml:"app_version_id" yaml:"app_version_id"`
	Platform     string    `boil:"platform" json:"platform" toml:"platform" yaml:"platform"`
	Version      string    `boil:"version" json:"version" toml:"version" yaml:"version"`
	ForceUpdate  bool      `boil:"force_update" json:"force_update" toml:"force_update" yaml:"force_update"`
	CreatedAt    null.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt    null.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
	DeletedAt    null.Time `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`

	R *appVersionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L appVersionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var AppVersionColumns = struct {
	AppVersionID string
	Platform     string
	Version      string
	ForceUpdate  string
	CreatedAt    string
	UpdatedAt    string
	DeletedAt    string
}{
	AppVersionID: "app_version_id",
	Platform:     "platform",
	Version:      "version",
	ForceUpdate:  "force_update",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

var AppVersionTableColumns = struct {
	AppVersionID string
	Platform     string
	Version      string
	ForceUpdate  string
	CreatedAt    string
	UpdatedAt    string
	DeletedAt    string
}{
	AppVersionID: "app_version.app_version_id",
	Platform:     "app_version.platform",
	Version:      "app_version.version",
	ForceUpdate:  "app_version.force_update",
	CreatedAt:    "app_version.created_at",
	UpdatedAt:    "app_version.updated_at",
	DeletedAt:    "app_version.deleted_at",
}

// Generated where

type whereHelperint64 struct{ field string }

func (w whereHelperint64) EQ(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint64) NEQ(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint64) LT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint64) LTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint64) GT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint64) GTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint64) IN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint64) NIN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod   { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod   { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod   { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) LIKE(x string) qm.QueryMod  { return qm.Where(w.field+" LIKE ?", x) }
func (w whereHelperstring) NLIKE(x string) qm.QueryMod { return qm.Where(w.field+" NOT LIKE ?", x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelperbool struct{ field string }

func (w whereHelperbool) EQ(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperbool) NEQ(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperbool) LT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperbool) LTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperbool) GT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperbool) GTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

type whereHelpernull_Time struct{ field string }

func (w whereHelpernull_Time) EQ(x null.Time) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Time) NEQ(x null.Time) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Time) LT(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Time) LTE(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Time) GT(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Time) GTE(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpernull_Time) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Time) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

var AppVersionWhere = struct {
	AppVersionID whereHelperint64
	Platform     whereHelperstring
	Version      whereHelperstring
	ForceUpdate  whereHelperbool
	CreatedAt    whereHelpernull_Time
	UpdatedAt    whereHelpernull_Time
	DeletedAt    whereHelpernull_Time
}{
	AppVersionID: whereHelperint64{field: "`app_version`.`app_version_id`"},
	Platform:     whereHelperstring{field: "`app_version`.`platform`"},
	Version:      whereHelperstring{field: "`app_version`.`version`"},
	ForceUpdate:  whereHelperbool{field: "`app_version`.`force_update`"},
	CreatedAt:    whereHelpernull_Time{field: "`app_version`.`created_at`"},
	UpdatedAt:    whereHelpernull_Time{field: "`app_version`.`updated_at`"},
	DeletedAt:    whereHelpernull_Time{field: "`app_version`.`deleted_at`"},
}

// AppVersionRels is where relationship names are stored.
var AppVersionRels = struct {
}{}

// appVersionR is where relationships are stored.
type appVersionR struct {
}

// NewStruct creates a new relationship struct
func (*appVersionR) NewStruct() *appVersionR {
	return &appVersionR{}
}

// appVersionL is where Load methods for each relationship are stored.
type appVersionL struct{}

var (
	appVersionAllColumns            = []string{"app_version_id", "platform", "version", "force_update", "created_at", "updated_at", "deleted_at"}
	appVersionColumnsWithoutDefault = []string{"platform", "version", "deleted_at"}
	appVersionColumnsWithDefault    = []string{"app_version_id", "force_update", "created_at", "updated_at"}
	appVersionPrimaryKeyColumns     = []string{"app_version_id"}
	appVersionGeneratedColumns      = []string{}
)

type (
	// AppVersionSlice is an alias for a slice of pointers to AppVersion.
	// This should almost always be used instead of []AppVersion.
	AppVersionSlice []*AppVersion
	// AppVersionHook is the signature for custom AppVersion hook methods
	AppVersionHook func(context.Context, boil.ContextExecutor, *AppVersion) error

	appVersionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	appVersionType                 = reflect.TypeOf(&AppVersion{})
	appVersionMapping              = queries.MakeStructMapping(appVersionType)
	appVersionPrimaryKeyMapping, _ = queries.BindMapping(appVersionType, appVersionMapping, appVersionPrimaryKeyColumns)
	appVersionInsertCacheMut       sync.RWMutex
	appVersionInsertCache          = make(map[string]insertCache)
	appVersionUpdateCacheMut       sync.RWMutex
	appVersionUpdateCache          = make(map[string]updateCache)
	appVersionUpsertCacheMut       sync.RWMutex
	appVersionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var appVersionAfterSelectHooks []AppVersionHook

var appVersionBeforeInsertHooks []AppVersionHook
var appVersionAfterInsertHooks []AppVersionHook

var appVersionBeforeUpdateHooks []AppVersionHook
var appVersionAfterUpdateHooks []AppVersionHook

var appVersionBeforeDeleteHooks []AppVersionHook
var appVersionAfterDeleteHooks []AppVersionHook

var appVersionBeforeUpsertHooks []AppVersionHook
var appVersionAfterUpsertHooks []AppVersionHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *AppVersion) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range appVersionAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *AppVersion) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range appVersionBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *AppVersion) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range appVersionAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *AppVersion) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range appVersionBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *AppVersion) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range appVersionAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *AppVersion) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range appVersionBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *AppVersion) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range appVersionAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *AppVersion) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range appVersionBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *AppVersion) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range appVersionAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddAppVersionHook registers your hook function for all future operations.
func AddAppVersionHook(hookPoint boil.HookPoint, appVersionHook AppVersionHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		appVersionAfterSelectHooks = append(appVersionAfterSelectHooks, appVersionHook)
	case boil.BeforeInsertHook:
		appVersionBeforeInsertHooks = append(appVersionBeforeInsertHooks, appVersionHook)
	case boil.AfterInsertHook:
		appVersionAfterInsertHooks = append(appVersionAfterInsertHooks, appVersionHook)
	case boil.BeforeUpdateHook:
		appVersionBeforeUpdateHooks = append(appVersionBeforeUpdateHooks, appVersionHook)
	case boil.AfterUpdateHook:
		appVersionAfterUpdateHooks = append(appVersionAfterUpdateHooks, appVersionHook)
	case boil.BeforeDeleteHook:
		appVersionBeforeDeleteHooks = append(appVersionBeforeDeleteHooks, appVersionHook)
	case boil.AfterDeleteHook:
		appVersionAfterDeleteHooks = append(appVersionAfterDeleteHooks, appVersionHook)
	case boil.BeforeUpsertHook:
		appVersionBeforeUpsertHooks = append(appVersionBeforeUpsertHooks, appVersionHook)
	case boil.AfterUpsertHook:
		appVersionAfterUpsertHooks = append(appVersionAfterUpsertHooks, appVersionHook)
	}
}

// One returns a single appVersion record from the query.
func (q appVersionQuery) One(ctx context.Context, exec boil.ContextExecutor) (*AppVersion, error) {
	o := &AppVersion{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysql: failed to execute a one query for app_version")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all AppVersion records from the query.
func (q appVersionQuery) All(ctx context.Context, exec boil.ContextExecutor) (AppVersionSlice, error) {
	var o []*AppVersion

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "mysql: failed to assign all query results to AppVersion slice")
	}

	if len(appVersionAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all AppVersion records in the query.
func (q appVersionQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to count app_version rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q appVersionQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "mysql: failed to check if app_version exists")
	}

	return count > 0, nil
}

// AppVersions retrieves all the records using an executor.
func AppVersions(mods ...qm.QueryMod) appVersionQuery {
	mods = append(mods, qm.From("`app_version`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`app_version`.*"})
	}

	return appVersionQuery{q}
}

// FindAppVersion retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindAppVersion(ctx context.Context, exec boil.ContextExecutor, appVersionID int64, selectCols ...string) (*AppVersion, error) {
	appVersionObj := &AppVersion{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `app_version` where `app_version_id`=?", sel,
	)

	q := queries.Raw(query, appVersionID)

	err := q.Bind(ctx, exec, appVersionObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysql: unable to select from app_version")
	}

	if err = appVersionObj.doAfterSelectHooks(ctx, exec); err != nil {
		return appVersionObj, err
	}

	return appVersionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *AppVersion) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("mysql: no app_version provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(appVersionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	appVersionInsertCacheMut.RLock()
	cache, cached := appVersionInsertCache[key]
	appVersionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			appVersionAllColumns,
			appVersionColumnsWithDefault,
			appVersionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(appVersionType, appVersionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(appVersionType, appVersionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `app_version` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `app_version` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `app_version` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, appVersionPrimaryKeyColumns))
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
		return errors.Wrap(err, "mysql: unable to insert into app_version")
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

	o.AppVersionID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == appVersionMapping["app_version_id"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.AppVersionID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to populate default values for app_version")
	}

CacheNoHooks:
	if !cached {
		appVersionInsertCacheMut.Lock()
		appVersionInsertCache[key] = cache
		appVersionInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the AppVersion.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *AppVersion) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	appVersionUpdateCacheMut.RLock()
	cache, cached := appVersionUpdateCache[key]
	appVersionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			appVersionAllColumns,
			appVersionPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("mysql: unable to update app_version, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `app_version` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, appVersionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(appVersionType, appVersionMapping, append(wl, appVersionPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "mysql: unable to update app_version row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by update for app_version")
	}

	if !cached {
		appVersionUpdateCacheMut.Lock()
		appVersionUpdateCache[key] = cache
		appVersionUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q appVersionQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update all for app_version")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to retrieve rows affected for app_version")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o AppVersionSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), appVersionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `app_version` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, appVersionPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update all in appVersion slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to retrieve rows affected all in update all appVersion")
	}
	return rowsAff, nil
}

var mySQLAppVersionUniqueColumns = []string{
	"app_version_id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *AppVersion) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("mysql: no app_version provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(appVersionColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLAppVersionUniqueColumns, o)

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

	appVersionUpsertCacheMut.RLock()
	cache, cached := appVersionUpsertCache[key]
	appVersionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			appVersionAllColumns,
			appVersionColumnsWithDefault,
			appVersionColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			appVersionAllColumns,
			appVersionPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("mysql: unable to upsert app_version, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`app_version`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `app_version` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(appVersionType, appVersionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(appVersionType, appVersionMapping, ret)
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
		return errors.Wrap(err, "mysql: unable to upsert for app_version")
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

	o.AppVersionID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == appVersionMapping["app_version_id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(appVersionType, appVersionMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to retrieve unique values for app_version")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to populate default values for app_version")
	}

CacheNoHooks:
	if !cached {
		appVersionUpsertCacheMut.Lock()
		appVersionUpsertCache[key] = cache
		appVersionUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single AppVersion record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *AppVersion) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("mysql: no AppVersion provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), appVersionPrimaryKeyMapping)
	sql := "DELETE FROM `app_version` WHERE `app_version_id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete from app_version")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by delete for app_version")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q appVersionQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("mysql: no appVersionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete all from app_version")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by deleteall for app_version")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o AppVersionSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(appVersionBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), appVersionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `app_version` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, appVersionPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete all from appVersion slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by deleteall for app_version")
	}

	if len(appVersionAfterDeleteHooks) != 0 {
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
func (o *AppVersion) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindAppVersion(ctx, exec, o.AppVersionID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AppVersionSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := AppVersionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), appVersionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `app_version`.* FROM `app_version` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, appVersionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to reload all in AppVersionSlice")
	}

	*o = slice

	return nil
}

// AppVersionExists checks if the AppVersion row exists.
func AppVersionExists(ctx context.Context, exec boil.ContextExecutor, appVersionID int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `app_version` where `app_version_id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, appVersionID)
	}
	row := exec.QueryRowContext(ctx, sql, appVersionID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "mysql: unable to check if app_version exists")
	}

	return exists, nil
}

// Exists checks if the AppVersion row exists.
func (o *AppVersion) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return AppVersionExists(ctx, exec, o.AppVersionID)
}
