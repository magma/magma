/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package security

import (
	"fmt"
	"reflect"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/graphite/exporters"

	"github.com/stretchr/testify/assert"
)

const (
	basicTargetSeries  = `metric_name`
	longTargetSeries   = `top.middle.metric`
	regexSeries        = `top.*.metric`
	complexRegexSeries = `cpu[0-9]-{user,system}.value`

	absFunc        = `absolute(metric_name)`
	scaleFunc      = `scale(metric_name, 10)`
	aliasFunc      = `alias(metric_name, "new_name")`
	aggregateFunc  = `aggregateWithWildcards(metric_name, "sum", 1)`
	scaleWithRegex = `scale(cpu[0-9]-{user,system}.value, 10)`

	composedFunc         = `sumSeries(alias(metric_name, "new_name"))`
	composedMultiArgFunc = `sumSeries(alias(metric_name, "new_name"), absolute(metric_name))`
	composedMixedFunc    = `sumSeries(absolute(metric_name), scale(metric_name, 10), metric_name)`
	multiLayerFunc       = `sumSeries(scale(alias(metric_name, "new_name"), 10))`

	seriesByTagFunc      = `seriesByTag('name=metric_name')`
	seriesByTagMultiFunc = `seriesByTag('name=metric_name','networkID=testNetwork')`

	basicQuery               = `metric_name`
	basicQueryWithTags       = `metric_name,tag1=val1,tag2=val2`
	basicQueryWithNetworkTag = `metric_name,tag1=val1,networkID=otherNetwork`
	basicQueryWithRegexTag   = `metric_name,tag1=~val_id_.*`

	testNetwork = "testNetwork"

	seriesByTagTemplate = `seriesByTag('name=~^%s$','networkID=testNetwork')`

	seriesByTagInsideFunction = `sumSeries(seriesByTag('name=metric_name'))`
	notEqualOperator          = `seriesByTag('name!=metric_name')`
	notEqualRegexOperator     = `seriesByTag('name!=~metric_.*')`
)

var (
	restrictedBasicSeries        = fmt.Sprintf(seriesByTagTemplate, basicTargetSeries)
	restrictedLongSeries         = fmt.Sprintf(seriesByTagTemplate, longTargetSeries)
	restrictedRegexSeries        = fmt.Sprintf(seriesByTagTemplate, regexSeries)
	restrictedComplexRegexSeries = fmt.Sprintf(seriesByTagTemplate, complexRegexSeries)

	restrictedAbsFunc            = fmt.Sprintf("absolute(%s)", restrictedBasicSeries)
	restrictedScaleFunc          = fmt.Sprintf("scale(%s,10)", restrictedBasicSeries)
	restrictedAliasFunc          = fmt.Sprintf(`alias(%s,"new_name")`, restrictedBasicSeries)
	restrictedAggregateFunc      = fmt.Sprintf(`aggregateWithWildcards(%s,"sum",1)`, restrictedBasicSeries)
	restrictedScaleWithRegexFunc = fmt.Sprintf(`scale(%s,10)`, restrictedComplexRegexSeries)

	restrictedComposedFunc         = fmt.Sprintf(`sumSeries(alias(%s,"new_name"))`, restrictedBasicSeries)
	restrictedComposedMultiArgFunc = fmt.Sprintf(`sumSeries(alias(%s,"new_name"),absolute(%s))`, restrictedBasicSeries, restrictedBasicSeries)
	restrictedComposedMixedFunc    = fmt.Sprintf(`sumSeries(absolute(%s),scale(%s,10),%s)`, restrictedBasicSeries, restrictedBasicSeries, restrictedBasicSeries)
	restrictedMultiLayerFunc       = fmt.Sprintf(`sumSeries(scale(alias(%s,"new_name"),10))`, restrictedBasicSeries)

	restrictedSeriesByTagFunc = fmt.Sprintf(`seriesByTag('name=metric_name','networkID=%s')`, testNetwork)

	restrictedBasicQuery               = fmt.Sprintf(`seriesByTag('name=~^metric_name$','networkID=%s')`, testNetwork)
	restrictedBasicQueryWithTags       = fmt.Sprintf(`seriesByTag('name=~^metric_name$','networkID=%s','tag1=val1','tag2=val2')`, testNetwork)
	restrictedBasicQueryWithNetworkTag = fmt.Sprintf(`seriesByTag('name=~^metric_name$','networkID=%s','tag1=val1')`, testNetwork)
	restrictedBasicQueryWithRegexTag   = fmt.Sprintf(`seriesByTag('name=~^metric_name$','networkID=%s','tag1=~val_id_.*')`, testNetwork)

	restrictedSeriesByTagInsideFunction = fmt.Sprintf(`sumSeries(seriesByTag('name=metric_name','networkID=%s'))`, testNetwork)
	restrictedNotEqualOperator          = fmt.Sprintf(`seriesByTag('name!=metric_name','networkID=%s')`, testNetwork)
	restrictedNotEqualRegexOperator     = fmt.Sprintf(`seriesByTag('name!=~metric_.*','networkID=%s')`, testNetwork)
)

