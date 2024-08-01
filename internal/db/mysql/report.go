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

// Report is an object representing the database table.
type Report struct {
	ReportId     int64       `boil:"reportId" json:"reportId" toml:"reportId" yaml:"reportId"`
	CommentId    int64       `boil:"commentId" json:"commentId" toml:"commentId" yaml:"commentId"`
	ReporterId   int64       `boil:"reporterId" json:"reporterId" toml:"reporterId" yaml:"reporterId"`
	SubjectId    int64       `boil:"subjectId" json:"subjectId" toml:"subjectId" yaml:"subjectId"`
	ReportReason null.String `boil:"reportReason" json:"reportReason,omitempty" toml:"reportReason" yaml:"reportReason,omitempty"`
	CreatedAt    null.Time   `boil:"createdAt" json:"createdAt,omitempty" toml:"createdAt" yaml:"createdAt,omitempty"`
	UpdatedAt    null.Time   `boil:"updatedAt" json:"updatedAt,omitempty" toml:"updatedAt" yaml:"updatedAt,omitempty"`
	DeletedAt    null.Time   `boil:"deletedAt" json:"deletedAt,omitempty" toml:"deletedAt" yaml:"deletedAt,omitempty"`

	R *reportR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L reportL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ReportColumns = struct {
	ReportId     string
	CommentId    string
	ReporterId   string
	SubjectId    string
	ReportReason string
	CreatedAt    string
	UpdatedAt    string
	DeletedAt    string
}{
	ReportId:     "reportId",
	CommentId:    "commentId",
	ReporterId:   "reporterId",
	SubjectId:    "subjectId",
	ReportReason: "reportReason",
	CreatedAt:    "createdAt",
	UpdatedAt:    "updatedAt",
	DeletedAt:    "deletedAt",
}

var ReportTableColumns = struct {
	ReportId     string
	CommentId    string
	ReporterId   string
	SubjectId    string
	ReportReason string
	CreatedAt    string
	UpdatedAt    string
	DeletedAt    string
}{
	ReportId:     "report.reportId",
	CommentId:    "report.commentId",
	ReporterId:   "report.reporterId",
	SubjectId:    "report.subjectId",
	ReportReason: "report.reportReason",
	CreatedAt:    "report.createdAt",
	UpdatedAt:    "report.updatedAt",
	DeletedAt:    "report.deletedAt",
}

// Generated where

var ReportWhere = struct {
	ReportId     whereHelperint64
	CommentId    whereHelperint64
	ReporterId   whereHelperint64
	SubjectId    whereHelperint64
	ReportReason whereHelpernull_String
	CreatedAt    whereHelpernull_Time
	UpdatedAt    whereHelpernull_Time
	DeletedAt    whereHelpernull_Time
}{
	ReportId:     whereHelperint64{field: "`report`.`reportId`"},
	CommentId:    whereHelperint64{field: "`report`.`commentId`"},
	ReporterId:   whereHelperint64{field: "`report`.`reporterId`"},
	SubjectId:    whereHelperint64{field: "`report`.`subjectId`"},
	ReportReason: whereHelpernull_String{field: "`report`.`reportReason`"},
	CreatedAt:    whereHelpernull_Time{field: "`report`.`createdAt`"},
	UpdatedAt:    whereHelpernull_Time{field: "`report`.`updatedAt`"},
	DeletedAt:    whereHelpernull_Time{field: "`report`.`deletedAt`"},
}

// ReportRels is where relationship names are stored.
var ReportRels = struct {
}{}

// reportR is where relationships are stored.
type reportR struct {
}

// NewStruct creates a new relationship struct
func (*reportR) NewStruct() *reportR {
	return &reportR{}
}

// reportL is where Load methods for each relationship are stored.
type reportL struct{}

var (
	reportAllColumns            = []string{"reportId", "commentId", "reporterId", "subjectId", "reportReason", "createdAt", "updatedAt", "deletedAt"}
	reportColumnsWithoutDefault = []string{"commentId", "reporterId", "subjectId", "reportReason", "deletedAt"}
	reportColumnsWithDefault    = []string{"reportId", "createdAt", "updatedAt"}
	reportPrimaryKeyColumns     = []string{"reportId"}
	reportGeneratedColumns      = []string{}
)

type (
	// ReportSlice is an alias for a slice of pointers to Report.
	// This should almost always be used instead of []Report.
	ReportSlice []*Report
	// ReportHook is the signature for custom Report hook methods
	ReportHook func(context.Context, boil.ContextExecutor, *Report) error

	reportQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	reportType                 = reflect.TypeOf(&Report{})
	reportMapping              = queries.MakeStructMapping(reportType)
	reportPrimaryKeyMapping, _ = queries.BindMapping(reportType, reportMapping, reportPrimaryKeyColumns)
	reportInsertCacheMut       sync.RWMutex
	reportInsertCache          = make(map[string]insertCache)
	reportUpdateCacheMut       sync.RWMutex
	reportUpdateCache          = make(map[string]updateCache)
	reportUpsertCacheMut       sync.RWMutex
	reportUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var reportAfterSelectMu sync.Mutex
var reportAfterSelectHooks []ReportHook

var reportBeforeInsertMu sync.Mutex
var reportBeforeInsertHooks []ReportHook
var reportAfterInsertMu sync.Mutex
var reportAfterInsertHooks []ReportHook

var reportBeforeUpdateMu sync.Mutex
var reportBeforeUpdateHooks []ReportHook
var reportAfterUpdateMu sync.Mutex
var reportAfterUpdateHooks []ReportHook

var reportBeforeDeleteMu sync.Mutex
var reportBeforeDeleteHooks []ReportHook
var reportAfterDeleteMu sync.Mutex
var reportAfterDeleteHooks []ReportHook

var reportBeforeUpsertMu sync.Mutex
var reportBeforeUpsertHooks []ReportHook
var reportAfterUpsertMu sync.Mutex
var reportAfterUpsertHooks []ReportHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Report) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Report) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Report) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Report) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Report) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Report) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Report) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Report) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Report) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddReportHook registers your hook function for all future operations.
func AddReportHook(hookPoint boil.HookPoint, reportHook ReportHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		reportAfterSelectMu.Lock()
		reportAfterSelectHooks = append(reportAfterSelectHooks, reportHook)
		reportAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		reportBeforeInsertMu.Lock()
		reportBeforeInsertHooks = append(reportBeforeInsertHooks, reportHook)
		reportBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		reportAfterInsertMu.Lock()
		reportAfterInsertHooks = append(reportAfterInsertHooks, reportHook)
		reportAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		reportBeforeUpdateMu.Lock()
		reportBeforeUpdateHooks = append(reportBeforeUpdateHooks, reportHook)
		reportBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		reportAfterUpdateMu.Lock()
		reportAfterUpdateHooks = append(reportAfterUpdateHooks, reportHook)
		reportAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		reportBeforeDeleteMu.Lock()
		reportBeforeDeleteHooks = append(reportBeforeDeleteHooks, reportHook)
		reportBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		reportAfterDeleteMu.Lock()
		reportAfterDeleteHooks = append(reportAfterDeleteHooks, reportHook)
		reportAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		reportBeforeUpsertMu.Lock()
		reportBeforeUpsertHooks = append(reportBeforeUpsertHooks, reportHook)
		reportBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		reportAfterUpsertMu.Lock()
		reportAfterUpsertHooks = append(reportAfterUpsertHooks, reportHook)
		reportAfterUpsertMu.Unlock()
	}
}

