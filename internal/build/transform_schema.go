//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
	"gopkg.in/yaml.v3"
)

// ApiGroup defines a group of API paths to be output to a specific file
type ApiGroup struct {
	Name       string
	PathPrefix []string
	Filename   string
	Group      string
}

// ApiGroups defines the output groups
var ApiGroups = []ApiGroup{
	{Name: "Fleet", PathPrefix: []string{"/api/fleet"}, Filename: "fleet.yml", Group: "fleet"},
	{Name: "DataViews", PathPrefix: []string{"/api/data_views"}, Filename: "data_views.yml", Group: "dataviews"},
	{Name: "Alerting", PathPrefix: []string{"/api/alerting"}, Filename: "alerting.yml", Group: "alerting"},
	{Name: "APM", PathPrefix: []string{"/api/apm"}, Filename: "apm.yml", Group: "apm"},
	{Name: "Cases", PathPrefix: []string{"/api/cases"}, Filename: "cases.yml", Group: "cases"},
	{Name: "Connectors", PathPrefix: []string{"/api/actions"}, Filename: "connectors.yml", Group: "connectors"},
	{Name: "DetectionEngine", PathPrefix: []string{"/api/detection_engine", "/api/exception_list", "/api/exceptions/shared"}, Filename: "detection_engine.yml", Group: "detection"},
	{Name: "Roles", PathPrefix: []string{"/api/security/role"}, Filename: "roles.yml", Group: "roles"},
	{Name: "ML", PathPrefix: []string{"/api/ml"}, Filename: "ml.yml", Group: "ml"},
	{Name: "SavedObjects", PathPrefix: []string{"/api/saved_objects", "/api/encrypted_saved_objects"}, Filename: "saved_objects.yml", Group: "savedobjects"},
	{Name: "SecurityAIAssistant", PathPrefix: []string{"/api/security_ai_assistant"}, Filename: "security_ai_assistant.yml", Group: "securityaiassistant"},
	{Name: "Endpoint", PathPrefix: []string{"/api/endpoint"}, Filename: "endpoint.yml", Group: "endpoint"},
	{Name: "OSquery", PathPrefix: []string{"/api/osquery"}, Filename: "osquery.yml", Group: "osquery"},
	{Name: "Spaces", PathPrefix: []string{"/api/spaces"}, Filename: "spaces.yml", Group: "spaces"},
	{Name: "Status", PathPrefix: []string{"/api/status"}, Filename: "status.yml", Group: "status"},
	{Name: "SLOs", PathPrefix: []string{"/s/{spaceId}/api/observability/slos", "/s/{spaceId}/internal/observability/slos/_definitions"}, Filename: "slos.yml", Group: "slo"},
	{Name: "Timeline", PathPrefix: []string{"/api/timeline", "/api/timelines", "/api/note", "/api/pinned_event"}, Filename: "timeline.yml", Group: "timeline"},
	{Name: "EntityEngine", PathPrefix: []string{"/api/entity_store", "/api/risk_score", "/api/asset_criticality"}, Filename: "entity_engine.yml", Group: "entityanalytics"},
	{Name: "Lists", PathPrefix: []string{"/api/lists"}, Filename: "lists.yml", Group: "lists"},
	{Name: "Streams", PathPrefix: []string{"/api/streams"}, Filename: "streams.yml", Group: "streams"},
	{Name: "Uptime", PathPrefix: []string{"/api/uptime"}, Filename: "uptime.yml", Group: "uptime"},
	{Name: "TaskManager", PathPrefix: []string{"/api/task_manager"}, Filename: "task_manager.yml", Group: "taskmanager"},
}

var numberOfPaths int

// ComponentUsage tracks which components are used by which groups
type ComponentUsage struct {
	SchemaUsage   map[string]map[string]bool
	ParamUsage    map[string]map[string]bool
	ResponseUsage map[string]map[string]bool
}

func main() {
	_inFile := flag.String("i", "", "input file")
	_oAPICodeGenVersion := flag.String("v", "2.4.1", "Open API Code Generator Version")
	flag.Parse()

	inFile := *_inFile
	oAPICodeGenVersion := *_oAPICodeGenVersion
	oAPICodeGenURL := "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v" + oAPICodeGenVersion

	if inFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	bytes, err := os.ReadFile(inFile)
	if err != nil {
		log.Fatalf("failed to read file %q: %v", inFile, err)
	}

	var schema Schema
	err = yaml.Unmarshal(bytes, &schema)
	if err != nil {
		log.Fatalf("failed to unmarshal schema from %q: %v", inFile, err)
	}

	// Run each transform
	for _, fn := range transformers {
		fn(&schema)
	}

	processedPaths := make(map[string]bool)

	type oAPICodeGen struct {
		Name  string
		Group string
	}

	tmplContent, err := os.ReadFile("cfg.tmpl")
	if err != nil {
		log.Fatalf("Error reading template file: %v", err)
	}

	// Process each API group
	for _, group := range ApiGroups {
		copySchema := schema

		// Create a new paths map filtering only the paths for this group
		copySchema.Paths = filterPathsByPrefixes(schema.Paths, group.PathPrefix)

		if len(copySchema.Paths) == 0 {
			continue
		}

		for pathURL := range copySchema.Paths {
			processedPaths[pathURL] = true
		}

		// Create oas file
		outFile := filepath.Join("./", group.Filename)
		saveFile(copySchema, outFile)

		// Create Open API Code Gen config.yaml
		nameLower := strings.ToLower(group.Name)
		oAPIData := oAPICodeGen{
			Name:  nameLower,
			Group: group.Group,
		}

		tmpl, err := template.New("configYAML").Parse(string(tmplContent))
		if err != nil {
			panic(err)
		}

		configFilePath := "./" + nameLower + ".config.yaml"
		configOutFile, err := os.Create(configFilePath)
		if err != nil {
			log.Fatalf("Error creating output file: %v", err)
		}
		err = tmpl.Execute(configOutFile, oAPIData)
		if err != nil {
			panic(err)
		}
		configOutFile.Close()

		cmd := exec.Command("go", "run", oAPICodeGenURL,
			"-config", configFilePath, outFile)

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Panicf("Error executing oapi-codegen for %s: %v\n%s", group.Group, err, string(output))
		}

		generatedCodePath := fmt.Sprintf("../../kbapi/models.%s.gen.go", nameLower)
		generatedCode, err := os.ReadFile(generatedCodePath)
		if err != nil {
			log.Fatalf("Error reading generated code: %v", err)
		}

		// Process the generated code
		processedCode := postProcessGeneratedCode(string(generatedCode))

		// Write the processed code back to the file
		err = os.WriteFile(generatedCodePath, []byte(processedCode), 0644)
		if err != nil {
			log.Fatalf("Error writing processed code: %v", err)
		}

		err = os.Remove(outFile)
		if err != nil {
			log.Printf("Warning: Failed to remove %s: %v", outFile, err)
		}

		err = os.Remove(configFilePath)
		if err != nil {
			log.Printf("Warning: Failed to remove %s: %v", configFilePath, err)
		}
	}

	// Create schema for all other paths and fail
	miscSchema := Schema{
		Version:    schema.Version,
		Info:       schema.Info,
		Servers:    schema.Servers,
		Security:   schema.Security,
		Tags:       schema.Tags,
		Components: schema.Components,
		Paths:      make(map[string]*Path),
	}

	for pathUrl, pathInfo := range schema.Paths {
		if !processedPaths[pathUrl] {
			miscSchema.Paths[pathUrl] = pathInfo
		}
	}

	if len(miscSchema.Paths) > 0 {
		outFile := filepath.Join("./", "misc.yml")
		saveFile(miscSchema, outFile)
		log.Println("Remaining paths:")
		for pathURL := range miscSchema.Paths {
			log.Printf("  %s", pathURL)
		}
		log.Panicf("Saved %s with %d remaining paths", outFile, len(miscSchema.Paths))
	}
}

