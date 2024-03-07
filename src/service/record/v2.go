package record

import (
	"errors"
	"qcg-center/src/dao"
	daoEntities "qcg-center/src/entities/dao"
	"qcg-center/src/entities/dto"
	"qcg-center/src/util"
	"strings"

	"github.com/gin-gonic/gin"

	"encoding/json"
	"reflect"
)

type RecordService struct {
	db dao.IDatabaseManager
}

func NewRecordService(db dao.IDatabaseManager) *RecordService {
	return &RecordService{
		db: db,
	}
}

func getBasicInfo(recordDTO interface{}) (dto.BasicInfo, error) {
	// 取出basic: BasicInfo
	basic := reflect.ValueOf(recordDTO).FieldByName("Basic")

	if !basic.IsValid() {
		return dto.BasicInfo{}, errors.New("basic not found")
	}

	return basic.Interface().(dto.BasicInfo), nil
}

func GetRemoteAddr(c *gin.Context) string {
	remoteAddr := c.Request.Header.Get("x-forwarded-for")

	if remoteAddr == "" {
		// 从remoteAddr取
		addr := c.Request.RemoteAddr
		// 分割IP和端口
		remoteAddrSlice := strings.Split(addr, ":")
		// 只取IP
		remoteAddr = remoteAddrSlice[0]
	}

	return remoteAddr
}

func (s *RecordService) InsertRecord(c *gin.Context, recordDTO interface{}) error {
	basic, err := getBasicInfo(recordDTO)

	if err != nil {
		return err
	}

	// 取出IP
	remoteAddr := GetRemoteAddr(c)

	// 插入 IP 地理信息
	ipGeoByte, err := util.GetIPGeoJSONBytes(remoteAddr)
	if err != nil {
		return err
	}

	// 反序列化转换成 dao.IPGeoInfoDAO
	ipGeoInfoDAO := daoEntities.IPGeoInfoDAO{
		CreatedAt: util.GetCSTTime(),
	}

	err = json.Unmarshal(ipGeoByte, &ipGeoInfoDAO)

	if err != nil {
		return err
	}

	if ipGeoInfoDAO.Status != "success" {
		return errors.New("ip geo info not found, status: " + ipGeoInfoDAO.Status + ", query: " + ipGeoInfoDAO.IP)
	}

	err = s.db.InsertIPGeoInfo(ipGeoInfoDAO)

	if err != nil {
		return err
	}

	// 插入 标识符元组
	identifierDAO := daoEntities.IdentifierTupleDAO{
		InstanceID: basic.InstanceID,
		HostID:     basic.HostID,
		IP:         remoteAddr,
		CreatedAt:  util.GetCSTTime(),
	}

	err = s.db.InsertIdentifierTuple(identifierDAO)
	if err != nil {
		return err
	}

	// 插入记录
	recordDAO := daoEntities.CommonRecordDAO{
		RemoteAddr: remoteAddr,
		Time:       util.GetCSTTime(),
		Data:       recordDTO,
	}

	return s.db.InsertRecord(recordDAO)
}
