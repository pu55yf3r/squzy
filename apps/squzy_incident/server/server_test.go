package server

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"squzy/apps/squzy_incident/database"
	"squzy/apps/squzy_incident/expression"
	"testing"
)

type mockStorage struct {
}

func (m mockStorage) SaveResponseFromScheduler(ctx context.Context, in *apiPb.SchedulerResponse, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, nil
}

func (m mockStorage) SaveResponseFromAgent(ctx context.Context, in *apiPb.Metric, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, nil
}

func (m mockStorage) SaveTransaction(ctx context.Context, in *apiPb.TransactionInfo, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, nil
}

func (m mockStorage) GetSchedulerInformation(ctx context.Context, in *apiPb.GetSchedulerInformationRequest, opts ...grpc.CallOption) (*apiPb.GetSchedulerInformationResponse, error) {
	return nil, nil
}

func (m mockStorage) GetSchedulerUptime(ctx context.Context, in *apiPb.GetSchedulerUptimeRequest, opts ...grpc.CallOption) (*apiPb.GetSchedulerUptimeResponse, error) {
	return nil, nil
}

func (m mockStorage) GetAgentInformation(ctx context.Context, in *apiPb.GetAgentInformationRequest, opts ...grpc.CallOption) (*apiPb.GetAgentInformationResponse, error) {
	return nil, nil
}

func (m mockStorage) GetTransactionsGroup(ctx context.Context, in *apiPb.GetTransactionGroupRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionGroupResponse, error) {
	return nil, nil
}

func (m mockStorage) GetTransactions(ctx context.Context, in *apiPb.GetTransactionsRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionsResponse, error) {
	return nil, nil
}

func (m mockStorage) GetTransactionById(ctx context.Context, in *apiPb.GetTransactionByIdRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionByIdResponse, error) {
	return nil, nil
}

func (m mockStorage) SaveIncident(ctx context.Context, in *apiPb.Incident, opts ...grpc.CallOption) (*empty.Empty, error) {
	if in.RuleId == isIncidentExistwasIncident.Hex() {
		return nil, errors.New("ERROR")
	}
	return nil, nil
}

func (m mockStorage) UpdateIncidentStatus(ctx context.Context, in *apiPb.UpdateIncidentStatusRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	if in.IncidentId == incidentExistIncidentOpenedIncident.Hex() {
		return nil, errors.New("ERROR")
	}
	return nil, nil
}

func (m mockStorage) GetIncidentById(ctx context.Context, in *apiPb.IncidentIdRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	return nil, nil
}

func (m mockStorage) GetIncidentByRuleId(ctx context.Context, in *apiPb.RuleIdRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	if in.RuleId == incidentExistIncidentOpenedIncident.Hex() {
		return &apiPb.Incident{
			Id:     incidentExistIncidentOpenedIncident.Hex(),
			Status: apiPb.IncidentStatus_INCIDENT_STATUS_OPENED,
		}, nil
	}
	return nil, nil
}

func (m mockStorage) GetIncidentsList(ctx context.Context, in *apiPb.GetIncidentsListRequest, opts ...grpc.CallOption) (*apiPb.GetIncidentsListResponse, error) {
	return nil, nil
}

type mockDatabase struct {
}

func (m mockDatabase) SaveRule(context.Context, *database.Rule) error {
	return nil
}

func (m mockDatabase) FindRuleById(context.Context, primitive.ObjectID) (*database.Rule, error) {
	return &database.Rule{}, nil
}

func (m mockDatabase) FindRulesByOwnerId(ctx context.Context, ownerType apiPb.RuleOwnerType, ownerId primitive.ObjectID) ([]*database.Rule, error) {
	if ownerId == ruleIsNotActive {
		return []*database.Rule{
			{
				Status: apiPb.RuleStatus_RULE_STATUS_UNSPECIFIED,
				Rule:   "len(Last(10)) > 0",
			},
		}, nil
	}
	if ownerId == incidentExistIncidentOpenedIncident {
		return []*database.Rule{
			{
				Id:     incidentExistIncidentOpenedIncident,
				Status: apiPb.RuleStatus_RULE_STATUS_ACTIVE,
				Rule:   "len(Last(10)) > 0",
			},
		}, nil
	}
	if ownerId == isIncidentExistwasIncident {
		return []*database.Rule{
			{
				Id:     isIncidentExistwasIncident,
				Status: apiPb.RuleStatus_RULE_STATUS_ACTIVE,
				Rule:   "5 == 5",
			},
		}, nil
	}
	return []*database.Rule{
		{
			Status: apiPb.RuleStatus_RULE_STATUS_ACTIVE,
			Rule:   "len(Last(10)) > 0",
		},
	}, nil
}

