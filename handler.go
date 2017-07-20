// Package dao provides a data access object library.
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
package dao

// TransactionFunc definition of a function wrapped in a database transaction context.
type TransactionFunc func(m Manager, ctx *Context, args ...interface{}) (interface{}, error)

// Process wrap database transaction handling around given TransactionFunc.
func Process(m Manager, f TransactionFunc, args ...interface{}) (interface{}, error) {
	ctx, err := m.StartTransaction()
	if err != nil {
		return nil, err
	}
	defer m.EndTransaction(ctx)

	v, err := f(m, ctx, args...)
	if err != nil {
		m.RollbackTransaction(ctx)
		return nil, err
	}

	if err = m.CommitTransaction(ctx); err != nil {
		m.RollbackTransaction(ctx)
		return nil, err
	}

	return v, nil
}

// TransactionMapFunc definition of a function wrapped in a database transaction context.
type TransactionMapFunc func(m Manager, ctx *Context, data map[string]interface{}) (interface{}, error)

// ProcessMap wrap database transaction handling around given TransactionMapFunc.
func ProcessMap(m Manager, f TransactionMapFunc, data []map[string]interface{}) ([]interface{}, error) {
	ctx, err := m.StartTransaction()
	if err != nil {
		return nil, err
	}
	defer m.EndTransaction(ctx)

	var models []interface{}
	for _, val := range data {
		v, err := f(m, ctx, val)
		if err != nil {
			m.RollbackTransaction(ctx)
			return nil, err
		}
		models = append(models, v)
	}

	if err = m.CommitTransaction(ctx); err != nil {
		m.RollbackTransaction(ctx)
		return nil, err
	}

	return models, nil
}