func postProcessGeneratedCode(code string) string {
	// Parse the code preserving comments
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		log.Fatalf("Error parsing generated code: %v", err)
		return code
	}

	// Create a new AST with only declarations we want to keep
	newF := &ast.File{
		Name:    f.Name,
		Decls:   []ast.Decl{},
		Package: f.Package,
	}

	// Track positions of kept declarations to determine which comments to keep
	keptPositions := make(map[token.Pos]token.Pos)

	// Filter declarations
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.IMPORT {
				// Keep imports for now, will be fixed later
				newF.Decls = append(newF.Decls, d)
				keptPositions[d.Pos()] = d.End()
			} else if d.Tok == token.TYPE {
				shouldKeep := true
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						if ts.Name.Name == "ClientOption" ||
							ts.Name.Name == "Client" ||
							ts.Name.Name == "HttpRequestDoer" ||
							ts.Name.Name == "ClientWithResponses" ||
							ts.Name.Name == "ClientWithResponsesInterface" ||
							ts.Name.Name == "ClientInterface" ||
							ts.Name.Name == "RequestEditorFn" {
							shouldKeep = false
							break
						}
						if strings.HasSuffix(ts.Name.Name, "Response") {
							if st, ok := ts.Type.(*ast.StructType); ok {
								// Find the JSON200 field
								for i, field := range st.Fields.List {
									if len(field.Names) > 0 && field.Names[0].Name == "JSON200" {
										// Replace the type with the contents of JSON200
										if starExpr, ok := field.Type.(*ast.StarExpr); ok {
											// Case 1: JSON200 is a struct literal
											if structType, ok := starExpr.X.(*ast.StructType); ok {
												// Replace the entire struct with the contents of JSON200
												ts.Type = structType
											} else if ident, ok := starExpr.X.(*ast.Ident); ok {
												// Case 2: JSON200 is a reference to another type
												// Just replace this type with the referenced type
												ts.Type = &ast.StarExpr{
													X: ident,
												}
											} else if selector, ok := starExpr.X.(*ast.SelectorExpr); ok {
												// Case 3: JSON200 is a reference to a type in another package
												ts.Type = &ast.StarExpr{
													X: selector,
												}
											}

											// Remove Body, HTTPResponse, and JSON* fields
											var newFields []*ast.Field
											for j, f := range st.Fields.List {
												if j != i && // Skip JSON200
													(len(f.Names) == 0 || // Keep anonymous fields
														(f.Names[0].Name != "Body" &&
															f.Names[0].Name != "HTTPResponse" &&
															!strings.HasPrefix(f.Names[0].Name, "JSON"))) {
													newFields = append(newFields, f)
												}
											}
											st.Fields.List = newFields
										}
										break
									}
								}
							}
						}
					}
				}
				if shouldKeep {
					newF.Decls = append(newF.Decls, d)
					keptPositions[d.Pos()] = d.End()
				}
			} else if d.Tok == token.VAR {
				newF.Decls = append(newF.Decls, d)
				keptPositions[d.Pos()] = d.End()
			}
			// Skip all other declarations (CONST, etc.)
		case *ast.FuncDecl:
			// Skip all function declarations
		}
	}

	// Filter comments to only keep those associated with kept declarations
	var newCommentGroups []*ast.CommentGroup
	for _, cg := range f.Comments {
		// Keep comment if it's a package-level comment
		if cg.End() < f.Package {
			newCommentGroups = append(newCommentGroups, cg)
			continue
		}

		// Check if comment is associated with a kept declaration
		shouldKeep := false
		cgPos := cg.Pos()

		for declPos, declEnd := range keptPositions {
			// Comment should be right before the declaration or within it
			if (cgPos < declPos && cg.End() >= declPos-1) ||
				(cgPos >= declPos && cgPos <= declEnd) {
				shouldKeep = true
				break
			}
		}

		if shouldKeep {
			newCommentGroups = append(newCommentGroups, cg)
		}
	}

	newF.Comments = newCommentGroups

	// Format the filtered AST back to code
	var buf bytes.Buffer
	something := printer.Config{Mode: printer.UseSpaces, Tabwidth: 8}
	something.Fprint(&buf, fset, newF)

	// Use goimports to fix imports
	processedCode, err := imports.Process("", buf.Bytes(), nil)
	if err != nil {
		log.Fatalf("Error fixing imports: %v", err)
		return buf.String()
	}

	return string(processedCode)
}

func filterPathsByPrefixes(paths map[string]*Path, prefixes []string) map[string]*Path {
	result := make(map[string]*Path)
	for pathURL, pathInfo := range paths {
		for _, prefix := range prefixes {
			if strings.HasPrefix(pathURL, prefix) {
				result[pathURL] = pathInfo
				break // Once we match a prefix, no need to check others
			}
		}
	}
	return result
}

// pathExists checks if path exists.
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// saveFile marshal and writes obj to path.
func saveFile(obj any, path string) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(obj); err != nil {
		log.Fatalf("failed to marshal to file %q: %v", path, err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0664); err != nil {
		log.Fatalf("failed to write file %q: %v", path, err)
	}
}

// ============================================================================

type Schema struct {
	Version    string           `yaml:"openapi"`
	Info       Map              `yaml:"info"`
	Servers    []Map            `yaml:"servers,omitempty"`
	Security   []Map            `yaml:"security,omitempty"`
	Tags       []Map            `yaml:"tags,omitempty"`
	Paths      map[string]*Path `yaml:"paths"`
	Components Map              `yaml:"components,omitempty"`
}

func (s Schema) GetPath(path string) *Path {
	return s.Paths[path]
}

func (s Schema) MustGetPath(path string) *Path {
	p := s.GetPath(path)
	if p == nil {
		log.Panicf("Path not found: %q", path)
	}
	return p
}

// ============================================================================

type Path struct {
	Parameters []Map `yaml:"parameters,omitempty"`
	Get        Map   `yaml:"get,omitempty"`
	Post       Map   `yaml:"post,omitempty"`
	Put        Map   `yaml:"put,omitempty"`
	Delete     Map   `yaml:"delete,omitempty"`
}

func (p Path) Endpoints(yield func(key string, endpoint Map) bool) {
	if p.Get != nil {
		yield("get", p.Get)
	}
	if p.Post != nil {
		yield("post", p.Post)
	}
	if p.Put != nil {
		yield("put", p.Put)
	}
	if p.Delete != nil {
		yield("delete", p.Delete)
	}
}

func (p Path) GetEndpoint(method string) Map {
	switch method {
	case "get":
		return p.Get
	case "post":
		return p.Post
	case "put":
		return p.Put
	case "delete":
		return p.Delete
	default:
		log.Panicf("Unhandled method: %q", method)
	}
	return nil
}

func (p Path) MustGetEndpoint(method string) Map {
	endpoint := p.GetEndpoint(method)
	if endpoint == nil {
		log.Panicf("Method not found: %q", method)
	}
	return endpoint
}

func (p *Path) SetEndpoint(method string, endpoint Map) {
	switch method {
	case "get":
		p.Get = endpoint
	case "post":
		p.Post = endpoint
	case "put":
		p.Put = endpoint
	case "delete":
		p.Delete = endpoint
	default:
		log.Panicf("Invalid method %q", method)
	}
}

// ============================================================================

type Map map[string]any

func (m Map) Keys() []string {
	keys := slices.Collect(maps.Keys(m))
	slices.Sort(keys)
	return keys
}

func (m Map) Has(key string) bool {
	_, ok := m.Get(key)
	return ok
}

func (m Map) Get(key string) (any, bool) {
	rootKey, subKeys, found := strings.Cut(key, ".")
	if found {
		switch t := m[rootKey].(type) {
		case Map:
			return t.Get(subKeys)
		case map[string]any:
			return Map(t).Get(subKeys)
		case Slice:
			return t.Get(subKeys)
		case []any:
			return Slice(t).Get(subKeys)
		default:
			rootKey = key
		}
	}

	value, ok := m[rootKey]
	return value, ok
}

func (m Map) MustGet(key string) any {
	v, ok := m.Get(key)
	if !ok {
		log.Panicf("%q not found", key)
	}
	return v
}