// One returns a single report record from the query.
func (q reportQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Report, error) {
	o := &Report{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysql: failed to execute a one query for report")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Report records from the query.
func (q reportQuery) All(ctx context.Context, exec boil.ContextExecutor) (ReportSlice, error) {
	var o []*Report

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "mysql: failed to assign all query results to Report slice")
	}

	if len(reportAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Report records in the query.
func (q reportQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to count report rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q reportQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "mysql: failed to check if report exists")
	}

	return count > 0, nil
}

// Reports retrieves all the records using an executor.
func Reports(mods ...qm.QueryMod) reportQuery {
	mods = append(mods, qm.From("`report`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`report`.*"})
	}

	return reportQuery{q}
}

// FindReport retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindReport(ctx context.Context, exec boil.ContextExecutor, reportId int64, selectCols ...string) (*Report, error) {
	reportObj := &Report{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `report` where `reportId`=?", sel,
	)

	q := queries.Raw(query, reportId)

	err := q.Bind(ctx, exec, reportObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysql: unable to select from report")
	}

	if err = reportObj.doAfterSelectHooks(ctx, exec); err != nil {
		return reportObj, err
	}

	return reportObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Report) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("mysql: no report provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(reportColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	reportInsertCacheMut.RLock()
	cache, cached := reportInsertCache[key]
	reportInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			reportAllColumns,
			reportColumnsWithDefault,
			reportColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(reportType, reportMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(reportType, reportMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `report` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `report` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `report` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, reportPrimaryKeyColumns))
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
		return errors.Wrap(err, "mysql: unable to insert into report")
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

	o.ReportId = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == reportMapping["reportId"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ReportId,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to populate default values for report")
	}

CacheNoHooks:
	if !cached {
		reportInsertCacheMut.Lock()
		reportInsertCache[key] = cache
		reportInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Report.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Report) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	reportUpdateCacheMut.RLock()
	cache, cached := reportUpdateCache[key]
	reportUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			reportAllColumns,
			reportPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("mysql: unable to update report, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `report` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, reportPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(reportType, reportMapping, append(wl, reportPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "mysql: unable to update report row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by update for report")
	}

	if !cached {
		reportUpdateCacheMut.Lock()
		reportUpdateCache[key] = cache
		reportUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q reportQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update all for report")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to retrieve rows affected for report")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ReportSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reportPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `report` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, reportPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update all in report slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to retrieve rows affected all in update all report")
	}
	return rowsAff, nil
}

