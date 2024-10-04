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

// TempSongInfo is an object representing the database table.
type TempSongInfo struct {
	TempSongInfoID int64       `boil:"temp_song_info_id" json:"temp_song_info_id" toml:"temp_song_info_id" yaml:"temp_song_info_id"`
	SongInfoID     int64       `boil:"song_info_id" json:"song_info_id" toml:"song_info_id" yaml:"song_info_id"`
	SongNumber     int         `boil:"song_number" json:"song_number" toml:"song_number" yaml:"song_number"`
	SongName       null.String `boil:"song_name" json:"song_name,omitempty" toml:"song_name" yaml:"song_name,omitempty"`
	ArtistName     null.String `boil:"artist_name" json:"artist_name,omitempty" toml:"artist_name" yaml:"artist_name,omitempty"`
	ArtistType     null.String `boil:"artist_type" json:"artist_type,omitempty" toml:"artist_type" yaml:"artist_type,omitempty"`
	IsMR           null.Bool   `boil:"is_mr" json:"is_mr,omitempty" toml:"is_mr" yaml:"is_mr,omitempty"`
	IsChosen22000  null.Bool   `boil:"is_chosen_22000" json:"is_chosen_22000,omitempty" toml:"is_chosen_22000" yaml:"is_chosen_22000,omitempty"`
	Country        null.String `boil:"country" json:"country,omitempty" toml:"country" yaml:"country,omitempty"`
	Album          null.String `boil:"album" json:"album,omitempty" toml:"album" yaml:"album,omitempty"`
	OctaveCrawling null.String `boil:"octave_crawling" json:"octave_crawling,omitempty" toml:"octave_crawling" yaml:"octave_crawling,omitempty"`
	OctaveLibrosa  null.String `boil:"octave_librosa" json:"octave_librosa,omitempty" toml:"octave_librosa" yaml:"octave_librosa,omitempty"`
	VideoLink      null.String `boil:"video_link" json:"video_link,omitempty" toml:"video_link" yaml:"video_link,omitempty"`
	MelonSongID    null.String `boil:"melon_song_id" json:"melon_song_id,omitempty" toml:"melon_song_id" yaml:"melon_song_id,omitempty"`
	IsLive         null.Bool   `boil:"is_live" json:"is_live,omitempty" toml:"is_live" yaml:"is_live,omitempty"`
	CreatedAt      null.Time   `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt      null.Time   `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
	DeletedAt      null.Time   `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`
	HighTag        null.Bool   `boil:"high_tag" json:"high_tag,omitempty" toml:"high_tag" yaml:"high_tag,omitempty"`
	LowTag         null.Bool   `boil:"low_tag" json:"low_tag,omitempty" toml:"low_tag" yaml:"low_tag,omitempty"`
	AudioFileURL   null.String `boil:"audio_file_url" json:"audio_file_url,omitempty" toml:"audio_file_url" yaml:"audio_file_url,omitempty"`

	R *tempSongInfoR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L tempSongInfoL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var TempSongInfoColumns = struct {
	TempSongInfoID string
	SongInfoID     string
	SongNumber     string
	SongName       string
	ArtistName     string
	ArtistType     string
	IsMR           string
	IsChosen22000  string
	Country        string
	Album          string
	OctaveCrawling string
	OctaveLibrosa  string
	VideoLink      string
	MelonSongID    string
	IsLive         string
	CreatedAt      string
	UpdatedAt      string
	DeletedAt      string
	HighTag        string
	LowTag         string
	AudioFileURL   string
}{
	TempSongInfoID: "temp_song_info_id",
	SongInfoID:     "song_info_id",
	SongNumber:     "song_number",
	SongName:       "song_name",
	ArtistName:     "artist_name",
	ArtistType:     "artist_type",
	IsMR:           "is_mr",
	IsChosen22000:  "is_chosen_22000",
	Country:        "country",
	Album:          "album",
	OctaveCrawling: "octave_crawling",
	OctaveLibrosa:  "octave_librosa",
	VideoLink:      "video_link",
	MelonSongID:    "melon_song_id",
	IsLive:         "is_live",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
	HighTag:        "high_tag",
	LowTag:         "low_tag",
	AudioFileURL:   "audio_file_url",
}

var TempSongInfoTableColumns = struct {
	TempSongInfoID string
	SongInfoID     string
	SongNumber     string
	SongName       string
	ArtistName     string
	ArtistType     string
	IsMR           string
	IsChosen22000  string
	Country        string
	Album          string
	OctaveCrawling string
	OctaveLibrosa  string
	VideoLink      string
	MelonSongID    string
	IsLive         string
	CreatedAt      string
	UpdatedAt      string
	DeletedAt      string
	HighTag        string
	LowTag         string
	AudioFileURL   string
}{
	TempSongInfoID: "temp_song_info.temp_song_info_id",
	SongInfoID:     "temp_song_info.song_info_id",
	SongNumber:     "temp_song_info.song_number",
	SongName:       "temp_song_info.song_name",
	ArtistName:     "temp_song_info.artist_name",
	ArtistType:     "temp_song_info.artist_type",
	IsMR:           "temp_song_info.is_mr",
	IsChosen22000:  "temp_song_info.is_chosen_22000",
	Country:        "temp_song_info.country",
	Album:          "temp_song_info.album",
	OctaveCrawling: "temp_song_info.octave_crawling",
	OctaveLibrosa:  "temp_song_info.octave_librosa",
	VideoLink:      "temp_song_info.video_link",
	MelonSongID:    "temp_song_info.melon_song_id",
	IsLive:         "temp_song_info.is_live",
	CreatedAt:      "temp_song_info.created_at",
	UpdatedAt:      "temp_song_info.updated_at",
	DeletedAt:      "temp_song_info.deleted_at",
	HighTag:        "temp_song_info.high_tag",
	LowTag:         "temp_song_info.low_tag",
	AudioFileURL:   "temp_song_info.audio_file_url",
}

// Generated where

var TempSongInfoWhere = struct {
	TempSongInfoID whereHelperint64
	SongInfoID     whereHelperint64
	SongNumber     whereHelperint
	SongName       whereHelpernull_String
	ArtistName     whereHelpernull_String
	ArtistType     whereHelpernull_String
	IsMR           whereHelpernull_Bool
	IsChosen22000  whereHelpernull_Bool
	Country        whereHelpernull_String
	Album          whereHelpernull_String
	OctaveCrawling whereHelpernull_String
	OctaveLibrosa  whereHelpernull_String
	VideoLink      whereHelpernull_String
	MelonSongID    whereHelpernull_String
	IsLive         whereHelpernull_Bool
	CreatedAt      whereHelpernull_Time
	UpdatedAt      whereHelpernull_Time
	DeletedAt      whereHelpernull_Time
	HighTag        whereHelpernull_Bool
	LowTag         whereHelpernull_Bool
	AudioFileURL   whereHelpernull_String
}{
	TempSongInfoID: whereHelperint64{field: "`temp_song_info`.`temp_song_info_id`"},
	SongInfoID:     whereHelperint64{field: "`temp_song_info`.`song_info_id`"},
	SongNumber:     whereHelperint{field: "`temp_song_info`.`song_number`"},
	SongName:       whereHelpernull_String{field: "`temp_song_info`.`song_name`"},
	ArtistName:     whereHelpernull_String{field: "`temp_song_info`.`artist_name`"},
	ArtistType:     whereHelpernull_String{field: "`temp_song_info`.`artist_type`"},
	IsMR:           whereHelpernull_Bool{field: "`temp_song_info`.`is_mr`"},
	IsChosen22000:  whereHelpernull_Bool{field: "`temp_song_info`.`is_chosen_22000`"},
	Country:        whereHelpernull_String{field: "`temp_song_info`.`country`"},
	Album:          whereHelpernull_String{field: "`temp_song_info`.`album`"},
	OctaveCrawling: whereHelpernull_String{field: "`temp_song_info`.`octave_crawling`"},
	OctaveLibrosa:  whereHelpernull_String{field: "`temp_song_info`.`octave_librosa`"},
	VideoLink:      whereHelpernull_String{field: "`temp_song_info`.`video_link`"},
	MelonSongID:    whereHelpernull_String{field: "`temp_song_info`.`melon_song_id`"},
	IsLive:         whereHelpernull_Bool{field: "`temp_song_info`.`is_live`"},
	CreatedAt:      whereHelpernull_Time{field: "`temp_song_info`.`created_at`"},
	UpdatedAt:      whereHelpernull_Time{field: "`temp_song_info`.`updated_at`"},
	DeletedAt:      whereHelpernull_Time{field: "`temp_song_info`.`deleted_at`"},
	HighTag:        whereHelpernull_Bool{field: "`temp_song_info`.`high_tag`"},
	LowTag:         whereHelpernull_Bool{field: "`temp_song_info`.`low_tag`"},
	AudioFileURL:   whereHelpernull_String{field: "`temp_song_info`.`audio_file_url`"},
}

// TempSongInfoRels is where relationship names are stored.
var TempSongInfoRels = struct {
}{}

// tempSongInfoR is where relationships are stored.
type tempSongInfoR struct {
}

// NewStruct creates a new relationship struct
func (*tempSongInfoR) NewStruct() *tempSongInfoR {
	return &tempSongInfoR{}
}

// tempSongInfoL is where Load methods for each relationship are stored.
type tempSongInfoL struct{}

var (
	tempSongInfoAllColumns            = []string{"temp_song_info_id", "song_info_id", "song_number", "song_name", "artist_name", "artist_type", "is_mr", "is_chosen_22000", "country", "album", "octave_crawling", "octave_librosa", "video_link", "melon_song_id", "is_live", "created_at", "updated_at", "deleted_at", "high_tag", "low_tag", "audio_file_url"}
	tempSongInfoColumnsWithoutDefault = []string{"song_info_id", "song_number", "song_name", "artist_name", "artist_type", "country", "album", "octave_crawling", "octave_librosa", "video_link", "melon_song_id", "deleted_at", "audio_file_url"}
	tempSongInfoColumnsWithDefault    = []string{"temp_song_info_id", "is_mr", "is_chosen_22000", "is_live", "created_at", "updated_at", "high_tag", "low_tag"}
	tempSongInfoPrimaryKeyColumns     = []string{"temp_song_info_id"}
	tempSongInfoGeneratedColumns      = []string{}
)

type (
	// TempSongInfoSlice is an alias for a slice of pointers to TempSongInfo.
	// This should almost always be used instead of []TempSongInfo.
	TempSongInfoSlice []*TempSongInfo
	// TempSongInfoHook is the signature for custom TempSongInfo hook methods
	TempSongInfoHook func(context.Context, boil.ContextExecutor, *TempSongInfo) error

	tempSongInfoQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	tempSongInfoType                 = reflect.TypeOf(&TempSongInfo{})
	tempSongInfoMapping              = queries.MakeStructMapping(tempSongInfoType)
	tempSongInfoPrimaryKeyMapping, _ = queries.BindMapping(tempSongInfoType, tempSongInfoMapping, tempSongInfoPrimaryKeyColumns)
	tempSongInfoInsertCacheMut       sync.RWMutex
	tempSongInfoInsertCache          = make(map[string]insertCache)
	tempSongInfoUpdateCacheMut       sync.RWMutex
	tempSongInfoUpdateCache          = make(map[string]updateCache)
	tempSongInfoUpsertCacheMut       sync.RWMutex
	tempSongInfoUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var tempSongInfoAfterSelectHooks []TempSongInfoHook

var tempSongInfoBeforeInsertHooks []TempSongInfoHook
var tempSongInfoAfterInsertHooks []TempSongInfoHook

var tempSongInfoBeforeUpdateHooks []TempSongInfoHook
var tempSongInfoAfterUpdateHooks []TempSongInfoHook

var tempSongInfoBeforeDeleteHooks []TempSongInfoHook
var tempSongInfoAfterDeleteHooks []TempSongInfoHook

var tempSongInfoBeforeUpsertHooks []TempSongInfoHook
var tempSongInfoAfterUpsertHooks []TempSongInfoHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *TempSongInfo) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range tempSongInfoAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *TempSongInfo) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range tempSongInfoBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *TempSongInfo) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range tempSongInfoAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *TempSongInfo) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range tempSongInfoBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *TempSongInfo) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range tempSongInfoAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *TempSongInfo) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range tempSongInfoBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *TempSongInfo) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range tempSongInfoAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *TempSongInfo) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range tempSongInfoBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *TempSongInfo) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range tempSongInfoAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddTempSongInfoHook registers your hook function for all future operations.
func AddTempSongInfoHook(hookPoint boil.HookPoint, tempSongInfoHook TempSongInfoHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		tempSongInfoAfterSelectHooks = append(tempSongInfoAfterSelectHooks, tempSongInfoHook)
	case boil.BeforeInsertHook:
		tempSongInfoBeforeInsertHooks = append(tempSongInfoBeforeInsertHooks, tempSongInfoHook)
	case boil.AfterInsertHook:
		tempSongInfoAfterInsertHooks = append(tempSongInfoAfterInsertHooks, tempSongInfoHook)
	case boil.BeforeUpdateHook:
		tempSongInfoBeforeUpdateHooks = append(tempSongInfoBeforeUpdateHooks, tempSongInfoHook)
	case boil.AfterUpdateHook:
		tempSongInfoAfterUpdateHooks = append(tempSongInfoAfterUpdateHooks, tempSongInfoHook)
	case boil.BeforeDeleteHook:
		tempSongInfoBeforeDeleteHooks = append(tempSongInfoBeforeDeleteHooks, tempSongInfoHook)
	case boil.AfterDeleteHook:
		tempSongInfoAfterDeleteHooks = append(tempSongInfoAfterDeleteHooks, tempSongInfoHook)
	case boil.BeforeUpsertHook:
		tempSongInfoBeforeUpsertHooks = append(tempSongInfoBeforeUpsertHooks, tempSongInfoHook)
	case boil.AfterUpsertHook:
		tempSongInfoAfterUpsertHooks = append(tempSongInfoAfterUpsertHooks, tempSongInfoHook)
	}
}