func (m Map) GetSlice(key string) (Slice, bool) {
	value, ok := m.Get(key)
	if !ok {
		return nil, false
	}

	switch t := value.(type) {
	case Slice:
		return t, true
	case []any:
		return t, true
	}

	log.Panicf("%q is not a slice", key)
	return nil, false
}

func (m Map) MustGetSlice(key string) Slice {
	v, ok := m.GetSlice(key)
	if !ok {
		log.Panicf("%q not found", key)
	}
	return v
}

func (m Map) GetMap(key string) (Map, bool) {
	value, ok := m.Get(key)
	if !ok {
		return nil, false
	}

	switch t := value.(type) {
	case Map:
		return t, true
	case map[string]any:
		return t, true
	}

	log.Panicf("%q is not a map", key)
	return nil, false
}

func (m Map) MustGetMap(key string) Map {
	v, ok := m.GetMap(key)
	if !ok {
		log.Panicf("%q not found", key)
	}
	return v
}

func (m Map) Set(key string, value any) {
	rootKey, subKeys, found := strings.Cut(key, ".")
	if found {
		if v, ok := m[rootKey]; ok {
			switch t := v.(type) {
			case Slice:
				t.Set(subKeys, value)
			case []any:
				Slice(t).Set(subKeys, value)
			case Map:
				t.Set(subKeys, value)
			case map[string]any:
				Map(t).Set(subKeys, value)
			}
		} else {
			subMap := Map{}
			subMap.Set(subKeys, value)
			m[rootKey] = subMap
		}
	} else {
		m[rootKey] = value
	}
}

func (m Map) Move(src string, dst string) {
	value := m.MustGet(src)
	m.Set(dst, value)
	m.Delete(src)
}

func (m Map) Delete(key string) bool {
	rootKey, subKeys, found := strings.Cut(key, ".")
	if found {
		if v, ok := m[rootKey]; ok {
			switch t := v.(type) {
			case Slice:
				return t.Delete(subKeys)
			case []any:
				return Slice(t).Delete(subKeys)
			case Map:
				return t.Delete(subKeys)
			case map[string]any:
				return Map(t).Delete(subKeys)
			}
		}
	} else {
		delete(m, rootKey)
		return true
	}
	return false
}

func (m Map) MustDelete(key string) {
	if !m.Delete(key) {
		log.Panicf("%q not found", key)
	}
}

func (m Map) CreateRef(schema *Schema, name string, key string) Map {
	refTarget := m.MustGet(key) // Check the full path
	refPath := fmt.Sprintf("schemas.%s", name)
	refValue := Map{"$ref": fmt.Sprintf("#/components/schemas/%s", name)}

	// If the component schema already exists and is not the same, panic
	writeComponent := true
	if existing, ok := schema.Components.Get(refPath); ok {
		if reflect.DeepEqual(refTarget, existing) {
			writeComponent = false
		} else {
			log.Panicf("Component schema key already in use and not an exact duplicate: %q", refPath)
			return nil
		}
	}

	var parent any
	var childKey string
	// Get the parent of the refTarget
	i := strings.LastIndex(key, ".")
	if i == -1 {
		parent = m
		childKey = key
	} else {
		parent = m.MustGet(key[:i])
		childKey = key[i+1:]
	}

	doMap := func(target Map, key string) {
		if writeComponent {
			schema.Components.Set(refPath, target.MustGet(key))
		}
		target.Set(key, refValue)
	}

	doSlice := func(target Slice, key string) {
		index := target.atoi(key)
		if writeComponent {
			schema.Components.Set(refPath, target[index])
		}
		target[index] = refValue
	}

	switch t := parent.(type) {
	case map[string]any:
		doMap(Map(t), childKey)
	case Map:
		doMap(t, childKey)
	case []any:
		doSlice(Slice(t), childKey)
	case Slice:
		doSlice(t, childKey)
	default:
		log.Panicf("Cannot create a ref of target type %T at %q", parent, key)
	}

	return refValue
}

func (m Map) Iterate(iteratee func(key string, node Map)) {
	joinPath := func(existing string, next string) string {
		if existing == "" {
			return next
		} else {
			return fmt.Sprintf("%s.%s", existing, next)
		}
	}
	joinIndex := func(existing string, next int) string {
		if existing == "" {
			return fmt.Sprintf("%d", next)
		} else {
			return fmt.Sprintf("%s.%d", existing, next)
		}
	}

	var iterate func(key string, val any)
	iterate = func(key string, val any) {
		switch tval := val.(type) {
		case []any:
			iterate(key, Slice(tval))
		case Slice:
			for i, v := range tval {
				iterate(joinIndex(key, i), v)
			}
		case map[string]any:
			iterate(key, Map(tval))
		case Map:
			for _, k := range tval.Keys() {
				iterate(joinPath(key, k), tval[k])
			}
			iteratee(key, tval)
		}
	}

	iterate("", m)
}

// ============================================================================

type Slice []any

func (s Slice) Get(key string) (any, bool) {
	rootKey, subKeys, found := strings.Cut(key, ".")
	index := s.atoi(rootKey)

	if found {
		switch t := s[index].(type) {
		case Slice:
			return t.Get(subKeys)
		case []any:
			return Slice(t).Get(subKeys)
		case Map:
			return t.Get(subKeys)
		case map[string]any:
			return Map(t).Get(subKeys)
		}
	}

	value := s[index]
	return value, true
}

func (s Slice) GetMap(key string) (Map, bool) {
	value, ok := s.Get(key)
	if !ok {
		return nil, false
	}

	switch t := value.(type) {
	case Map:
		return t, true
	case map[string]any:
		return t, true
	}

	log.Panicf("%q is not a map", key)
	return nil, false
}

func (s Slice) MustGetMap(key string) Map {
	v, ok := s.GetMap(key)
	if !ok {
		log.Panicf("%q not found", key)
	}
	return v
}

func (s Slice) Set(key string, value any) {
	rootKey, subKeys, found := strings.Cut(key, ".")
	index := s.atoi(rootKey)
	if found {
		v := s[index]
		switch t := v.(type) {
		case Slice:
			t.Set(subKeys, value)
		case []any:
			Slice(t).Set(subKeys, value)
		case Map:
			t.Set(subKeys, value)
		case map[string]any:
			Map(t).Set(subKeys, value)
		}
	} else {
		s[index] = value
	}
}

func (s Slice) Delete(key string) bool {
	rootKey, subKeys, found := strings.Cut(key, ".")
	index := s.atoi(rootKey)
	if found {
		item := (s)[index]
		switch t := item.(type) {
		case Slice:
			return t.Delete(subKeys)
		case []any:
			return Slice(t).Delete(subKeys)
		case Map:
			return t.Delete(subKeys)
		case map[string]any:
			return Map(t).Delete(subKeys)
		}
	} else {
		log.Panicf("Unable to delete from slice directly")
		return true
	}
	return false
}

func (s Slice) Contains(value string) bool {
	for _, v := range s {
		s, ok := v.(string)
		if !ok {
			continue
		}
		if value == s {
			return true
		}
	}

	return false
}

func (s Slice) atoi(key string) int {
	index, err := strconv.Atoi(key)
	if err != nil {
		log.Panicf("Failed to parse slice index key %q: %v", key, err)
	}
	if index < 0 || index >= len(s) {
		log.Panicf("Slice index is out of bounds (%d, target slice len: %d)", index, len(s))
	}
	return index
}

// ============================================================================

type TransformFunc func(schema *Schema)

var transformers = []TransformFunc{
	transformRemoveKbnXsrf,
	transformRemoveApiVersionParam,
	transformSimplifyContentType,
	transformAddMisingDescriptions,
	transformFleetPaths,
	transformRemoveEnums,
	transformRemoveExamples,
	transformRemoveUnusedComponents,
	transformRemoveDuplicateTags,
	fixArrayItemsDefinitions,
	removeTechnicalPreviewPaths,
	transformFixVersionFields,
	transformRemoveProblematicPaths,
	transformRunMessageEmail,
	transformEntityAnalyticsTypes,
	transformGetAllSpaces,
}

