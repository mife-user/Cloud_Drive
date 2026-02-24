package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/cache"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"fmt"
)

// GetFavorites 获取用户收藏的文件列表
func (r *fileRepo) GetFavorites(ctx context.Context, userID uint) ([]domain.File, error) {
	var err error
	var favorites []domain.FileFavorite
	// 从缓存中获取收藏列表
	favorites, err = r.getFavoriteRecord(ctx, userID)
	if err != nil {
		logger.Error("查询收藏列表失败", logger.C(err))
		return nil, err
	}
	// 从缓存中获取文件信息
	files, err := r.getFavoriteFiles(ctx, favorites, userID)
	if err != nil {
		logger.Error("查询收藏文件失败", logger.C(err))
		return nil, err
	}
	return files, nil
}

// getFavoriteRecord 从缓存中获取收藏列表
func (r *fileRepo) getFavoriteRecord(ctx context.Context, userID uint) ([]domain.FileFavorite, error) {
	var err error
	var favorites []domain.FileFavorite
	var lovesJSONs map[string]string
	// 从缓存中获取收藏列表
	loverKey := fmt.Sprintf("lover:%d", userID)
	mapcmd := r.rd.HGetAll(ctx, loverKey)
	if mapcmd.Err() != nil {
		// 缓存中不存在，查询数据库
		if err = r.db.Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
			logger.Error("查询收藏列表失败", logger.C(err))
			return nil, err
		}
		for _, favorite := range favorites {
			fileIDSTR, err := exc.UintToStr(favorite.FileID)
			fileKey := fmt.Sprintf("file:%s", fileIDSTR)
			if err != nil {
				logger.Error("转换文件ID失败", logger.C(err))
				continue
			}
			favoriteJSON, err := exc.ExcFileToJSON(favorite)
			if err != nil {
				logger.Error("序列化收藏记录失败", logger.C(err))
				continue
			}
			if err = r.rd.HSet(ctx, loverKey, fileKey, favoriteJSON).Err(); err != nil {
				logger.Error("缓存收藏记录失败", logger.C(err))
				continue
			}
		}
	} else {
		lovesJSONs, err = mapcmd.Result()
		if err != nil {
			logger.Error("查询缓存收藏列表失败", logger.C(err))
			return nil, err
		}
		for _, value := range lovesJSONs {
			var favorite domain.FileFavorite
			if err = exc.ExcJSONToFile(value, &favorite); err != nil {
				logger.Error("反序列化收藏记录失败", logger.C(err))
				continue
			}
			favorites = append(favorites, favorite)
		}
	}
	return favorites, nil
}

// getFavoriteFiles 从缓存中获取收藏文件
func (r *fileRepo) getFavoriteFiles(ctx context.Context, favorites []domain.FileFavorite, userID uint) ([]domain.File, error) {
	userKey := fmt.Sprintf("files:%d", userID)
	var files []domain.File
	for _, favorite := range favorites {
		var file domain.File
		// 从缓存中获取文件信息
		fileIDSTR := fmt.Sprintf("file:%d", favorite.FileID)
		fileJSON, err := r.rd.HGet(ctx, userKey, fileIDSTR).Result()
		// 缓存中不存在，查询数据库
		if err != nil {
			if err = r.db.Where("id = ?", favorite.FileID).First(&file).Error; err != nil {
				logger.Debug("查询文件失败，跳过", logger.C(err))
				continue
			}
			if fileJSON, err = exc.ExcFileToJSON(file); err != nil {
				logger.Error("序列化文件信息失败", logger.C(err))
				continue
			}
			if err = r.rd.HSet(ctx, userKey, fileIDSTR, fileJSON).Err(); err != nil {
				logger.Error("缓存文件信息失败", logger.C(err))
				continue
			}
			if file.UserID != userID && file.Permissions != "public" {
				logger.Debug("用户不再有文件访问权限，跳过", logger.U("user_id", userID), logger.S("file_id", fmt.Sprintf("%d", file.ID)))
				continue
			}
			// 设置缓存过期时间，使用带随机偏移的缓存策略
			ttl := cache.FileCacheConfig.RandomTTL()
			if err = r.rd.Expire(ctx, userKey, ttl).Err(); err != nil {
				logger.Error("设置缓存过期时间失败", logger.C(err))
				continue
			}
		} else {
			if err = exc.ExcJSONToFile(fileJSON, &file); err != nil {
				logger.Error("反序列化文件信息失败", logger.C(err))
				continue
			}
			if file.DeletedAt.Valid {
				logger.Error("文件已被删除", logger.U("file_id", favorite.FileID))
				continue
			}
		}
		files = append(files, file)
	}
	return files, nil
}