func (m mockDatabase) RemoveRule(ctx context.Context, ruleId primitive.ObjectID) (*database.Rule, error) {
	return &database.Rule{}, nil
}

func (m mockDatabase) ActivateRule(ctx context.Context, ruleId primitive.ObjectID) (*database.Rule, error) {
	return &database.Rule{}, nil
}

func (m mockDatabase) DeactivateRule(ctx context.Context, ruleId primitive.ObjectID) (*database.Rule, error) {
	return &database.Rule{}, nil
}

type mockErrorStorage struct {
}

func (m mockErrorStorage) SaveResponseFromScheduler(ctx context.Context, in *apiPb.SchedulerResponse, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) SaveResponseFromAgent(ctx context.Context, in *apiPb.Metric, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) SaveTransaction(ctx context.Context, in *apiPb.TransactionInfo, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetSchedulerInformation(ctx context.Context, in *apiPb.GetSchedulerInformationRequest, opts ...grpc.CallOption) (*apiPb.GetSchedulerInformationResponse, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetSchedulerUptime(ctx context.Context, in *apiPb.GetSchedulerUptimeRequest, opts ...grpc.CallOption) (*apiPb.GetSchedulerUptimeResponse, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetAgentInformation(ctx context.Context, in *apiPb.GetAgentInformationRequest, opts ...grpc.CallOption) (*apiPb.GetAgentInformationResponse, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetTransactionsGroup(ctx context.Context, in *apiPb.GetTransactionGroupRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionGroupResponse, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetTransactions(ctx context.Context, in *apiPb.GetTransactionsRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionsResponse, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetTransactionById(ctx context.Context, in *apiPb.GetTransactionByIdRequest, opts ...grpc.CallOption) (*apiPb.GetTransactionByIdResponse, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) SaveIncident(ctx context.Context, in *apiPb.Incident, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) UpdateIncidentStatus(ctx context.Context, in *apiPb.UpdateIncidentStatusRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetIncidentById(ctx context.Context, in *apiPb.IncidentIdRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetIncidentByRuleId(ctx context.Context, in *apiPb.RuleIdRequest, opts ...grpc.CallOption) (*apiPb.Incident, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorStorage) GetIncidentsList(ctx context.Context, in *apiPb.GetIncidentsListRequest, opts ...grpc.CallOption) (*apiPb.GetIncidentsListResponse, error) {
	return nil, errors.New("ERROR")
}

type mockErrorDatabase struct {
}

func (m mockErrorDatabase) SaveRule(context.Context, *database.Rule) error {
	return errors.New("ERROR")
}

func (m mockErrorDatabase) FindRuleById(context.Context, primitive.ObjectID) (*database.Rule, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorDatabase) FindRulesByOwnerId(ctx context.Context, ownerType apiPb.RuleOwnerType, ownerId primitive.ObjectID) ([]*database.Rule, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorDatabase) RemoveRule(ctx context.Context, ruleId primitive.ObjectID) (*database.Rule, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorDatabase) ActivateRule(ctx context.Context, ruleId primitive.ObjectID) (*database.Rule, error) {
	return nil, errors.New("ERROR")
}

func (m mockErrorDatabase) DeactivateRule(ctx context.Context, ruleId primitive.ObjectID) (*database.Rule, error) {
	return nil, errors.New("ERROR")
}

var (
	ctx = context.Background()

	s = &server{
		ruleDb:  &mockDatabase{},
		storage: &mockStorage{},
		expr:    expression.NewExpression(&mockStorage{}),
	}
	sErr = &server{
		ruleDb:  &mockErrorDatabase{},
		storage: &mockErrorStorage{},
		expr:    expression.NewExpression(&mockErrorStorage{}),
	}

	incidentExistIncidentOpenedIncident = primitive.NewObjectID()
	isIncidentExistwasIncident          = primitive.NewObjectID()
	ruleIsNotActive                     = primitive.NewObjectID()
)

func TestNewIncidentServer(t *testing.T) {
	t.Run("Should: not nil", func(t *testing.T) {
		assert.NotNil(t, NewIncidentServer(&mockStorage{}, &mockDatabase{}))
	})
}

func TestServer_ActivateRule(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		_, err := s.ActivateRule(ctx, &apiPb.RuleIdRequest{
			RuleId: "",
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, err := sErr.ActivateRule(ctx, &apiPb.RuleIdRequest{
			RuleId: primitive.NewObjectID().Hex(),
		})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		_, err := s.ActivateRule(ctx, &apiPb.RuleIdRequest{
			RuleId: primitive.NewObjectID().Hex(),
		})
		assert.NoError(t, err)
	})
}

func TestServer_DeactivateRule(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		_, err := s.DeactivateRule(ctx, &apiPb.RuleIdRequest{
			RuleId: "",
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, err := sErr.DeactivateRule(ctx, &apiPb.RuleIdRequest{
			RuleId: primitive.NewObjectID().Hex(),
		})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		_, err := s.DeactivateRule(ctx, &apiPb.RuleIdRequest{
			RuleId: primitive.NewObjectID().Hex(),
		})
		assert.NoError(t, err)
	})
}

func TestServer_CreateRule(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		_, err := sErr.CreateRule(ctx, &apiPb.CreateRuleRequest{
			OwnerId:   "",
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, err := sErr.CreateRule(ctx, &apiPb.CreateRuleRequest{
			OwnerId:   primitive.NewObjectID().Hex(),
			OwnerType: apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT,
			Rule:      "len(Last(10)) > 0",
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, err := s.CreateRule(ctx, &apiPb.CreateRuleRequest{
			OwnerId:   primitive.NewObjectID().Hex(),
			OwnerType: apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT,
			Rule:      "wrongRule",
		})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		_, err := s.CreateRule(ctx, &apiPb.CreateRuleRequest{
			OwnerId:   primitive.NewObjectID().Hex(),
			OwnerType: apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT,
			Rule:      "len(Last(10)) > 0",
		})
		assert.NoError(t, err)
	})
}

func TestServer_GetRuleById(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		_, err := s.GetRuleById(ctx, &apiPb.RuleIdRequest{
			RuleId: "",
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, err := sErr.GetRuleById(ctx, &apiPb.RuleIdRequest{
			RuleId: primitive.NewObjectID().Hex(),
		})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		_, err := s.GetRuleById(ctx, &apiPb.RuleIdRequest{
			RuleId: primitive.NewObjectID().Hex(),
		})
		assert.NoError(t, err)
	})
}

func TestServer_GetRulesByOwnerId(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		_, err := s.GetRulesByOwnerId(ctx, &apiPb.GetRulesByOwnerIdRequest{
			OwnerId: "",
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, err := s.GetRulesByOwnerId(ctx, &apiPb.GetRulesByOwnerIdRequest{
			OwnerId: primitive.NewObjectID().Hex(),
		})
		assert.NoError(t, err)
	})
}

func TestServer_RemoveRule(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		_, err := s.RemoveRule(ctx, &apiPb.RuleIdRequest{
			RuleId: "",
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, err := sErr.RemoveRule(ctx, &apiPb.RuleIdRequest{
			RuleId: primitive.NewObjectID().Hex(),
		})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		_, err := s.RemoveRule(ctx, &apiPb.RuleIdRequest{
			RuleId: primitive.NewObjectID().Hex(),
		})
		assert.NoError(t, err)
	})
}

func TestServer_ProcessRecordFromStorage(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		_, err := s.ProcessRecordFromStorage(ctx, &apiPb.StorageRecord{
			Record: &apiPb.StorageRecord_Agent{},
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		_, err := sErr.ProcessRecordFromStorage(ctx, &apiPb.StorageRecord{
			Record: &apiPb.StorageRecord_Agent{
				Agent: &apiPb.Metric{
					AgentId: primitive.NewObjectID().Hex(),
				},
			},
		})
		assert.Error(t, err)
	})
	t.Run("Should: return error", func(t *testing.T) {
		sWithErrorStorage := &server{
			ruleDb:  &mockDatabase{},
			storage: &mockErrorStorage{},
			expr:    expression.NewExpression(&mockStorage{}),
		}
		_, err := sWithErrorStorage.ProcessRecordFromStorage(ctx, &apiPb.StorageRecord{
			Record: &apiPb.StorageRecord_Agent{
				Agent: &apiPb.Metric{
					AgentId: primitive.NewObjectID().Hex(),
				},
			},
		})
		assert.Error(t, err)
	})
	// isIncidentExist(incident) && isIncidentOpened(incident) && !wasIncident
	t.Run("Should: return error", func(t *testing.T) {
		_, err := s.ProcessRecordFromStorage(ctx, &apiPb.StorageRecord{
			Record: &apiPb.StorageRecord_Agent{
				Agent: &apiPb.Metric{
					AgentId: incidentExistIncidentOpenedIncident.Hex(),
				},
			},
		})
		assert.Error(t, err)
	})
	// !isIncidentExist(incident) && wasIncident
	t.Run("Should: return error", func(t *testing.T) {
		_, err := s.ProcessRecordFromStorage(ctx, &apiPb.StorageRecord{
			Record: &apiPb.StorageRecord_Agent{
				Agent: &apiPb.Metric{
					AgentId: isIncidentExistwasIncident.Hex(),
				},
			},
		})
		assert.Error(t, err)
	})

	t.Run("Should: return no error", func(t *testing.T) {
		_, err := s.ProcessRecordFromStorage(ctx, &apiPb.StorageRecord{
			Record: &apiPb.StorageRecord_Agent{
				Agent: &apiPb.Metric{
					AgentId: ruleIsNotActive.Hex(),
				},
			},
		})
		assert.NoError(t, err)
	})
}

func TestServer_CloseIncident(t *testing.T) {
	t.Run("Should: return no error", func(t *testing.T) {
		_, err := s.CloseIncident(ctx, &apiPb.IncidentIdRequest{})
		assert.NoError(t, err)
	})
}

func TestServer_StudyIncident(t *testing.T) {
	t.Run("Should: return no error", func(t *testing.T) {
		_, err := s.StudyIncident(ctx, &apiPb.IncidentIdRequest{})
		assert.NoError(t, err)
	})
}

func Test_getOwnerTypeAndId(t *testing.T) {
	t.Run("Should: return no error", func(t *testing.T) {
		_, _, err := getOwnerTypeAndId(&apiPb.StorageRecord{
			Record: &apiPb.StorageRecord_Scheduler{
				Scheduler: &apiPb.Scheduler{},
			},
		})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		_, _, err := getOwnerTypeAndId(&apiPb.StorageRecord{
			Record: &apiPb.StorageRecord_Scheduler{
				Scheduler: &apiPb.Scheduler{
					Id: primitive.NewObjectID().Hex(),
				},
			},
		})
		assert.NoError(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		_, _, err := getOwnerTypeAndId(&apiPb.StorageRecord{
			Record: &apiPb.StorageRecord_Agent{
				Agent: &apiPb.Metric{},
			},
		})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		_, _, err := getOwnerTypeAndId(&apiPb.StorageRecord{
			Record: &apiPb.StorageRecord_Transaction{
				Transaction: &apiPb.TransactionInfo{},
			},
		})
		assert.Error(t, err)
	})
	t.Run("Should: return no error", func(t *testing.T) {
		_, _, err := getOwnerTypeAndId(&apiPb.StorageRecord{
			Record: &apiPb.StorageRecord_Transaction{
				Transaction: &apiPb.TransactionInfo{
					ApplicationId: primitive.NewObjectID().Hex(),
				},
			},
		})
		assert.NoError(t, err)
	})
}

func Test_isIncidentOpened(t *testing.T) {
	t.Run("Should: return false", func(t *testing.T) {
		assert.False(t, isIncidentOpened(nil))
	})
}

func TestServer_tryCloseIncident(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		assert.Error(t, sErr.tryCloseIncident(ctx, true, &apiPb.Incident{}))
	})
}
