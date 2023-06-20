package logic

import (
	"chesss/internal/db"
	"chesss/pkg/ent"
	"chesss/pkg/ent/chess"
	errors "errors"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func Move(ctx *gin.Context) {
	before := ctx.Param("before")
	client := db.InitializeDB()
	chess, err := client.Chess.Query().
		Where(chess.BeforeEQ(before)).
		Order(ent.Desc(chess.FieldCount)).
		First(ctx)
	if ent.IsNotFound(err) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "暂未记录此棋面",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "数据库异常",
		})
		return
	}
	check(chess)
	ctx.JSON(http.StatusOK, gin.H{
		"after": chess.After,
	})
}

func Memorize(ctx *gin.Context) {
	client := db.InitializeDB()

	before := ctx.PostForm("before")
	after := ctx.PostForm("after")

	if len(before) != 64 || len(after) != 64 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "传入参数不符合规则",
		})
		return
	}

	_, err := client.Chess.
		Create().
		SetBefore(before).
		SetAfter(after).
		SetCount(1).
		Save(ctx)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "数据库异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "保存成功",
	})
}

func IsValid(ctx *gin.Context) {
	chess := &ent.Chess{
		Before: ctx.Param("before"),
		After:  ctx.Param("after"),
	}
	res, _ := check(chess)
	if res {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "合法移动",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "不合法移动",
		})
	}
}

// 校验是否符合规则
func check(chess *ent.Chess) (bool, error) {
	var err error
	if len(chess.Before) != 64 || len(chess.After) != 64 {
		err = errors.New("暂未记录此棋面")
		return false, err
	}

	count, positions := getDifferentPosition(chess)

	if count > 2 {
		err = errors.New("暂未记录此棋面")
		return false, err
	}
	var move int
	if chess.After[positions[0]:positions[0]+1] != "9" {
		move = positions[0]
	} else {
		move = positions[1]
	}

	isMoveValid(chess, move)

	return true, nil
}

func isMoveValid(chess *ent.Chess, positions int) bool {
	var valid bool
	switch positions {
	case 0, 16, 32, 48:
		valid = isCarMoveValid(chess, positions)
	case 2, 14, 34, 46:
		valid = isHouseMoveValid(chess, positions)
	case 4, 12, 36, 44:
		valid = isEMoveValid(chess, positions)
	case 6, 10, 38, 42:
		valid = isSamuraiMoveValid(chess, positions)
	case 8, 40:
		valid = isWillMoveValid(chess, positions)
	case 18, 20, 50, 52:
		valid = isGunMoveValid(chess, positions)
	case 22, 24, 26, 28, 30, 54, 56, 58, 60, 62:
		valid = isSoldierMoveValid(chess, positions)
	}
	return valid
}

func getDifferentPosition(chess *ent.Chess) (int, []int) {
	var count int
	var positions []int

	for i := 0; i < len(chess.Before); i += 2 {
		if chess.After[i:(i+1)] != chess.After[i:i+1] {
			count++
			positions = append(positions, i)
		}
	}

	return count, positions
}

// 车
func isCarMoveValid(chess *ent.Chess, positions int) bool {
	fromX, _ := strconv.Atoi(chess.Before[positions : positions+1])
	fromY, _ := strconv.Atoi(chess.Before[positions+1 : positions+2])
	toX, _ := strconv.Atoi(chess.After[positions : positions+1])
	toY, _ := strconv.Atoi(chess.After[positions+1 : positions+2])
	if fromX == toX || fromY == toY {
		return true
	}
	return false
}

// 马
func isHouseMoveValid(chess *ent.Chess, positions int) bool {
	fromX, _ := strconv.Atoi(chess.Before[positions : positions+1])
	fromY, _ := strconv.Atoi(chess.Before[positions+1 : positions+2])
	toX, _ := strconv.Atoi(chess.After[positions : positions+1])
	toY, _ := strconv.Atoi(chess.After[positions+1 : positions+2])
	disY := toY - fromY
	disX := toX - fromX
	if math.Abs(float64(disY))*math.Abs(float64(disX)) != 2 {
		return false
	}
	if disY == 2 {
		location := strconv.Itoa((fromY+1)*10 + fromY)
		return isNotStumbling(chess, location)
	} else if disY == -2 {
		location := strconv.Itoa((fromY-1)*10 + fromY)
		return isNotStumbling(chess, location)
	} else if disX == 2 {
		location := strconv.Itoa((fromX+1)*10 + fromY)
		return isNotStumbling(chess, location)
	} else if disX == -2 {
		location := strconv.Itoa((fromX-1)*10 + fromY)
		return isNotStumbling(chess, location)
	}
	return true
}

