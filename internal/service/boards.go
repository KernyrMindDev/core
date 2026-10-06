package service

import (
	"errors"

	"github.com/KernyrMindDev/core/internal/apperr"
	"github.com/KernyrMindDev/core/internal/model"

	"gorm.io/gorm"
)

// BoardService 封装白板资源的业务逻辑
type BoardService struct {
	DB *gorm.DB
}

func NewBoardService(db *gorm.DB) *BoardService {
	return &BoardService{DB: db}
}

// getOwnedBoard 查询白板并校验调用者是否为所有者
func (s *BoardService) getOwnedBoard(uid, boardID string) (*model.Board, error) {
	var board model.Board
	if err := s.DB.First(&board, "id = ?", boardID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("白板不存在")
		}
		return nil, apperr.Internal("查询白板失败", err)
	}
	if board.UserID != uid {
		return nil, apperr.Forbidden("无权限访问此白板")
	}
	return &board, nil
}

// ListBoards 获取当前用户名下的所有白板
func (s *BoardService) ListBoards(uid string) ([]model.Board, error) {
	var boards []model.Board
	if err := s.DB.Where("user_id = ?", uid).Find(&boards).Error; err != nil {
		return nil, apperr.Internal("查询白板列表失败", err)
	}
	return boards, nil
}

// CreateBoard 创建新白板
func (s *BoardService) CreateBoard(uid, title string) (*model.Board, error) {
	board := model.Board{
		Title:  title,
		UserID: uid,
	}
	if err := s.DB.Create(&board).Error; err != nil {
		return nil, apperr.Internal("创建白板失败", err)
	}
	return &board, nil
}

// BoardDetail 白板详情聚合结果
type BoardDetail struct {
	Board       model.Board
	Nodes       []model.BoardNode
	Connections []model.BoardConnection
}

// GetBoardDetail 获取白板详情(含所有节点与连线)
func (s *BoardService) GetBoardDetail(uid, boardID string) (*BoardDetail, error) {
	board, err := s.getOwnedBoard(uid, boardID)
	if err != nil {
		return nil, err
	}

	var nodes []model.BoardNode
	var connections []model.BoardConnection
	if err := s.DB.Where("board_id = ?", boardID).Find(&nodes).Error; err != nil {
		return nil, apperr.Internal("查询节点失败", err)
	}
	if err := s.DB.Where("board_id = ?", boardID).Find(&connections).Error; err != nil {
		return nil, apperr.Internal("查询连线失败", err)
	}

	return &BoardDetail{
		Board:       *board,
		Nodes:       nodes,
		Connections: connections,
	}, nil
}

// UpdateBoard 更新白板标题
func (s *BoardService) UpdateBoard(uid, boardID, title string) (*model.Board, error) {
	board, err := s.getOwnedBoard(uid, boardID)
	if err != nil {
		return nil, err
	}

	board.Title = title
	if err := s.DB.Save(board).Error; err != nil {
		return nil, apperr.Internal("更新白板失败", err)
	}
	return board, nil
}

// DeleteBoard 删除白板及其关联的节点、连线
func (s *BoardService) DeleteBoard(uid, boardID string) error {
	board, err := s.getOwnedBoard(uid, boardID)
	if err != nil {
		return err
	}

	txErr := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("board_id = ?", boardID).Delete(&model.BoardNode{}).Error; err != nil {
			return err
		}
		if err := tx.Where("board_id = ?", boardID).Delete(&model.BoardConnection{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(board).Error; err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return apperr.Internal("删除白板及关联数据失败", txErr)
	}
	return nil
}
