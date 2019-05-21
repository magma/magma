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
	"regexp"
	"sort"
	"strings"

	"magma/orc8r/cloud/go/services/metricsd/graphite/exporters"

	"github.com/go-graphite/carbonapi/pkg/parser"
)

const (
	// matches graphite metric name. Alphanumeric plus '.' and '_'
	metricNameRegexString = `[a-zA-Z_\.\+\*\d]+`

	// URLs use ';' as a delimiter between parameters so we use commas
	urlTagDelimiter = ","

	seriesByTagFuncName = "seriesByTag"
)

var (
	simpleTagRegexString = fmt.Sprintf("%s=[~]?%s", metricNameRegexString, metricNameRegexString)
	SimpleTagRegex       = regexp.MustCompile(simpleTagRegexString)

	// regex to match an list of tags
	tagListRegexString = fmt.Sprintf(",%s", simpleTagRegexString)

	//regex to match a metric name followed by an optional list of tags
	BasicQueryRegex = regexp.MustCompile(fmt.Sprintf("^(%s)(%s)*$", metricNameRegexString, tagListRegexString))
)

// RestrictQuery takes a graphite query string and replaces all seriesLists with
// a function call to `seriesByTag()` where the arguments include the series
// name, any other tags in the series, and networkID which ensures that only
// metrics belonging to that network can be retrieved.
func RestrictQuery(queryStr, networkID string) (string, error) {
	if BasicQueryRegex.MatchString(queryStr) {
		return restrictSeries(queryStr, networkID).ToString(), nil
	}
	query, _, err := parser.Parse(queryStr)
	if err != nil {
		return queryStr, fmt.Errorf("Could not parse query: %v", err)
	}
	if query.IsFunc() {
		query = restrictFunction(query, networkID)
		return query.ToString(), nil
	}
	if strings.HasPrefix(query.Target(), seriesByTagFuncName) {
		return restrictSeriesByTagFunc(query.Target(), networkID).ToString(), nil
	}
	query = restrictSeries(query.Target(), networkID)
	return query.ToString(), nil
}

// Builds a seriesByTag function for a given series and set of tags. Wraps the
// name in a regex that anchors the beginning and end to ensure that no side
// effects of the regex function are encountered.
func restrictSeries(series, networkID string) parser.Expr {
	name, tags := splitQueryNameTags(series)
	tags.Insert(exporters.NetworkTagName, networkID)

	var keys []string
	tags.Insert("name", fmt.Sprintf("~^%s$", name))
	for key := range tags {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var tagArgs []interface{}
	for _, key := range keys {
		tagArgs = append(tagArgs, parser.NewValueExpr(fmt.Sprintf("'%s=%s'", key, tags[key])))
	}
	return parser.NewExpr(seriesByTagFuncName, tagArgs...)
}

func restrictFunction(function parser.Expr, networkID string) parser.Expr {
	var newFuncArgs []parser.Expr
	var argsStrings []string
	for _, arg := range function.Args() {
		if arg.IsFunc() {
			newFunc := restrictFunction(arg, networkID)
			newFuncArgs = append(newFuncArgs, newFunc)
			argsStrings = append(argsStrings, newFunc.ToString())
		} else if arg.IsName() {
			if strings.HasPrefix(arg.Target(), seriesByTagFuncName) {
				return restrictSeriesByTagFunc(arg.Target(), networkID)
			}
			newSeries := restrictSeries(arg.Target(), networkID)
			seriesTarget := parser.NewTargetExpr(newSeries.ToString())
			newFuncArgs = append(newFuncArgs, seriesTarget)
			argsStrings = append(argsStrings, seriesTarget.Target())
		} else if arg.IsConst() {
			newFuncArgs = append(newFuncArgs, arg)
			argsStrings = append(argsStrings, fmt.Sprintf("%g", arg.FloatValue()))
		} else if arg.IsString() {
			newFuncArgs = append(newFuncArgs, arg)
			argsStrings = append(argsStrings, fmt.Sprintf(`"%s"`, arg.StringValue()))
		}
	}
	retFunc := parser.NewExprTyped(function.Target(), newFuncArgs)
	retFunc.SetRawArgs(strings.Join(argsStrings, ","))
	return retFunc
}

// Adds a 'networkID' tag to the arguments of this 'seriesByTag' function. If
// 'networkID' already exists, this enforces that the correct ID is used
func restrictSeriesByTagFunc(funcStr, networkID string) parser.Expr {
	tagList := SimpleTagRegex.FindAllString(funcStr, -1)
	tags := fillTagSetFromList(tagList)
	tags.Insert(exporters.NetworkTagName, networkID)

	var tagArgs []interface{}
	for _, tag := range tags.SortedTags() {
		tagArgs = append(tagArgs, parser.NewValueExpr(fmt.Sprintf(`'%s'`, tag)))
	}
	return parser.NewExpr(seriesByTagFuncName, tagArgs...)
}

func splitQueryNameTags(query string) (string, exporters.TagSet) {
	tagList := SimpleTagRegex.FindAllString(query, -1)
	if len(tagList) == 0 {
		return query, exporters.TagSet{}
	}
	nameEndIdx := strings.Index(query, urlTagDelimiter)
	name := query[:nameEndIdx]
	tags := fillTagSetFromList(tagList)
	return name, tags
}

func fillTagSetFromList(tagList []string) exporters.TagSet {
	tags := make(exporters.TagSet)
	for _, tag := range tagList {
		equalsIndex := strings.Index(tag, "=")
		key := tag[:equalsIndex]
		val := tag[equalsIndex+1:]
		tags.Insert(key, val)
	}
	return tags
}
