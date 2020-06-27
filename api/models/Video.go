package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Video  . . 
type Video struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"text;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
// Prepare . . .
func (p *Video) Prepare() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validate . . .
func (p *Video) Validate() map[string]string {

	var err error

	var errorMessages = make(map[string]string)

	if p.Title == "" {
		err = errors.New("Required Title")
		errorMessages["Required_title"] = err.Error()

	}
	if p.Content == "" {
		err = errors.New("Required Content")
		errorMessages["Required_content"] = err.Error()

	}
	if p.AuthorID < 1 {
		err = errors.New("Required Author")
		errorMessages["Required_author"] = err.Error()
	}
	return errorMessages
}
// SaveVideo . . 
func (p *Video) SaveVideo(db *gorm.DB) (*Video, error) {
	var err error
	err = db.Debug().Model(&Video{}).Create(&p).Error
	if err != nil {
		return &Video{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Video{}, err
		}
	}
	return p, nil
}

// FindAllVideos . . .
func (p *Video) FindAllVideos(db *gorm.DB) (*[]Video, error) {
	var err error
	videos := []Video{}
	err = db.Debug().Model(&Video{}).Limit(100).Order("created_at desc").Find(&videos).Error
	if err != nil {
		return &[]Video{}, err
	}
	if len(videos) > 0 {
		for i, _ := range videos {
			err := db.Debug().Model(&User{}).Where("id = ?", videos[i].AuthorID).Take(&videos[i].Author).Error
			if err != nil {
				return &[]Video{}, err
			}
		}
	}
	return &videos, nil
}

// FindVideoByID  . .. 
func (p *Video) FindVideoByID(db *gorm.DB, pid uint64) (*Video, error) {
	var err error
	err = db.Debug().Model(&Video{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Video{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Video{}, err
		}
	}
	return p, nil
}



// UpdateAVideo . . 
func (p *Video) UpdateAVideo(db *gorm.DB) (*Video, error) {

	var err error

	err = db.Debug().Model(&Video{}).Where("id = ?", p.ID).Updates(Video{Title: p.Title, Content: p.Content, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Video{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Video{}, err
		}
	}
	return p, nil
}

// DeleteAVideo . . 
func (p *Video) DeleteAVideo(db *gorm.DB) (int64, error) {

	db = db.Debug().Model(&Video{}).Where("id = ?", p.ID).Take(&Video{}).Delete(&Video{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}


// FindUserVideos . .
func (p *Video) FindUserVideos(db *gorm.DB, uid uint32) (*[]Video, error) {

	var err error
	videos := []Video{}
	err = db.Debug().Model(&Video{}).Where("author_id = ?", uid).Limit(100).Order("created_at desc").Find(&videos).Error
	if err != nil {
		return &[]Video{}, err
	}
	if len(videos) > 0 {
		for i, _ := range videos {
			err := db.Debug().Model(&User{}).Where("id = ?", videos[i].AuthorID).Take(&videos[i].Author).Error
			if err != nil {
				return &[]Video{}, err
			}
		}
	}
	return &videos, nil
}

//When a user is deleted, we also delete the post that the user had

// DeleteUserVideos . . .
func (p *Video) DeleteUserVideos(db *gorm.DB, uid uint32) (int64, error) {
	videos := []Video{}
	db = db.Debug().Model(&Video{}).Where("author_id = ?", uid).Find(&videos).Delete(&videos)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