// transformRemoveKbnXsrf removes the kbn-xsrf header as it	is already applied
// in the client.
func transformRemoveKbnXsrf(schema *Schema) {
	removeKbnXsrf := func(node any) bool {
		param := node.(Map)
		if v, ok := param["name"]; ok {
			name := v.(string)
			if strings.HasSuffix(name, "kbn_xsrf") || strings.HasSuffix(name, "kbn-xsrf") {
				return true
			}
		}
		// Data_views_kbn_xsrf, Saved_objects_kbn_xsrf, etc
		if v, ok := param["$ref"]; ok {
			ref := v.(string)
			if strings.HasSuffix(ref, "kbn_xsrf") || strings.HasSuffix(ref, "kbn-xsrf") {
				return true
			}
		}
		return false
	}

	for _, pathInfo := range schema.Paths {
		for _, endpoint := range pathInfo.Endpoints {
			if params, ok := endpoint.GetSlice("parameters"); ok {
				params = slices.DeleteFunc(params, removeKbnXsrf)
				endpoint["parameters"] = params
			}
		}
	}
}

// transformRemoveApiVersionParam removes the Elastic API Version query
// parameter header.
func transformRemoveApiVersionParam(schema *Schema) {
	removeApiVersion := func(node any) bool {
		param := node.(Map)
		if name, ok := param["name"]; ok && name == "elastic-api-version" {
			return true
		}
		return false
	}

	for _, pathInfo := range schema.Paths {
		for _, endpoint := range pathInfo.Endpoints {
			if params, ok := endpoint.GetSlice("parameters"); ok {
				params = slices.DeleteFunc(params, removeApiVersion)
				endpoint["parameters"] = params
			}
		}
	}
}

// transformSimplifyContentType simplifies Content-Type headers such as
// 'application/json; Elastic-Api-Version=2023-10-31' by stripping everything
// after the ';'.
func transformSimplifyContentType(schema *Schema) {
	simplifyContentType := func(fields Map) {
		if content, ok := fields.GetMap("content"); ok {
			for key := range content {
				newKey, _, found := strings.Cut(key, ";")
				if found {
					content.Move(key, newKey)
				}
			}
		}
	}

	for _, pathInfo := range schema.Paths {
		for _, endpoint := range pathInfo.Endpoints {
			if req, ok := endpoint.GetMap("requestBody"); ok {
				simplifyContentType(req)
			}
			if resp, ok := endpoint.GetMap("responses"); ok {
				for code := range resp {
					simplifyContentType(resp.MustGetMap(code))
				}
			}
		}
	}

	if responses, ok := schema.Components.GetMap("responses"); ok {
		for key := range responses {
			resp := responses.MustGetMap(key)
			simplifyContentType(resp)
		}
	}
}

// transformAddMisingDescriptions adds descriptions to each path missing one.
func transformAddMisingDescriptions(schema *Schema) {
	for _, pathInfo := range schema.Paths {
		for _, endpoint := range pathInfo.Endpoints {
			responses := endpoint.MustGetMap("responses")
			for code := range responses {
				response := responses.MustGetMap(code)
				if _, ok := response["description"]; !ok {
					response["description"] = ""
				}
			}
		}
	}
}

// transformKibanaPaths fixes the Kibana paths.
func transformKibanaPaths(schema *Schema) {
	// Convert any paths needing it to /s/{spaceId} variants
	spaceIdPaths := []string{
		"/api/data_views",
		"/api/data_views/data_view",
		"/api/data_views/data_view/{viewId}",
	}

	// Add a spaceId parameter if not already present
	if _, ok := schema.Components.Get("parameters.spaceId"); !ok {
		schema.Components.Set("parameters.spaceId", Map{
			"in":          "path",
			"name":        "spaceId",
			"description": "An identifier for the space. If `/s/` and the identifier are omitted from the path, the default space is used.",
			"required":    true,
			"schema":      Map{"type": "string", "example": "default"},
		})
	}

	for _, path := range spaceIdPaths {
		pathInfo := schema.Paths[path]
		schema.Paths[fmt.Sprintf("/s/{spaceId}%s", path)] = pathInfo
		delete(schema.Paths, path)

		// Add the spaceId parameter
		param := Map{"$ref": "#/components/parameters/spaceId"}
		for _, endpoint := range pathInfo.Endpoints {
			if params, ok := endpoint.GetSlice("parameters"); ok {
				params = append(params, param)
				endpoint.Set("parameters", params)
			} else {
				params = Slice{param}
				endpoint.Set("parameters", params)
			}
		}
	}

	// Data views
	// https://github.com/elastic/kibana/blob/main/src/plugins/data_views/server/rest_api_routes/schema.ts

	dataViewsPath := schema.MustGetPath("/s/{spaceId}/api/data_views")

	dataViewsPath.Get.CreateRef(schema, "get_data_views_response_item", "responses.200.content.application/json.schema.properties.data_view.items")

	schema.Components.CreateRef(schema, "Data_views_data_view_response_object_inner", "schemas.Data_views_data_view_response_object.properties.data_view")
	schema.Components.CreateRef(schema, "Data_views_sourcefilter_item", "schemas.Data_views_sourcefilters.items")
	schema.Components.CreateRef(schema, "Data_views_runtimefieldmap_script", "schemas.Data_views_runtimefieldmap.properties.script")

	schema.Components.Set("schemas.Data_views_fieldformats.additionalProperties", Map{
		"$ref": "#/components/schemas/Data_views_fieldformat",
	})
	schema.Components.Set("schemas.Data_views_fieldformat", Map{
		"type": "object",
		"properties": Map{
			"id":     Map{"type": "string"},
			"params": Map{"$ref": "#/components/schemas/Data_views_fieldformat_params"},
		},
	})
	schema.Components.Set("schemas.Data_views_fieldformat_params", Map{
		"type": "object",
		"properties": Map{
			"pattern":                Map{"type": "string"},
			"urlTemplate":            Map{"type": "string"},
			"labelTemplate":          Map{"type": "string"},
			"inputFormat":            Map{"type": "string"},
			"outputFormat":           Map{"type": "string"},
			"outputPrecision":        Map{"type": "integer"},
			"includeSpaceWithSuffix": Map{"type": "boolean"},
			"useShortSuffix":         Map{"type": "boolean"},
			"timezone":               Map{"type": "string"},
			"fieldType":              Map{"type": "string"},
			"colors": Map{
				"type":  "array",
				"items": Map{"$ref": "#/components/schemas/Data_views_fieldformat_params_color"},
			},
			"fieldLength": Map{"type": "integer"},
			"transform":   Map{"type": "string"},
			"lookupEntries": Map{
				"type":  "array",
				"items": Map{"$ref": "#/components/schemas/Data_views_fieldformat_params_lookup"},
			},
			"unknownKeyValue": Map{"type": "string"},
			"type":            Map{"type": "string"},
			"width":           Map{"type": "integer"},
			"height":          Map{"type": "integer"},
		},
	})
	schema.Components.Set("schemas.Data_views_fieldformat_params_color", Map{
		"type": "object",
		"properties": Map{
			"range":      Map{"type": "string"},
			"regex":      Map{"type": "string"},
			"text":       Map{"type": "string"},
			"background": Map{"type": "string"},
		},
	})
	schema.Components.Set("schemas.Data_views_fieldformat_params_lookup", Map{
		"type": "object",
		"properties": Map{
			"key":   Map{"type": "string"},
			"value": Map{"type": "string"},
		},
	})

	schema.Components.CreateRef(schema, "Data_views_create_data_view_request_object_inner", "schemas.Data_views_create_data_view_request_object.properties.data_view")
	schema.Components.CreateRef(schema, "Data_views_update_data_view_request_object_inner", "schemas.Data_views_update_data_view_request_object.properties.data_view")
}

