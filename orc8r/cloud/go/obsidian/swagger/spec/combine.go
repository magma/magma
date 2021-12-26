/*
 Copyright 2020 The Magma Authors.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package spec

import (
	"regexp"
	"sort"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"magma/orc8r/cloud/go/swagger"
)

// Combine multiple Swagger specs, giving precedence to the "common" spec.
// First error contains merge warnings for overwritten fields and
// incompatible Swagger specs.
// This custom-built functionality mirrors the "official" implementation.
// See: https://github.com/go-openapi/analysis/blob/master/mixin.go
func Combine(yamlCommon string, yamlSpecs []string) (string, error, error) {
	warnings := &multierror.Error{}

	common, specs, errs := unmarshalToSwagger(yamlCommon, yamlSpecs)
	if errs != nil {
		warnings = multierror.Append(warnings, errs)
	}

	combined, errs := combine(common, specs)
	if errs != nil {
		warnings = multierror.Append(warnings, errs)
	}

	out, err := marshalToYAML(combined)
	if err != nil {
		return "", nil, err
	}

	return out, warnings.ErrorOrNil(), nil
}

func combine(common swagger.Spec, specs []swagger.Spec) (swagger.Spec, error) {
	var tags [][]swagger.TagDefinition
	var paths []map[string]interface{}
	var responses []map[string]interface{}
	var parameters []map[string]interface{}
	var definitions []map[string]interface{}

	for _, s := range specs {
		tags = append(tags, s.Tags)
		paths = append(paths, s.Paths)
		responses = append(responses, s.Responses)
		parameters = append(parameters, s.Parameters)
		definitions = append(definitions, s.Definitions)
	}

	errs := &multierror.Error{}

	var combined = swagger.Spec{
		Swagger:             common.Swagger,
		Info:                common.Info,
		BasePath:            common.BasePath,
		Consumes:            common.Consumes,
		Produces:            common.Produces,
		Schemes:             common.Schemes,
		SecurityDefinitions: common.SecurityDefinitions,
		Security:            common.Security,
		Tags:                combineTags(common.Tags, tags, errs),
		Paths:               combineSubSpec(common.Paths, paths, "paths", errs),
		Responses:           combineSubSpec(common.Responses, responses, "responses", errs),
		Parameters:          combineSubSpec(common.Parameters, parameters, "parameters", errs),
		Definitions:         combineSubSpec(common.Definitions, definitions, "definitions", errs),
	}
	return combined, errs.ErrorOrNil()
}

// unmarshalToSwagger converts a list of specs and a common spec
// from YAML format to Swagger structs.
// Returned error is a list of errors from incompatible swagger specs.
func unmarshalToSwagger(yamlCommon string, yamlSpecs []string) (swagger.Spec, []swagger.Spec, error) {
	errs := &multierror.Error{}

	editedYAMLSpecs := makeAllYAMLReferencesLocal(yamlSpecs)
	specs := unmarshalManyFromYAML(editedYAMLSpecs, errs)

	common, err := unmarshalFromYAML(yamlCommon)
	if err != nil {
		return swagger.Spec{}, nil, err
	}

	return common, specs, errs.ErrorOrNil()
}

// unmarshalManyFromYAML maps the passed strings to their respective
// Swagger specs.
func unmarshalManyFromYAML(yamlSpecs []string, errs error) []swagger.Spec {
	var specs []swagger.Spec
	for _, yamlSpec := range yamlSpecs {
		s, err := unmarshalFromYAML(yamlSpec)
		if err != nil {
			errs = multierror.Append(errs, err)
		} else {
			specs = append(specs, s)
		}
	}
	return specs
}

// unmarshalFromYAML unmarshals the passed string to a Swagger spec.
func unmarshalFromYAML(yamlSpec string) (swagger.Spec, error) {
	spec := swagger.Spec{}
	err := yaml.Unmarshal([]byte(yamlSpec), &spec)
	if err != nil {
		return swagger.Spec{}, err
	}
	return spec, nil
}

// marshalToYAML marshals the passed Swagger spec to a YAML-formatted string.
func marshalToYAML(spec swagger.Spec) (string, error) {
	d, err := yaml.Marshal(&spec)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

// makeAllYAMLReferencesLocal rewrites all cross-file references to local,
// intra-file references.
// E.g.
//	- Before: `$ref: './orc8r-swagger-common.yml#/responses/UnexpectedError'`
//	- After:  `$ref: '#/responses/UnexpectedError'`
func makeAllYAMLReferencesLocal(yamlContents []string) []string {
	var rewritten []string
	// Match on any yml reference to file_name_here.foo.bar#/baz
	// and change those references to #/baz (strip the prefix)
	// e.g. $ref: 'foo_bar_baz.blah#/asdf' -> $ref: '#/asdf'
	ymlRefRe := regexp.MustCompile(`(\$ref:\s+)["']?.*(#/[^"'\s]+)["']?`)
	for _, yamlContent := range yamlContents {
		rewritten = append(rewritten, ymlRefRe.ReplaceAllString(yamlContent, "$1'$2'"))
	}
	return rewritten
}

func combineSubSpec(common map[string]interface{}, others []map[string]interface{}, name string, errs error) map[string]interface{} {
	combinedSpec := map[string]interface{}{}
	for _, cfg := range others {
		merge(combinedSpec, cfg, name, errs)
	}
	merge(combinedSpec, common, name, errs) // prefer common spec's fields
	return combinedSpec
}

func combineTags(common []swagger.TagDefinition, others [][]swagger.TagDefinition, errs error) []swagger.TagDefinition {
	combinedTagsByName := map[string]string{}
	for _, tags := range others {
		mergeTags(combinedTagsByName, tags, errs)
	}
	mergeTags(combinedTagsByName, common, errs) // prefer common tags

	var uniq []swagger.TagDefinition
	for name := range combinedTagsByName {
		t := swagger.TagDefinition{Name: name, Description: combinedTagsByName[name]}
		uniq = append(uniq, t)
	}
	sort.Slice(uniq, func(i, j int) bool { return uniq[i].Name < uniq[j].Name })
	return uniq
}

// merge b's contents into a, recording merge warnings to errs.
func merge(a, b map[string]interface{}, fieldName string, errs error) {
	for k, v := range b {
		if _, ok := a[k]; ok {
			errs = multierror.Append(errs, errors.Errorf("overwriting spec key '%s' in field '%s'", k, fieldName))
		}
		a[k] = v
	}
}

// mergeTags merges b's contents into a, recording merge warnings to errs.
func mergeTags(a map[string]string, b []swagger.TagDefinition, errs error) {
	for _, tag := range b {
		if _, ok := a[tag.Name]; ok {
			errs = multierror.Append(errs, errors.Errorf("overwriting tag '%s'", tag.Name))
		}
		a[tag.Name] = tag.Description
	}
}
