package handler

import (
	"SingSong-Server/internal/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestAddBlacklist(t *testing.T) {
	db, mock, err := pkg.NewTestDBMock()
	require.NoError(t, err)
	defer db.Close()

	handler := AddBlacklist(db)

	t.Run("successfully adds to blacklist", func(t *testing.T) {
		c, w := pkg.NewTestGinContext("POST", "/v1/blacklist", `{"memberId": 12345}`)
		c.Set("memberId", int64(67890))
		// Blocker ID와 Blocked ID를 사용하여 차단 여부 확인
		pkg.ExpectCountQuery(mock,
			"blacklist",
			"(blocker_member_id = ? AND blocked_member_id = ?)",
			[]any{int64(67890), int64(12345)},
			0,
		)
		// 차단 추가
		pkg.ExpectInsert(mock, "blacklist", []any{int64(67890), int64(12345), nil})
		// 차단된 ID 조회
		pkg.ExpectSelectByID(mock, "blacklist", "blacklist_id", 1, []string{"blacklist_id", "created_at", "updated_at"})
		handler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"success"`)
	})
}
