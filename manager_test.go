//
// Copyright 2016 Pedro Salgado
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package dao_test

import (
	"github.com/stretchr/testify/suite"

	"github.com/steenzout/go-dao"
	"github.com/steenzout/go-dao/mock"
	mock_sql "github.com/steenzout/go-mock-database/sql"
)

const (
	// DAO_MOCK unique identifier for the mock.MockDAO interface implementations.
	DAO_MOCK = "mock.MockDAO"
)

type TestManager struct {
	dao.BaseManager
}

func (m *TestManager) CreateMockDAO(ctx *dao.Context) (mock.MockDAO, error) {
	f, found := m.Factories[DAO_MOCK]
	if !found {
		return nil, nil
	}
	return f.(*mock.MockFactory).CreateMockDAO(ctx)
}

// ManagerSuite test suite for the Manager struct.
type ManagerTestSuite struct {
	suite.Suite
}

func (s *ManagerTestSuite) Test() {
	ctx, err := manager.StartTransaction()
	if err != nil {
		s.Fail(err.Error())
	}
	s.NotNil(ctx)
	defer manager.EndTransaction(ctx)

	mockDAO, err := manager.CreateMockDAO(ctx)
	if err != nil {
		s.Fail(err.Error())
	}

	mockDAO.MockSomething()

	err = manager.CommitTransaction(ctx)
	if err != nil {
		s.Fail(err.Error())
	}
}

var factory *mock.MockFactory
var manager TestManager
var _ dao.Manager = (*TestManager)(nil)

func init() {
	manager = TestManager{*dao.NewBaseManager()}

	mtx := mock_sql.Tx{}
	mtx.On("CommitTransaction").Return(nil)
	ds1 := mock.NewDataSource().(*mock.DataSource)
	ds1.On("Begin").Return(&mtx, nil)

	manager.RegisterDataSource("mock", ds1)

	factory = mock.NewFactory()
	factory.SetDataSource(ds1)

	manager.RegisterFactory(DAO_MOCK, factory)
}