// One returns a single tempSongInfo record from the query.
func (q tempSongInfoQuery) One(ctx context.Context, exec boil.ContextExecutor) (*TempSongInfo, error) {
	o := &TempSongInfo{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysql: failed to execute a one query for temp_song_info")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all TempSongInfo records from the query.
func (q tempSongInfoQuery) All(ctx context.Context, exec boil.ContextExecutor) (TempSongInfoSlice, error) {
	var o []*TempSongInfo

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "mysql: failed to assign all query results to TempSongInfo slice")
	}

	if len(tempSongInfoAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all TempSongInfo records in the query.
func (q tempSongInfoQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to count temp_song_info rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q tempSongInfoQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "mysql: failed to check if temp_song_info exists")
	}

	return count > 0, nil
}

// TempSongInfos retrieves all the records using an executor.
func TempSongInfos(mods ...qm.QueryMod) tempSongInfoQuery {
	mods = append(mods, qm.From("`temp_song_info`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`temp_song_info`.*"})
	}

	return tempSongInfoQuery{q}
}

// FindTempSongInfo retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindTempSongInfo(ctx context.Context, exec boil.ContextExecutor, tempSongInfoID int64, selectCols ...string) (*TempSongInfo, error) {
	tempSongInfoObj := &TempSongInfo{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `temp_song_info` where `temp_song_info_id`=?", sel,
	)

	q := queries.Raw(query, tempSongInfoID)

	err := q.Bind(ctx, exec, tempSongInfoObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mysql: unable to select from temp_song_info")
	}

	if err = tempSongInfoObj.doAfterSelectHooks(ctx, exec); err != nil {
		return tempSongInfoObj, err
	}

	return tempSongInfoObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *TempSongInfo) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("mysql: no temp_song_info provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(tempSongInfoColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	tempSongInfoInsertCacheMut.RLock()
	cache, cached := tempSongInfoInsertCache[key]
	tempSongInfoInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			tempSongInfoAllColumns,
			tempSongInfoColumnsWithDefault,
			tempSongInfoColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(tempSongInfoType, tempSongInfoMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(tempSongInfoType, tempSongInfoMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `temp_song_info` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `temp_song_info` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `temp_song_info` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, tempSongInfoPrimaryKeyColumns))
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
		return errors.Wrap(err, "mysql: unable to insert into temp_song_info")
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

	o.TempSongInfoID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == tempSongInfoMapping["temp_song_info_id"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.TempSongInfoID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to populate default values for temp_song_info")
	}

CacheNoHooks:
	if !cached {
		tempSongInfoInsertCacheMut.Lock()
		tempSongInfoInsertCache[key] = cache
		tempSongInfoInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the TempSongInfo.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *TempSongInfo) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	tempSongInfoUpdateCacheMut.RLock()
	cache, cached := tempSongInfoUpdateCache[key]
	tempSongInfoUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			tempSongInfoAllColumns,
			tempSongInfoPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("mysql: unable to update temp_song_info, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `temp_song_info` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, tempSongInfoPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(tempSongInfoType, tempSongInfoMapping, append(wl, tempSongInfoPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "mysql: unable to update temp_song_info row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by update for temp_song_info")
	}

	if !cached {
		tempSongInfoUpdateCacheMut.Lock()
		tempSongInfoUpdateCache[key] = cache
		tempSongInfoUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q tempSongInfoQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update all for temp_song_info")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to retrieve rows affected for temp_song_info")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o TempSongInfoSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), tempSongInfoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `temp_song_info` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, tempSongInfoPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to update all in tempSongInfo slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to retrieve rows affected all in update all tempSongInfo")
	}
	return rowsAff, nil
}