type RestrictorTestCase struct {
	input     string
	networkID string
	expected  string
}

func NewBasicRestrictorTestCase(input string, expectedFunc string) RestrictorTestCase {
	return RestrictorTestCase{
		input:     input,
		networkID: testNetwork,
		expected:  expectedFunc,
	}
}

func (c RestrictorTestCase) runTest(t *testing.T) {
	query, err := RestrictQuery(c.input, c.networkID)
	assert.NoError(t, err)

	if ok := reflect.DeepEqual(c.expected, query); !ok {
		t.Errorf("Query %s did not match expected output", c.input)
		fmt.Printf("Expected: %+v\n", c.expected)
		fmt.Printf("Actual  : %+v\n", query)
		fmt.Printf("\n\n")
	}
}

func TestRestrictQuery(t *testing.T) {
	testCases := []RestrictorTestCase{
		NewBasicRestrictorTestCase(basicTargetSeries, restrictedBasicSeries),
		NewBasicRestrictorTestCase(longTargetSeries, restrictedLongSeries),
		NewBasicRestrictorTestCase(regexSeries, restrictedRegexSeries),
		NewBasicRestrictorTestCase(complexRegexSeries, restrictedComplexRegexSeries),

		NewBasicRestrictorTestCase(absFunc, restrictedAbsFunc),
		NewBasicRestrictorTestCase(scaleFunc, restrictedScaleFunc),
		NewBasicRestrictorTestCase(aliasFunc, restrictedAliasFunc),
		NewBasicRestrictorTestCase(aggregateFunc, restrictedAggregateFunc),
		NewBasicRestrictorTestCase(scaleWithRegex, restrictedScaleWithRegexFunc),

		NewBasicRestrictorTestCase(composedFunc, restrictedComposedFunc),
		NewBasicRestrictorTestCase(composedMultiArgFunc, restrictedComposedMultiArgFunc),
		NewBasicRestrictorTestCase(composedMixedFunc, restrictedComposedMixedFunc),
		NewBasicRestrictorTestCase(multiLayerFunc, restrictedMultiLayerFunc),

		NewBasicRestrictorTestCase(seriesByTagFunc, restrictedSeriesByTagFunc),
		NewBasicRestrictorTestCase(seriesByTagMultiFunc, restrictedSeriesByTagFunc),

		NewBasicRestrictorTestCase(basicQuery, restrictedBasicQuery),
		NewBasicRestrictorTestCase(basicQueryWithTags, restrictedBasicQueryWithTags),
		NewBasicRestrictorTestCase(basicQueryWithNetworkTag, restrictedBasicQueryWithNetworkTag),
		NewBasicRestrictorTestCase(basicQueryWithRegexTag, restrictedBasicQueryWithRegexTag),

		NewBasicRestrictorTestCase(seriesByTagInsideFunction, restrictedSeriesByTagInsideFunction),
		NewBasicRestrictorTestCase(notEqualOperator, restrictedNotEqualOperator),
		NewBasicRestrictorTestCase(notEqualRegexOperator, restrictedNotEqualRegexOperator),
	}
	for _, c := range testCases {
		c.runTest(t)
	}
}

func TestSplitQueryNameTags(t *testing.T) {
	query := "metric,tag1=val1"
	name, tags := splitQueryNameTags(query)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, tags.String(), ";tag1=val1")
	assert.Equal(t, "metric", name)

	query = "metric,tag1=val1,tag2=val2"
	name, tags = splitQueryNameTags(query)
	assert.Equal(t, 2, len(tags))
	assert.Equal(t, tags.String(), ";tag1=val1;tag2=val2")
	assert.Equal(t, "metric", name)

	query = "metric_name"
	name, tags = splitQueryNameTags(query)
	assert.Equal(t, exporters.TagSet{}, tags)
	assert.Equal(t, "metric_name", name)
}
