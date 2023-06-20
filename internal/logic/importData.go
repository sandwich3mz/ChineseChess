package logic

import (
	"chesss/internal/db"
	"chesss/pkg/ent"
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

func ImportData(ctx context.Context) {
	client := db.InitializeDB()
	defer client.Close()
	num, err := client.Chess.Query().Count(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	if num > 1 {
		return
	}
	file, err := os.Open("./resource/res.csv")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// bulk := make([]*ent.ChessCreate, 50)
	var bulk []*ent.ChessCreate
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			return
		}
		before := record[0][:64]
		after := record[0][64:]
		count := record[1]
		frequency, err := strconv.Atoi(count)
		if err != nil {
			log.Fatal(err)
			return
		}
		bulk = append(bulk, client.Chess.Create().
			SetBefore(before).
			SetAfter(after).
			SetCount(int64(frequency)))
		if len(bulk) > 49 {
			_, err := client.Chess.
				CreateBulk(bulk...).
				Save(ctx)
			if err != nil {
				log.Fatal(err)
				return
			}
			bulk = bulk[:0]
		}
	}
	// 如果 bulk 数组中还有剩余的实体对象，则保存它们
	if len(bulk) > 0 {
		if _, err := client.Chess.
			CreateBulk(bulk...).
			Save(ctx); err != nil {
			log.Fatal(err)
			return
		}
	}
}
