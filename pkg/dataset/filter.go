/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dataset

import (
	"strings"
)

import (
	"github.com/arana-db/arana/pkg/proto"
)

var (
	_ proto.Dataset = (*FilterDataset)(nil)
	_ proto.Dataset = (*FilterDatasetPrefix)(nil)
)

type PredicateFunc func(proto.Row) bool

type FilterDataset struct {
	proto.Dataset
	Predicate PredicateFunc
}

func (f FilterDataset) Next() (proto.Row, error) {
	if f.Predicate == nil {
		return f.Dataset.Next()
	}

	row, err := f.Dataset.Next()
	if err != nil {
		return nil, err
	}

	if !f.Predicate(row) {
		return f.Next()
	}

	return row, nil
}

type FilterDatasetPrefix struct {
	proto.Dataset
	Predicate PredicateFunc
	Prefix    string
}

func (f FilterDatasetPrefix) Next() (proto.Row, error) {
	if f.Predicate == nil {
		return f.Dataset.Next()
	}

	row, err := f.Dataset.Next()
	if err != nil {
		return nil, err
	}

	var (
		fields, _ = f.Fields()
		values    = make([]proto.Value, len(fields))
	)

	if err = row.Scan(values); err != nil {
		return nil, err
	}

	if strings.HasPrefix(values[0].(string), f.Prefix) {
		return f.Next()
	}

	if !f.Predicate(row) {
		return f.Next()
	}

	return row, nil
}
