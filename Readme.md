## swagger 문서화 후 실행하기
[참고문헌](https://github.com/swaggo/swag?tab=readme-ov-file#general-api-info)
1. handler 함수 위에 주석을 답니다. 
```
// GetRecommendation godoc
// @Summary      노래 추천 by 노래 번호 목록
// @Description  노래 번호 목록을 보내면 유사한 노래들을 추천합니다.
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Param        songs   body      RecommendRequest  true  "노래 번호 목록"
// @Success      200 {object} RecommendResponse "성공"
// @Router       /recommend [post]
```
2. 루트 경로에서 다음 명령어를 실행하면 docs 폴더가 갱신됩니다.
```
swag init -g ./cmd/main.go -o ./docs  
```
3. 서버를 실행합니다.
```
go run cmd/main.go
```