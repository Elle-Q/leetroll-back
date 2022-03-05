package service

import (
	"gomo/admin/service/vo"
	"gomo/db/handlers"
	"gomo/db/models"
)

type ItemService struct {
	ItemHandler *handlers.ItemHandler
	FileHandler *handlers.FileHandler
	Error error
}


func (e *ItemService) GetFilesById(ID int, vo *vo.ItemFilesVO) *ItemService{
	itemHandler := e.ItemHandler
	fileHandler := e.FileHandler

	//获取item的相关信息(名称)
	item := models.MakeItem()
	err := itemHandler.Get(ID, item).Error
	if err != nil {
		e.Error = err
		return e
	}

	files := make([]models.File, 0)
	//获取文件信息
	fileHandler.QueryByItemId(ID, &files)

	main := make([]models.File, 0)
	prev := make([]models.File, 0)
	refs := make([]models.File, 0)

	for _, f := range files{
		switch f.Type {
		case "main":
			main = append(main, f)
			break
		case "preview":
			prev = append(prev, f)
			break
		case "refs":
			refs = append(refs, f)
			break
		}
	}
	vo.Main = main
	vo.Preview = prev
	vo.Refs = refs
	vo.ID = int64(ID)
	vo.ItemName = item.Name

	return e
}