var mySQLTempSongInfoUniqueColumns = []string{
	"temp_song_info_id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *TempSongInfo) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("mysql: no temp_song_info provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(tempSongInfoColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLTempSongInfoUniqueColumns, o)

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

	tempSongInfoUpsertCacheMut.RLock()
	cache, cached := tempSongInfoUpsertCache[key]
	tempSongInfoUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			tempSongInfoAllColumns,
			tempSongInfoColumnsWithDefault,
			tempSongInfoColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			tempSongInfoAllColumns,
			tempSongInfoPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("mysql: unable to upsert temp_song_info, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`temp_song_info`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `temp_song_info` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(tempSongInfoType, tempSongInfoMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(tempSongInfoType, tempSongInfoMapping, ret)
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
		return errors.Wrap(err, "mysql: unable to upsert for temp_song_info")
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

	o.TempSongInfoID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == tempSongInfoMapping["temp_song_info_id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(tempSongInfoType, tempSongInfoMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to retrieve unique values for temp_song_info")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to populate default values for temp_song_info")
	}

CacheNoHooks:
	if !cached {
		tempSongInfoUpsertCacheMut.Lock()
		tempSongInfoUpsertCache[key] = cache
		tempSongInfoUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single TempSongInfo record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *TempSongInfo) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("mysql: no TempSongInfo provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), tempSongInfoPrimaryKeyMapping)
	sql := "DELETE FROM `temp_song_info` WHERE `temp_song_info_id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete from temp_song_info")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by delete for temp_song_info")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q tempSongInfoQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("mysql: no tempSongInfoQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete all from temp_song_info")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by deleteall for temp_song_info")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o TempSongInfoSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(tempSongInfoBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), tempSongInfoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `temp_song_info` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, tempSongInfoPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "mysql: unable to delete all from tempSongInfo slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "mysql: failed to get rows affected by deleteall for temp_song_info")
	}

	if len(tempSongInfoAfterDeleteHooks) != 0 {
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
func (o *TempSongInfo) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindTempSongInfo(ctx, exec, o.TempSongInfoID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *TempSongInfoSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := TempSongInfoSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), tempSongInfoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `temp_song_info`.* FROM `temp_song_info` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, tempSongInfoPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "mysql: unable to reload all in TempSongInfoSlice")
	}

	*o = slice

	return nil
}

// TempSongInfoExists checks if the TempSongInfo row exists.
func TempSongInfoExists(ctx context.Context, exec boil.ContextExecutor, tempSongInfoID int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `temp_song_info` where `temp_song_info_id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, tempSongInfoID)
	}
	row := exec.QueryRowContext(ctx, sql, tempSongInfoID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "mysql: unable to check if temp_song_info exists")
	}

	return exists, nil
}

// Exists checks if the TempSongInfo row exists.
func (o *TempSongInfo) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return TempSongInfoExists(ctx, exec, o.TempSongInfoID)
}