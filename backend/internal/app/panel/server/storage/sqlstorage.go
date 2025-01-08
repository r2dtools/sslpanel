package storage

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

const guidPrefix = "server_guid_prefix"

func (s *Server) AfterFind(tx *gorm.DB) (err error) {
	guidRaw := fmt.Sprintf("%s_%d", guidPrefix, s.ID)
	s.Guid = base64.RawStdEncoding.EncodeToString([]byte(guidRaw))

	return
}

type sqlStorage struct {
	db *gorm.DB
}

func (s sqlStorage) FindByID(id int) (*Server, error) {
	var server Server
	err := s.db.First(&server, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("could not find server with ID %d: %v", id, err)
	}

	return &server, nil
}

func (s sqlStorage) FindByGuid(guid string) (*Server, error) {
	decodedGuid, err := base64.RawStdEncoding.DecodeString(guid)

	if err != nil {
		return nil, fmt.Errorf("invalid server guid: %v", err)
	}

	idStr := strings.TrimPrefix(string(decodedGuid), fmt.Sprintf("%s_", guidPrefix))
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return nil, fmt.Errorf("invalid server guid: %v", err)
	}

	return s.FindByID(id)
}

func (s sqlStorage) FindAllByAccountID(accountID int) ([]Server, error) {
	var servers []Server
	err := s.db.Where("account_id = ?", accountID).Order("id desc").Find(&servers).Error

	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (s sqlStorage) Save(server *Server) error {
	if server.ID == 0 {
		return s.db.Create(server).Error
	}

	return s.db.Save(server).Error
}

func (s sqlStorage) Remove(server *Server) error {
	err := s.db.Delete(server).Error

	if err != nil {
		return fmt.Errorf("failed to delete server with ID %d: %v", server.ID, err)
	}

	return nil
}

func (s sqlStorage) FindCountByIP(ipv4, ipv6 string, excludeIds []int) (int, error) {
	var count int64

	if ipv4 == "" && ipv6 == "" {
		return 0, errors.New("ipv4 or ipv6 must be specified")
	}

	db := s.db

	if ipv4 != "" {
		db = db.Where("ipv4_address = ?", ipv4)
	}

	if ipv6 != "" {
		db = db.Or("ipv6_address = ?", ipv6)
	}

	db = s.db.Where(db)

	if len(excludeIds) != 0 {
		db = db.Where("id NOT IN ?", excludeIds)
	}

	err := db.Table("servers").Count(&count).Error

	return int(count), err
}

func NewServerSqlStorage(db *gorm.DB) ServerStorage {
	return sqlStorage{
		db: db,
	}
}

func (*Server) TableName() string {
	return "servers"
}