// transformFleetPaths fixes the fleet paths.
func transformFleetPaths(schema *Schema) {
	// Agent policies
	// https://github.com/elastic/kibana/blob/main/x-pack/plugins/fleet/common/types/models/agent_policy.ts
	// https://github.com/elastic/kibana/blob/main/x-pack/plugins/fleet/common/types/rest_spec/agent_policy.ts

	agentPoliciesPath := schema.MustGetPath("/api/fleet/agent_policies")
	agentPolicyPath := schema.MustGetPath("/api/fleet/agent_policies/{agentPolicyId}")

	agentPoliciesPath.Get.CreateRef(schema, "agent_policy", "responses.200.content.application/json.schema.properties.items.items")
	agentPoliciesPath.Post.CreateRef(schema, "agent_policy", "responses.200.content.application/json.schema.properties.item")
	agentPolicyPath.Get.CreateRef(schema, "agent_policy", "responses.200.content.application/json.schema.properties.item")
	agentPolicyPath.Put.CreateRef(schema, "agent_policy", "responses.200.content.application/json.schema.properties.item")

	// See: https://github.com/elastic/kibana/issues/197155
	// [request body.keep_monitoring_alive]: expected value of type [boolean] but got [null]
	// [request body.supports_agentless]: expected value of type [boolean] but got [null]
	// [request body.overrides]: expected value of type [boolean] but got [null]
	// [request body.required_versions]: definition for this key is missing"}
	for _, key := range []string{"keep_monitoring_alive", "supports_agentless", "overrides", "required_versions"} {
		agentPoliciesPath.Post.Set(fmt.Sprintf("requestBody.content.application/json.schema.properties.%s.x-omitempty", key), true)
		agentPolicyPath.Put.Set(fmt.Sprintf("requestBody.content.application/json.schema.properties.%s.x-omitempty", key), true)
	}

	// Enrollment api keys
	// https://github.com/elastic/kibana/blob/main/x-pack/plugins/fleet/common/types/models/enrollment_api_key.ts
	// https://github.com/elastic/kibana/blob/main/x-pack/plugins/fleet/common/types/rest_spec/enrollment_api_key.ts

	apiKeysPath := schema.MustGetPath("/api/fleet/enrollment_api_keys")
	apiKeysPath.Get.CreateRef(schema, "enrollment_api_key", "responses.200.content.application/json.schema.properties.items.items")

	// EPM
	// https://github.com/elastic/kibana/blob/main/x-pack/plugins/fleet/common/types/models/epm.ts
	// https://github.com/elastic/kibana/blob/main/x-pack/plugins/fleet/common/types/rest_spec/epm.ts

	packagesPath := schema.MustGetPath("/api/fleet/epm/packages")
	packagePath := schema.MustGetPath("/api/fleet/epm/packages/{pkgName}/{pkgVersion}")
	packagesPath.Get.CreateRef(schema, "package_list_item", "responses.200.content.application/json.schema.properties.items.items")
	packagePath.Get.CreateRef(schema, "package_info", "responses.200.content.application/json.schema.properties.item")

	// Server hosts
	// https://github.com/elastic/kibana/blob/main/x-pack/plugins/fleet/common/types/models/fleet_server_policy_config.ts
	// https://github.com/elastic/kibana/blob/main/x-pack/plugins/fleet/common/types/rest_spec/fleet_server_hosts.ts

	hostsPath := schema.MustGetPath("/api/fleet/fleet_server_hosts")
	hostPath := schema.MustGetPath("/api/fleet/fleet_server_hosts/{itemId}")

	hostsPath.Get.CreateRef(schema, "server_host", "responses.200.content.application/json.schema.properties.items.items")
	hostsPath.Post.CreateRef(schema, "server_host", "responses.200.content.application/json.schema.properties.item")
	hostPath.Get.CreateRef(schema, "server_host", "responses.200.content.application/json.schema.properties.item")
	hostPath.Put.CreateRef(schema, "server_host", "responses.200.content.application/json.schema.properties.item")

	// 8.6.2 regression
	// [request body.proxy_id]: definition for this key is missing
	// See: https://github.com/elastic/kibana/issues/197155
	hostsPath.Post.Set("requestBody.content.application/json.schema.properties.proxy_id.x-omitempty", true)
	hostPath.Put.Set("requestBody.content.application/json.schema.properties.proxy_id.x-omitempty", true)

	// Outputs
	// https://github.com/elastic/kibana/blob/main/x-pack/plugins/fleet/common/types/models/output.ts
	// https://github.com/elastic/kibana/blob/main/x-pack/plugins/fleet/common/types/rest_spec/output.ts

	outputByIdPath := schema.MustGetPath("/api/fleet/outputs/{outputId}")
	outputsPath := schema.MustGetPath("/api/fleet/outputs")

	outputsPath.Post.CreateRef(schema, "new_output_union", "requestBody.content.application/json.schema")
	outputByIdPath.Put.CreateRef(schema, "update_output_union", "requestBody.content.application/json.schema")
	outputsPath.Get.CreateRef(schema, "output_union", "responses.200.content.application/json.schema.properties.items.items")
	outputByIdPath.Get.CreateRef(schema, "output_union", "responses.200.content.application/json.schema.properties.item")
	outputsPath.Post.CreateRef(schema, "output_union", "responses.200.content.application/json.schema.properties.item")
	outputByIdPath.Put.CreateRef(schema, "output_union", "responses.200.content.application/json.schema.properties.item")

	for _, name := range []string{"output", "new_output", "update_output"} {
		// Ref each index in the anyOf union
		schema.Components.CreateRef(schema, fmt.Sprintf("%s_elasticsearch", name), fmt.Sprintf("schemas.%s_union.anyOf.0", name))
		schema.Components.CreateRef(schema, fmt.Sprintf("%s_remote_elasticsearch", name), fmt.Sprintf("schemas.%s_union.anyOf.1", name))
		schema.Components.CreateRef(schema, fmt.Sprintf("%s_logstash", name), fmt.Sprintf("schemas.%s_union.anyOf.2", name))
		schema.Components.CreateRef(schema, fmt.Sprintf("%s_kafka", name), fmt.Sprintf("schemas.%s_union.anyOf.3", name))

		// Extract child structs
		for _, typ := range []string{"elasticsearch", "remote_elasticsearch", "logstash", "kafka"} {
			schema.Components.CreateRef(schema, fmt.Sprintf("%s_shipper", name), fmt.Sprintf("schemas.%s_%s.properties.shipper", name, typ))
			schema.Components.CreateRef(schema, fmt.Sprintf("%s_ssl", name), fmt.Sprintf("schemas.%s_%s.properties.ssl", name, typ))
		}

		// Ideally just remove the "anyOf", however then we would need to make
		// refs for each of the "oneOf" options. So turn them into an "any" instead.
		// See: https://github.com/elastic/kibana/issues/197153
		/*
			anyOf:
			  - items: {}
			    type: array
			  - type: boolean
			  - type: number
			  - type: object
			  - type: string
			nullable: true
			oneOf:
			  - type: number
			  - not: {}
		*/

		props := schema.Components.MustGetMap(fmt.Sprintf("schemas.%s_kafka.properties", name))
		for _, key := range []string{"compression_level", "connection_type", "password", "username"} {
			props.Set(key, Map{})
		}
	}

	// Add the missing discriminator to the response union
	// See: https://github.com/elastic/kibana/issues/181994
	schema.Components.Set("schemas.output_union.discriminator", Map{
		"propertyName": "type",
		"mapping": Map{
			"elasticsearch":        "#/components/schemas/output_elasticsearch",
			"remote_elasticsearch": "#/components/schemas/output_remote_elasticsearch",
			"logstash":             "#/components/schemas/output_logstash",
			"kafka":                "#/components/schemas/output_kafka",
		},
	})

	for _, name := range []string{"new_output", "update_output"} {
		for _, typ := range []string{"elasticsearch", "remote_elasticsearch", "logstash", "kafka"} {
			// [request body.1.ca_sha256]: expected value of type [string] but got [null]"
			// See: https://github.com/elastic/kibana/issues/197155
			schema.Components.Set(fmt.Sprintf("schemas.%s_%s.properties.ca_sha256.x-omitempty", name, typ), true)

			// [request body.1.ca_trusted_fingerprint]: expected value of type [string] but got [null]
			// See: https://github.com/elastic/kibana/issues/197155
			schema.Components.Set(fmt.Sprintf("schemas.%s_%s.properties.ca_trusted_fingerprint.x-omitempty", name, typ), true)

			// 8.6.2 regression
			// [request body.proxy_id]: definition for this key is missing"
			// See: https://github.com/elastic/kibana/issues/197155
			schema.Components.Set(fmt.Sprintf("schemas.%s_%s.properties.proxy_id.x-omitempty", name, typ), true)
		}

		// [request body.1.shipper]: expected a plain object value, but found [null] instead
		// See: https://github.com/elastic/kibana/issues/197155
		schema.Components.Set(fmt.Sprintf("schemas.%s_shipper.x-omitempty", name), true)

		// [request body.1.ssl]: expected a plain object value, but found [null] instead
		// See: https://github.com/elastic/kibana/issues/197155
		schema.Components.Set(fmt.Sprintf("schemas.%s_ssl.x-omitempty", name), true)

	}

	for _, typ := range []string{"elasticsearch", "remote_elasticsearch", "logstash", "kafka"} {
		// strict_dynamic_mapping_exception: [1:345] mapping set to strict, dynamic introduction of [id] within [ingest-outputs] is not allowed"
		// See: https://github.com/elastic/kibana/issues/197155
		schema.Components.MustDelete(fmt.Sprintf("schemas.update_output_%s.properties.id", typ))
	}

	// Package policies
	// https://github.com/elastic/kibana/blob/main/x-pack/plugins/fleet/common/types/models/package_policy.ts
	// https://github.com/elastic/kibana/blob/main/x-pack/plugins/fleet/common/types/rest_spec/package_policy.ts

	epmPoliciesPath := schema.MustGetPath("/api/fleet/package_policies")
	epmPolicyPath := schema.MustGetPath("/api/fleet/package_policies/{packagePolicyId}")

	epmPoliciesPath.Get.CreateRef(schema, "package_policy", "responses.200.content.application/json.schema.properties.items.items")
	epmPoliciesPath.Post.CreateRef(schema, "package_policy", "responses.200.content.application/json.schema.properties.item")

	epmPoliciesPath.Post.Move("requestBody.content.application/json.schema.anyOf.1", "requestBody.content.application/json.schema") // anyOf.0 is the deprecated array format
	epmPolicyPath.Put.Move("requestBody.content.application/json.schema.anyOf.1", "requestBody.content.application/json.schema")    // anyOf.0 is the deprecated array format
	epmPoliciesPath.Post.CreateRef(schema, "package_policy_request", "requestBody.content.application/json.schema")
	epmPolicyPath.Put.CreateRef(schema, "package_policy_request", "requestBody.content.application/json.schema")

	epmPolicyPath.Get.CreateRef(schema, "package_policy", "responses.200.content.application/json.schema.properties.item")
	epmPolicyPath.Put.CreateRef(schema, "package_policy", "responses.200.content.application/json.schema.properties.item")

	schema.Components.CreateRef(schema, "package_policy_secret_ref", "schemas.package_policy.properties.secret_references.items")
	schema.Components.Move("schemas.package_policy.properties.inputs.anyOf.1", "schemas.package_policy.properties.inputs") // anyOf.0 is the deprecated array format

	schema.Components.CreateRef(schema, "package_policy_input", "schemas.package_policy.properties.inputs.additionalProperties")
	schema.Components.CreateRef(schema, "package_policy_input_stream", "schemas.package_policy_input.properties.streams.additionalProperties")

	schema.Components.CreateRef(schema, "package_policy_request_package", "schemas.package_policy_request.properties.package")
	schema.Components.CreateRef(schema, "package_policy_request_input", "schemas.package_policy_request.properties.inputs.additionalProperties")
	schema.Components.CreateRef(schema, "package_policy_request_input_stream", "schemas.package_policy_request_input.properties.streams.additionalProperties")

	// Simplify all of the vars
	schema.Components.Set("schemas.package_policy.properties.vars", Map{"type": "object"})
	schema.Components.Set("schemas.package_policy_input.properties.vars", Map{"type": "object"})
	schema.Components.Set("schemas.package_policy_input_stream.properties.vars", Map{"type": "object"})
	schema.Components.Set("schemas.package_policy_request.properties.vars", Map{"type": "object"})
	schema.Components.Set("schemas.package_policy_request_input.properties.vars", Map{"type": "object"})
	schema.Components.Set("schemas.package_policy_request_input_stream.properties.vars", Map{"type": "object"})

	// [request body.0.output_id]: expected value of type [string] but got [null]
	// [request body.1.output_id]: definition for this key is missing"
	// See: https://github.com/elastic/kibana/issues/197155
	schema.Components.Set("schemas.package_policy_request.properties.output_id.x-omitempty", true)
}

