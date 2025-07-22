package helpers

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
)

type DataTableConfig struct {
	Model     any
	Query     orm.Query
	Search    []string
	FormatRow func(index int, model any) map[string]any
}

func RenderDataTable(ctx http.Context, config DataTableConfig) http.Response {
	// Ambil parameter datatables
	draw := ctx.Request().Query("draw", "1")
	start := ctx.Request().Query("start", "0")
	length := ctx.Request().Query("length", "10")
	searchValue := ctx.Request().Query("search[value]", "")

	startInt, _ := strconv.Atoi(start)
	lengthInt, _ := strconv.Atoi(length)
	drawInt, _ := strconv.Atoi(draw)

	// Total semua data
	queryTotal := config.Query.Model(config.Model)
	var totalRecords int64
	totalRecords, err := queryTotal.Count()
	if err != nil {
		fmt.Printf("Error counting users: %v\n", err)
	}

	// Query filter pencarian
	queryFiltered := config.Query.Model(config.Model)
	if searchValue != "" && len(config.Search) > 0 {
		for i, col := range config.Search {
			if i == 0 {
				queryFiltered = queryFiltered.Where(col+" LIKE ?", "%"+searchValue+"%")
			} else {
				queryFiltered = queryFiltered.OrWhere(col+" LIKE ?", "%"+searchValue+"%")
			}
		}
	}

	var filteredRecords int64
	filteredRecords, err = queryFiltered.Count()
	if err != nil {
		fmt.Printf("Error counting users: %v\n", err)
	}

	// Siapkan slice model untuk result
	modelType := reflect.TypeOf(config.Model).Elem()
	sliceType := reflect.SliceOf(reflect.PointerTo(modelType))
	slicePtr := reflect.New(sliceType).Interface()

	// Ambil data dengan limit + offset
	queryFiltered.Offset(startInt).Limit(lengthInt).Find(slicePtr)

	// Format setiap row
	sliceValue := reflect.ValueOf(slicePtr).Elem()
	var data []map[string]any
	for i := 0; i < sliceValue.Len(); i++ {
		model := sliceValue.Index(i).Interface()
		data = append(data, config.FormatRow(startInt+i, model))
	}

	// Kirim response JSON
	return ctx.Response().Json(http.StatusOK, map[string]any{
		"draw":            drawInt,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            data,
	})
}
