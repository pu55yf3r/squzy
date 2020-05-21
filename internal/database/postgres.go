package database

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"time"
)

type postgres struct {
	db *gorm.DB
}

type Snapshot struct {
	gorm.Model
	SchedulerID string    `gorm:"column:schedulerId"`
	Code        string    `gorm:"column:code"`
	Type        string    `gorm:"column:type"`
	Error       string    `gorm:"column:error"`
	Meta        *MetaData `gorm:"column:meta"`
}

type MetaData struct {
	gorm.Model
	SnapshotID uint      `gorm:"column:snapshotId"`
	StartTime  time.Time `gorm:"column:startTime"`
	EndTime    time.Time `gorm:"column:endTime"`
	Value      []byte    `gorm:"column:value"` //TODO: google
}

//Agent gorm description
type StatRequest struct {
	gorm.Model
	AgentID    string `gorm:"column:agentID"`
	AgentName  string `gorm:"column:agentName"`
	CPUInfo    []*CPUInfo
	MemoryInfo *MemoryInfo `gorm:"column:memoryInfo"`
	DiskInfo   []*DiskInfo `gorm:"column:diskInfo"`
	NetInfo    []*NetInfo  `gorm:"column:netInfo"`
	Time       time.Time   `gorm:"column:time"`
}

const (
	cpuInfoKey  = "CPUInfo"
	diskInfoKey = "DiskInfo"
	netInfoKey  = "NetInfo"
)

type CPUInfo struct {
	gorm.Model
	StatRequestID uint    `gorm:"column:statRequestId"`
	Load          float64 `gorm:"column:load"`
}

type MemoryInfo struct {
	gorm.Model
	StatRequestID uint        `gorm:"column:statRequestId"`
	Mem           *MemoryMem  `gorm:"column:mem"`
	Swap          *MemorySwap `gorm:"column:swap"`
}

type MemoryMem struct {
	gorm.Model
	MemoryInfoID uint    `gorm:"column:memoryInfoId"`
	Total        uint64  `gorm:"column:total"`
	Used         uint64  `gorm:"column:used"`
	Free         uint64  `gorm:"column:free"`
	Shared       uint64  `gorm:"column:shared"`
	UsedPercent  float64 `gorm:"column:usedPercent"`
}

type MemorySwap struct {
	gorm.Model
	MemoryInfoID uint    `gorm:"column:memoryInfoId"`
	Total        uint64  `gorm:"column:total"`
	Used         uint64  `gorm:"column:used"`
	Free         uint64  `gorm:"column:free"`
	Shared       uint64  `gorm:"column:shared"`
	UsedPercent  float64 `gorm:"column:usedPercent"`
}

type DiskInfo struct {
	gorm.Model
	StatRequestID uint    `gorm:"column:statRequestId"`
	Name          string  `gorm:"column:name"`
	Total         uint64  `gorm:"column:total"`
	Free          uint64  `gorm:"column:free"`
	Used          uint64  `gorm:"column:used"`
	UsedPercent   float64 `gorm:"column:usedPercent"`
}

type NetInfo struct {
	gorm.Model
	StatRequestID uint   `gorm:"column:statRequestId"`
	Name          string `gorm:"column:name"`
	BytesSent     uint64 `gorm:"column:bytesSent"`
	BytesRecv     uint64 `gorm:"column:bytesRecv"`
	PacketsSent   uint64 `gorm:"column:packetsSent"`
	PacketsRecv   uint64 `gorm:"column:packetsRecv"`
	ErrIn         uint64 `gorm:"column:errIn"`
	ErrOut        uint64 `gorm:"column:errOut"`
	DropIn        uint64 `gorm:"column:dropIn"`
	DropOut       uint64 `gorm:"column:dropOut"`
}

const (
	dbSnapshotCollection    = "snapshots"     //TODO: check
	dbStatRequestCollection = "stat_requests" //TODO: check
)

var (
	errorDataBase = errors.New("ERROR_DATABASE_OPERATION")
)