// transformRemoveEnums remove all enums.
func transformRemoveEnums(schema *Schema) {
	deleteEnumFn := func(key string, node Map) {
		if node.Has("enum") {
			delete(node, "enum")
		}
	}

	for _, pathInfo := range schema.Paths {
		for _, methInfo := range pathInfo.Endpoints {
			methInfo.Iterate(deleteEnumFn)
		}
	}
	schema.Components.Iterate(deleteEnumFn)
}

// transformRemoveExamples removes all examples.
func transformRemoveExamples(schema *Schema) {
	deleteExampleFn := func(key string, node Map) {
		if node.Has("example") {
			delete(node, "example")
		}
		if node.Has("examples") {
			delete(node, "examples")
		}
	}

	for _, pathInfo := range schema.Paths {
		for _, methInfo := range pathInfo.Endpoints {
			methInfo.Iterate(deleteExampleFn)
		}
	}
	schema.Components.Iterate(deleteExampleFn)
	schema.Components.Set("examples", Map{})
}

// transformAddOptionalPointersFlag adds a x-go-type-skip-optional-pointer
// flag to maps and arrays, since they are already nullable types.
func transformAddOptionalPointersFlag(schema *Schema) {
	addFlagFn := func(key string, node Map) {
		if node["type"] == "array" {
			node["x-go-type-skip-optional-pointer"] = true
		} else if node["type"] == "object" {
			if _, ok := node["properties"]; !ok {
				node["x-go-type-skip-optional-pointer"] = true
			}
		}
	}

	for _, pathInfo := range schema.Paths {
		for _, methInfo := range pathInfo.Endpoints {
			methInfo.Iterate(addFlagFn)
		}
	}
	schema.Components.Iterate(addFlagFn)
}

// transformRemoveUnusedComponents removes all unused schema components.
func transformRemoveUnusedComponents(schema *Schema) {
	var refs map[string]any
	collectRefsFn := func(key string, node Map) {
		if ref, ok := node["$ref"].(string); ok {
			i := strings.LastIndex(ref, "/")
			ref = ref[i+1:]
			refs[ref] = nil
		}
	}

	componentParams := schema.Components.MustGetMap("parameters")
	componentSchemas := schema.Components.MustGetMap("schemas")

	for {
		// Collect refs
		refs = make(map[string]any)
		for _, pathInfo := range schema.Paths {
			for _, methInfo := range pathInfo.Endpoints {
				methInfo.Iterate(collectRefsFn)
			}
		}
		schema.Components.Iterate(collectRefsFn)

		loop := false
		for key := range componentSchemas {
			if _, ok := refs[key]; !ok {
				delete(componentSchemas, key)
				loop = true
			}
		}
		for key := range componentParams {
			if _, ok := refs[key]; !ok {
				delete(componentParams, key)
				loop = true
			}
		}
		if !loop {
			break
		}
	}
}

