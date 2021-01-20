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

package swagger

import (
	"github.com/golang/glog"

	"regexp"
	"sort"

	"magma/orc8r/cloud/go/tools/swaggergen/generate"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Load specs from YAML format to Swagger structs.
func Load(common string, configs []string) ([]generate.SwaggerSpec, generate.SwaggerSpec, error) {
	specs, err := loadSwaggerSpecs(configs)
	if err != nil {
		return nil, generate.SwaggerSpec{}, err
	}

	commonSpec, err := loadCommonSpec(common)
	if err != nil {
		return nil, generate.SwaggerSpec{}, err
	}

	return specs, commonSpec, nil
}

/*
CombineStr is a wrapper function around the function Combine()
which serves to combine multiple swagger specs whose contents
are passed in. The output is the merged copy of the swagger spec
*/
func CombineStr(common string, configs []string) (string, error) {
	cfgs, commonCfg, err := Load(common, configs)
	if err != nil {
		return "Error: Load of Swagger Specs Failed", err
	}
	
	combinedSpec, warnings := Combine(commonCfg, cfgs)
	glog.V(2).Infof("Combining Swagger Specs Together \n")
	merrs, _ := warnings.(*multierror.Error)
	if len(merrs.Errors) != 0 {
		glog.Warningf("%+v merge errors: %+v", len(merrs.Errors), merrs)
	}

	out, err := MarshalToYAML(combinedSpec)
	if err != nil {
		return "Error: Marshaling of Combined Swagger Spec Failed", err
	}

	return out, merrs.ErrorOrNil()
}

// Combine multiple Swagger specs, giving precedence to the "common" spec.
// Returned "error" contains warnings for any overwritten fields.
//
// This custom-built functionality mirrors the "official" implementation.
// See: https://github.com/go-openapi/analysis/blob/master/mixin.go
func Combine(common generate.SwaggerSpec, specs []generate.SwaggerSpec) (generate.SwaggerSpec, error) {
	var tags [][]generate.TagDefinition
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

	out := generate.SwaggerSpec{
		Swagger:  common.Swagger,
		Info:     common.Info,
		BasePath: common.BasePath,
		Consumes: common.Consumes,
		Produces: common.Produces,
		Schemes:  common.Schemes,

		Tags:        combineTags(common.Tags, tags, errs),
		Paths:       combineSubSpec(common.Paths, paths, "paths", errs),
		Responses:   combineSubSpec(common.Responses, responses, "responses", errs),
		Parameters:  combineSubSpec(common.Parameters, parameters, "parameters", errs),
		Definitions: combineSubSpec(common.Definitions, definitions, "definitions", errs),
	}

	return out, errs.ErrorOrNil()
}

// loadSwaggerSpecs unmarshals all of the passed in Swagger Specs YAML strings
// to Swagger structs.
func loadSwaggerSpecs(configs []string) ([]generate.SwaggerSpec, error){
	editedContents := makeAllYAMLReferencesLocal(configs)

	specs, err := unmarshalManyFromYAML(editedContents)
	if err != nil {
		return nil, err
	}

	return specs, err
}

// loadCommonSpec unmarshals the common Swagger Specs YAML to struct.
func loadCommonSpec(common string) (generate.SwaggerSpec, error) {
	spec, err := UnmarshalFromYAML(common)
	if err != nil {
		return generate.SwaggerSpec{}, err
	}
	return spec ,err
}

// UnmarshalFromYAML maps the passed strings to their respective
// Swagger specs.
func unmarshalManyFromYAML(swaggerYAMLs []string) ([]generate.SwaggerSpec, error) {
	var specs []generate.SwaggerSpec
	for _, swaggerYAML := range swaggerYAMLs {
		s, err := UnmarshalFromYAML(swaggerYAML)
		if err != nil {
			return nil, err
		}
		specs = append(specs, s)
	}
	return specs, nil
}

// UnmarshalFromYAML unmarshals the passed string to a Swagger spec.
func UnmarshalFromYAML(swaggerYAML string) (generate.SwaggerSpec, error) {
	spec := generate.SwaggerSpec{}
	err := yaml.Unmarshal([]byte(swaggerYAML), &spec)
	return spec, err
}

// MarshalToYAML marshals the passed Swagger spec to a YAML-formatted string.
func MarshalToYAML(spec generate.SwaggerSpec) (string, error) {
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
	ymlRefRe := regexp.MustCompile(`(\$ref:\s*)['"].+(#/.+)['"]`)
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

func combineTags(common []generate.TagDefinition, others [][]generate.TagDefinition, errs error) []generate.TagDefinition {
	combinedTagsByName := map[string]string{}
	for _, tags := range others {
		mergeTags(combinedTagsByName, tags, errs)
	}
	mergeTags(combinedTagsByName, common, errs) // prefer common tags

	var uniq []generate.TagDefinition
	for name := range combinedTagsByName {
		t := generate.TagDefinition{Name: name, Description: combinedTagsByName[name]}
		uniq = append(uniq, t)
	}
	sort.Slice(uniq, func(i, j int) bool { return uniq[i].Name < uniq[j].Name })
	return uniq
}

// merge b's contents into a, recording merge warnings to errs.
func merge(a, b map[string]interface{}, fieldName string, errs error) {
	for k, v := range b {
		if _, ok := a[k]; ok {
			multierror.Append(errs, errors.Errorf("overwriting spec key '%s' in field '%s'", k, fieldName))
		}
		a[k] = v
	}
}

// mergeTags merges b's contents into a, recording merge warnings to errs.
func mergeTags(a map[string]string, b []generate.TagDefinition, errs error) {
	for _, tag := range b {
		if _, ok := a[tag.Name]; ok {
			multierror.Append(errs, errors.Errorf("overwriting tag '%s'", tag.Name))
		}
		a[tag.Name] = tag.Description
	}
}