var mySQLReportUniqueColumns = []string{
	"reportId",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Report) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("mysql: no report provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(reportColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLReportUniqueColumns, o)

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

	reportUpsertCacheMut.RLock()
	cache, cached := reportUpsertCache[key]
	reportUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			reportAllColumns,
			reportColumnsWithDefault,
			reportColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			reportAllColumns,
			reportPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("mysql: unable to upsert report, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`report`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `report` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(reportType, reportMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(reportType, reportMapping, ret)
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
		return errors.Wrap(err, "mysql: unable to upsert for report")
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

	o.ReportId = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == reportMapping["reportId"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(reportType, reportMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to retrieve unique values for report")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to populate default values for report")
	}

CacheNoHooks:
	if !cached {
		reportUpsertCacheMut.Lock()
		reportUpsertCache[key] = cache
		reportUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Report record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Report) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("mysql: no Report provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), reportPrimaryKeyMapping)
	sql := "DELETE FROM `report` WHERE `reportId`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete from report")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by delete for report")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q reportQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("mysql: no reportQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete all from report")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by deleteall for report")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ReportSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(reportBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reportPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `report` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, reportPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete all from report slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by deleteall for report")
	}

	if len(reportAfterDeleteHooks) != 0 {
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
func (o *Report) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindReport(ctx, exec, o.ReportId)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ReportSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ReportSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reportPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `report`.* FROM `report` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, reportPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to reload all in ReportSlice")
	}

	*o = slice

	return nil
}

// ReportExists checks if the Report row exists.
func ReportExists(ctx context.Context, exec boil.ContextExecutor, reportId int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `report` where `reportId`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, reportId)
	}
	row := exec.QueryRowContext(ctx, sql, reportId)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "mysql: unable to check if report exists")
	}

	return exists, nil
}

// Exists checks if the Report row exists.
func (o *Report) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return ReportExists(ctx, exec, o.ReportId)
}