func (p *postgres) Migrate() error {
	models := []interface{}{
		&Snapshot{},
		&MetaData{},
		&StatRequest{},
		&CPUInfo{},
		&MemoryInfo{},
		&MemoryMem{},
		&MemorySwap{},
		&DiskInfo{},
		&NetInfo{},
	}

	var err error
	for _, model := range models {
		err = p.db.AutoMigrate(model).Error // migrate models one-by-one
	}

	return err
}

func (p *postgres) InsertSnapshot(data *apiPb.SchedulerResponse) error {
	snapshot, err := ConvertToPostgresSnapshot(data)
	if err != nil {
		return err
	}
	if err := p.db.Table(dbSnapshotCollection).Create(snapshot).Error; err != nil {
		return errorDataBase
	}
	return nil
}

func (p *postgres) GetSnapshots(schedulerID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.SchedulerSnapshot, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int
	err = p.db.Table(dbSnapshotCollection).
		Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), schedulerID).
		Where(fmt.Sprintf(`"%s"."created_at" BETWEEN ? and ?`, dbSnapshotCollection), timeFrom, timeTo).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	//TODO: test if it works
	var dbSnapshots []*Snapshot
	err = p.db.
		Table(dbSnapshotCollection).
		Set("gorm:auto_preload", true).
		Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), schedulerID).
		Where(fmt.Sprintf(`"%s"."created_at" BETWEEN ? and ?`, dbSnapshotCollection), timeFrom, timeTo).
		Order("created_at").
		Offset(offset).
		Limit(limit).
		Find(&dbSnapshots).Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgresSnapshots(dbSnapshots), int32(count), nil
}

func (p *postgres) InsertStatRequest(data *apiPb.Metric) error {
	pgData, err := ConvertToPostgressStatRequest(data)
	if err != nil {
		return err
	}
	if err := p.db.Table(dbStatRequestCollection).Create(pgData).Error; err != nil {
		//TODO: log?
		return errorDataBase
	}
	return nil
}

func (p *postgres) GetStatRequest(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int
	err = p.db.Table(dbStatRequestCollection).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	//TODO: test if it works
	var statRequests []*StatRequest
	err = p.db.
		Set("gorm:auto_preload", true).
		//Preload("disk_infos").
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Order("time").
		Offset(offset).
		Limit(limit).
		Find(&statRequests).Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgressStatRequests(statRequests), int32(count), nil
}

func (p *postgres) GetCPUInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, cpuInfoKey)
}

func (p *postgres) GetMemoryInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int
	err = p.db.Table(dbStatRequestCollection).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	//TODO: test if it works
	var statRequests []*StatRequest
	err = p.db.
		Preload("MemoryInfo").
		Preload("MemoryInfo.Mem").
		Preload("MemoryInfo.Swap").
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Order("time").
		Offset(offset).
		Limit(limit).
		Find(&statRequests).
		Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgressStatRequests(statRequests), int32(count), nil
}

func (p *postgres) GetDiskInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, diskInfoKey)
}

func (p *postgres) GetNetInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, netInfoKey)
}

func (p *postgres) getSpecialRecords(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter, key string) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int
	err = p.db.Table(dbStatRequestCollection).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	//TODO: test if it works
	var statRequests []*StatRequest
	err = p.db.
		Preload(key).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Order("time").
		Offset(offset).
		Limit(limit).
		Find(&statRequests).
		Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgressStatRequests(statRequests), int32(count), nil
}

func getTime(filter *apiPb.TimeFilter) (time.Time, time.Time, error) {
	timeFrom := time.Unix(0, 0)
	timeTo := time.Now()
	var err error
	if filter != nil {
		if filter.GetFrom() != nil {
			timeFrom, err = ptypes.Timestamp(filter.From)
		}
		if filter.GetTo() != nil {
			timeTo, err = ptypes.Timestamp(filter.To)
		}
	}
	return timeFrom, timeTo, err
}

//Return offset and limit
func getOffsetAndLimit(count int, pagination *apiPb.Pagination) (int32, int32) {
	if pagination != nil {
		if pagination.Page == -1 {
			return int32(count) - pagination.Limit, pagination.Limit
		}
		return pagination.GetLimit() * pagination.GetPage(), pagination.GetLimit()
	}
	return int32(0), int32(count)
}