// 是否绊马腿
func isNotStumbling(chess *ent.Chess, location string) bool {
	if strings.Contains(chess.Before, location) {
		return false
	}
	return true
}

// 象
func isEMoveValid(chess *ent.Chess, positions int) bool {
	fromX, _ := strconv.Atoi(chess.Before[positions : positions+1])
	fromY, _ := strconv.Atoi(chess.Before[positions+1 : positions+2])
	toX, _ := strconv.Atoi(chess.After[positions : positions+1])
	toY, _ := strconv.Atoi(chess.After[positions+1 : positions+2])
	disY := toY - fromY
	disX := toX - fromX
	if math.Abs(float64(disY)) == 2 && math.Abs(float64(disX)) == 2 {
		return false
	}
	// 是否过河
	if positions < 32 {
		if toY < 5 {
			return false
		}
	} else {
		if toY > 4 {
			return false
		}
	}
	// 是否塞象眼
	location := strconv.Itoa((fromX+toX)/2*10 + (fromY+toY)/2*10)
	return isNotStumbling(chess, location)
	return true
}

// 是否塞象眼
func isNotStuffedElephantEye(chess *ent.Chess, location string) bool {
	if strings.Contains(chess.Before, location) {
		return false
	}
	return true
}

// 士
func isSamuraiMoveValid(chess *ent.Chess, positions int) bool {
	fromX, _ := strconv.Atoi(chess.Before[positions : positions+1])
	fromY, _ := strconv.Atoi(chess.Before[positions+1 : positions+2])
	toX, _ := strconv.Atoi(chess.After[positions : positions+1])
	toY, _ := strconv.Atoi(chess.After[positions+1 : positions+2])
	disY := toY - fromY
	disX := toX - fromX
	if math.Abs(float64(disY)) == 1 && math.Abs(float64(disX)) == 1 {
		return false
	}
	if positions < 32 {
		if toX < 3 || toX > 5 || toY < 7 || toY > 9 {
			return false
		}
	} else {
		if toX < 3 || toX > 5 || toY < 0 || toY > 2 {
			return false
		}
	}
	return true
}

// 将
func isWillMoveValid(chess *ent.Chess, positions int) bool {
	fromX, _ := strconv.Atoi(chess.Before[positions : positions+1])
	fromY, _ := strconv.Atoi(chess.Before[positions+1 : positions+2])
	toX, _ := strconv.Atoi(chess.After[positions : positions+1])
	toY, _ := strconv.Atoi(chess.After[positions+1 : positions+2])
	disY := toY - fromY
	disX := toX - fromX
	if disY == 0 && math.Abs(float64(disX)) == 1 {
		return false
	}
	if disX == 0 && math.Abs(float64(disY)) == 1 {
		return false
	}
	if positions < 32 {
		if toX < 3 || toX > 5 || toY < 7 || toY > 9 {
			return false
		}
	} else {
		if toX < 3 || toX > 5 || toY < 0 || toY > 2 {
			return false
		}
	}
	return true
}

// 炮
func isGunMoveValid(chess *ent.Chess, positions int) bool {
	fromX, _ := strconv.Atoi(chess.Before[positions : positions+1])
	fromY, _ := strconv.Atoi(chess.Before[positions+1 : positions+2])
	toX, _ := strconv.Atoi(chess.After[positions : positions+1])
	toY, _ := strconv.Atoi(chess.After[positions+1 : positions+2])
	if fromX == toX || fromY == toY {
		return true
	}
	return false
}

// 兵
func isSoldierMoveValid(chess *ent.Chess, positions int) bool {
	fromX, _ := strconv.Atoi(chess.Before[positions : positions+1])
	fromY, _ := strconv.Atoi(chess.Before[positions+1 : positions+2])
	toX, _ := strconv.Atoi(chess.After[positions : positions+1])
	toY, _ := strconv.Atoi(chess.After[positions+1 : positions+2])
	disY := toY - fromY
	disX := toX - fromX
	if positions < 32 {
		// 红方没过河
		if fromY > 4 {
			if disX == 0 && disY == -1 {
				return true
			}
		} else {
			if disY == -1 && math.Abs(float64(disX)) == 0 {
				return true
			} else if math.Abs(float64(disX)) == 1 && disY == 0 {
				return true
			}
		}
	} else {
		// 黑方没过河
		if fromY < 5 {
			if disX == 0 && disY == -1 {
				return true
			}
		} else {
			if disY == 1 && math.Abs(float64(disX)) == 0 {
				return true
			} else if math.Abs(float64(disX)) == 1 && disY == 0 {
				return true
			}
		}
	}
	return false
}