func transformRemoveDuplicateTags(schema *Schema) {
	if schema.Tags == nil || len(schema.Tags) == 0 {
		return
	}

	seenTags := make(map[string]bool)
	var uniqueTags []Map

	for _, tag := range schema.Tags {
		name, ok := tag["name"].(string)
		if !ok {
			uniqueTags = append(uniqueTags, tag)
			continue
		}

		if _, exists := seenTags[name]; !exists {
			seenTags[name] = true
			uniqueTags = append(uniqueTags, tag)
		}
	}

	schema.Tags = uniqueTags
}

func transformFixRequiredPathParameters(schema *Schema) {
	for pathURL, pathInfo := range schema.Paths {
		pathParamRegex := regexp.MustCompile(`\{([^}]+)\}`)
		matches := pathParamRegex.FindAllStringSubmatch(pathURL, -1)

		if len(matches) == 0 {
			continue
		}

		// Create a set of parameter names we need to fix
		paramNamesToFix := make(map[string]bool)
		for _, match := range matches {
			if len(match) > 1 {
				paramNamesToFix[match[1]] = true
			}
		}

		// Function to ensure a parameter has required=true
		fixParameter := func(param map[string]interface{}) bool {
			if name, ok := param["name"].(string); ok {
				if paramNamesToFix[name] {
					param["required"] = true
					return true
				}
			}
			return false
		}

		// Fix endpoint-level parameters
		for _, method := range []string{"Get", "Post", "Put", "Delete"} {
			var endpoint map[string]interface{}

			switch method {
			case "Get":
				endpoint = pathInfo.Get
			case "Post":
				endpoint = pathInfo.Post
			case "Put":
				endpoint = pathInfo.Put
			case "Delete":
				endpoint = pathInfo.Delete
			}

			if endpoint == nil {
				continue
			}

			// Check for parameters in this endpoint
			if paramsAny, ok := endpoint["parameters"]; ok {
				if params, ok := paramsAny.([]interface{}); ok {
					for _, paramAny := range params {
						if param, ok := paramAny.(map[string]interface{}); ok {
							fixParameter(param)
						}
					}
				}
			}
		}
	}
}

func transformFixAllValidationIssues(schema *Schema) {
	// Fix required arrays to ensure they only contain strings
	fixRequiredArray := func(node Map) {
		if required, ok := node["required"]; ok {
			switch req := required.(type) {
			case []interface{}:
				// Ensure all items are strings
				validItems := []string{}
				hasInvalidItems := false

				for _, item := range req {
					if str, ok := item.(string); ok {
						validItems = append(validItems, str)
					} else {
						hasInvalidItems = true
					}
				}

				if hasInvalidItems {
					node["required"] = validItems
				}
			case []string:
				// Already valid
			case string:
				// Convert single string to array
				node["required"] = []string{req}
			default:
				// Invalid type, remove it
				delete(node, "required")
			}
		}
	}

	// Fix properties named 'required'
	renameRequiredProperty := func(node Map) {
		if props, ok := node.GetMap("properties"); ok {
			if requiredProp, ok := props.Get("required"); ok {
				props.Delete("required")
				props.Set("isRequired", requiredProp)
			}
		}
	}

	// Fix missing responses
	fixEmptyResponses := func(key string, node Map) {
		if strings.HasSuffix(key, ".responses") && len(node) == 0 {
			// Add a default response
			node["200"] = Map{
				"description": "Success response",
				"content": Map{
					"application/json": Map{
						"schema": Map{
							"type":       "object",
							"properties": Map{},
						},
					},
				},
			}
		}
	}

	// Fix oneOf/anyOf issues
	fixSchemaRefs := func(node Map) {
		// If a schema has oneOf but doesn't have $ref, fix it
		if node.Has("oneOf") || node.Has("anyOf") {
			items, isOneOf := node.GetSlice("oneOf")
			if !isOneOf {
				items, _ = node.GetSlice("anyOf")
				isOneOf = false
			}

			// Check if items are valid
			for _, item := range items {
				if itemMap, ok := item.(map[string]interface{}); ok {
					// If it has enum with 0 items, fix it
					if enum, hasEnum := itemMap["enum"]; hasEnum {
						if enumSlice, isSlice := enum.([]interface{}); isSlice && len(enumSlice) == 0 {
							delete(itemMap, "enum")
						}
					}
				}
			}
		}
	}

	// Fix parameter issues
	fixParameters := func(key string, node Map) {
		if strings.Contains(key, "/parameters/") {
			// Ensure required path parameters have required=true
			if in, ok := node["in"].(string); ok && in == "path" {
				node["required"] = true
			}

			// Ensure 'in' is a valid value
			if in, ok := node["in"]; ok {
				inStr, isStr := in.(string)
				if !isStr || (inStr != "path" && inStr != "query" && inStr != "header" && inStr != "cookie") {
					node["in"] = "path" // Default to path
				}
			}
		}
	}

	// Apply all fixes to the entire schema
	applyAllFixes := func(key string, node Map) {
		fixRequiredArray(node)
		renameRequiredProperty(node)
		fixEmptyResponses(key, node)
		fixSchemaRefs(node)
		fixParameters(key, node)
	}

	// Process the entire schema
	for _, pathInfo := range schema.Paths {
		for _, endpoint := range pathInfo.Endpoints {
			endpoint.Iterate(applyAllFixes)
		}
	}
	schema.Components.Iterate(applyAllFixes)
}

func fixItemsTypes(schema *Schema) {
	fixObjectItems := func(key string, node Map) {
		if strings.HasSuffix(key, ".items") {
			// If there's a validation error about items needing to be objects
			if node.Has("type") && node["type"] != "object" {
				// Convert to a proper object schema
				node["type"] = "object"
				if !node.Has("properties") {
					node["properties"] = Map{}
				}
			}
		}
	}

	// Apply the fix
	for _, pathInfo := range schema.Paths {
		for _, endpoint := range pathInfo.Endpoints {
			endpoint.Iterate(fixObjectItems)
		}
	}
	schema.Components.Iterate(fixObjectItems)
}

// fixArrayItemsDefinitions fixes arrays where items is incorrectly defined as an array
func fixArrayItemsDefinitions(schema *Schema) {
	fixArrayItemsInNode := func(key string, node Map) {
		// Check if this node is an array type with an items field
		if nodeType, ok := node["type"].(string); ok && nodeType == "array" {
			if items, ok := node["items"]; ok {
				// If items is an array (incorrect format), fix it
				if itemsArray, isArray := items.([]interface{}); isArray {
					if len(itemsArray) > 0 {
						// Take the first item from the array and use it as the direct items value
						node["items"] = itemsArray[0]
					} else {
						// Empty array, set to a default schema
						node["items"] = Map{"type": "object"}
					}
				}
			}
		}
	}

	// Apply the fix to the entire schema
	for _, pathInfo := range schema.Paths {
		for _, endpoint := range pathInfo.Endpoints {
			endpoint.Iterate(fixArrayItemsInNode)
		}
	}
	schema.Components.Iterate(fixArrayItemsInNode)

}

func fixRequiredWithAnyOf(schema *Schema) {
	// This function fixes schemas where 'anyOf' is incorrectly placed inside 'required'
	fixRequiredInNode := func(key string, node Map) {
		if required, ok := node["required"]; ok {
			if requiredArray, isArray := required.([]interface{}); isArray {
				// Look for anyOf objects within the required array
				var anyOfItem map[string]interface{}
				var anyOfIndex int
				var hasAnyOf bool

				for i, item := range requiredArray {
					if itemMap, isMap := item.(map[string]interface{}); isMap {
						if _, hasAnyOfKey := itemMap["anyOf"]; hasAnyOfKey {
							anyOfItem = itemMap
							anyOfIndex = i
							hasAnyOf = true
							break
						}
					}
				}

				if hasAnyOf {
					// Remove the anyOf item from the required array
					newRequired := append(requiredArray[:anyOfIndex], requiredArray[anyOfIndex+1:]...)
					node["required"] = newRequired

					// Add anyOf at the schema level
					anyOfArray := anyOfItem["anyOf"].([]interface{})
					newAnyOf := make([]interface{}, len(anyOfArray))

					for i, option := range anyOfArray {
						newAnyOf[i] = map[string]interface{}{
							"required": []string{option.(string)},
						}
					}

					node["anyOf"] = newAnyOf
				}
			}
		}
	}

	// Apply the fix to the entire schema, focusing on component schemas
	if schemas, ok := schema.Components.GetMap("schemas"); ok {
		for schemaName, schemaObj := range schemas {
			if schemaMap, ok := schemaObj.(map[string]interface{}); ok {
				fixRequiredInNode(fmt.Sprintf("component.schemas.%s", schemaName), Map(schemaMap))
			}
		}
	}

	// Also check all other schema definitions in the paths
	for _, pathInfo := range schema.Paths {
		for _, endpoint := range pathInfo.Endpoints {
			endpoint.Iterate(fixRequiredInNode)
		}
	}
}

func removeTechnicalPreviewPaths(schema *Schema) {
	// Collect paths to remove
	pathsToRemove := []string{}

	for pathURL, pathInfo := range schema.Paths {
		// Check each HTTP method (GET, POST, PUT, DELETE)
		for _, endpoint := range map[string]Map{
			"get":    pathInfo.Get,
			"post":   pathInfo.Post,
			"put":    pathInfo.Put,
			"delete": pathInfo.Delete,
		} {
			if endpoint == nil {
				continue
			}

			// Check if the endpoint is marked as Technical Preview
			if state, hasState := endpoint["x-state"]; hasState && state == "Technical Preview" {
				pathsToRemove = append(pathsToRemove, pathURL)
				break // No need to check other methods for this path
			}
		}
	}

	// Remove the collected paths
	for _, path := range pathsToRemove {
		numberOfPaths--
		delete(schema.Paths, path)
	}
}

func transformFixVersionFields(schema *Schema) {
	// Process all schemas
	if schemas, ok := schema.Components.GetMap("schemas"); ok {
		for _, schemaObj := range schemas {
			if schemaMap, ok := schemaObj.(Map); ok {
				if props, ok := schemaMap.GetMap("properties"); ok {
					// Check if this schema has both version and _version fields
					if props.Has("version") && props.Has("_version") {
						// Create a new field "versionUnderscored" to replace "_version"
						props["versionUnderscored"] = props["_version"]
						delete(props, "_version")

						// Also update "required" fields if necessary
						if required, ok := schemaMap["required"].([]interface{}); ok {
							for i, field := range required {
								if fieldName, ok := field.(string); ok && fieldName == "_version" {
									required[i] = "versionUnderscored"
								}
							}
						}
					}

					// Also handle @timestamp field
					if props.Has("@timestamp") {
						props["atTimestamp"] = props["@timestamp"]
						delete(props, "@timestamp")

						// Update "required" fields if necessary
						if required, ok := schemaMap["required"].([]interface{}); ok {
							for i, field := range required {
								if fieldName, ok := field.(string); ok && fieldName == "@timestamp" {
									required[i] = "atTimestamp"
								}
							}
						}
					}
				}
			}
		}
	}

	// Update references to versioned schemas in request bodies and responses
	updateReferencesInContent := func(content Map) {
		for _, mediaType := range content {
			if mediaTypeMap, ok := mediaType.(Map); ok {
				if mediaSchema, ok := mediaTypeMap.GetMap("schema"); ok {
					// Handle direct _version properties
					if props, ok := mediaSchema.GetMap("properties"); ok {
						if props.Has("version") && props.Has("_version") {
							props["versionUnderscored"] = props["_version"]
							delete(props, "_version")

							// Also update "required" fields if necessary
							if required, ok := mediaSchema["required"].([]interface{}); ok {
								for i, field := range required {
									if fieldName, ok := field.(string); ok && fieldName == "_version" {
										required[i] = "versionUnderscored"
									}
								}
							}
						}

						// Also handle @timestamp field
						if props.Has("@timestamp") {
							props["atTimestamp"] = props["@timestamp"]
							delete(props, "@timestamp")

							// Update "required" fields if necessary
							if required, ok := mediaSchema["required"].([]interface{}); ok {
								for i, field := range required {
									if fieldName, ok := field.(string); ok && fieldName == "@timestamp" {
										required[i] = "atTimestamp"
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Process request bodies
	if requestBodies, ok := schema.Components.GetMap("requestBodies"); ok {
		for _, requestBody := range requestBodies {
			if requestBodyMap, ok := requestBody.(Map); ok {
				if content, ok := requestBodyMap.GetMap("content"); ok {
					updateReferencesInContent(content)
				}
			}
		}
	}

	// Process responses
	if responses, ok := schema.Components.GetMap("responses"); ok {
		for _, response := range responses {
			if responseMap, ok := response.(Map); ok {
				if content, ok := responseMap.GetMap("content"); ok {
					updateReferencesInContent(content)
				}
			}
		}
	}

	// Process paths
	for _, pathInfo := range schema.Paths {
		pathInfo.Endpoints(func(method string, endpoint Map) bool {
			// Process request body
			if reqBody, ok := endpoint.GetMap("requestBody"); ok {
				if content, ok := reqBody.GetMap("content"); ok {
					updateReferencesInContent(content)
				}
			}

			// Process responses
			if responses, ok := endpoint.GetMap("responses"); ok {
				for _, response := range responses {
					if responseMap, ok := response.(Map); ok {
						if content, ok := responseMap.GetMap("content"); ok {
							updateReferencesInContent(content)
						}
					}
				}
			}

			return true
		})
	}
}

// transformRemoveProblematicPaths removes paths that are problematic with no easy fix
func transformRemoveProblematicPaths(schema *Schema) {
	pathsToRemove := []string{
		"/api/detection_engine/rules/preview",
		"/api/upgrade_assistant/reindex/batch",
		"/api/security/roles",
		"/api/security/role/{name}",
		//"/api/security/role/_query",
		//"/api/security/role",
	}

	for _, path := range pathsToRemove {
		if _, exists := schema.Paths[path]; exists {
			delete(schema.Paths, path)
			numberOfPaths--
		}
	}
}

// transformRunMessageEmail fixes the required field for anyOf fields
func transformRunMessageEmail(schema *Schema) {
	componentSchemas := schema.Components.MustGetMap("schemas")
	runMessagEmailMap, ok := componentSchemas.GetMap("run_message_email")
	if !ok {
		log.Printf("run_message_email schema does not exist")
		return
	}
	required, ok := runMessagEmailMap.GetSlice("required")
	required = required[:len(required)-1]
	runMessagEmailMap.Set("required", required)

	anyOf := []map[string][]string{
		{"required": []string{"to"}},
		{"required": []string{"cc"}},
		{"required": []string{"bcc"}},
	}

	runMessagEmailMap.Set("anyOf", anyOf)
}

// transformEntityAnalyticsTypes removes the format field from the schema as it is incorrect
func transformEntityAnalyticsTypes(schema *Schema) {
	componentSchemas := schema.Components.MustGetMap("schemas")
	entityRiskScoredRecordMap, ok := componentSchemas.GetMap("Security_Entity_Analytics_API_EntityRiskScoreRecord")
	if !ok {
		log.Print("Not exist")
	}
	properties, ok := entityRiskScoredRecordMap.GetMap("properties")
	for k := range properties {
		if k == "category_1_count" || k == "category_2_count" {
			property, _ := properties.GetMap(k)
			property.Delete("format")
			properties[k] = property
		}
	}
}

// transformGetAllSpaces removes the oneOf parameters has it redeclares duplicate types
func transformGetAllSpaces(schema *Schema) {
	space := schema.GetPath("/api/spaces/space")
	getSpace := space.GetEndpoint("get")
	parameters, _ := getSpace.GetSlice("parameters")

	parameter, ok := parameters.GetMap("1")
	if !ok {
	}

	parameterSchema, ok := parameter.GetMap("schema")
	if !ok {
	}

	if parameterSchema.Has("oneOf") {
		parameterSchema.Delete("oneOf")
		parameter.Set("schema", parameterSchema)
		getSpace.Set("1", parameter)
		schema.Paths["/api/spaces/space"].SetEndpoint("get", getSpace)
	}
}
