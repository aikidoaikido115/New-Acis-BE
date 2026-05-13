package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Field struct {
	Name     string
	Type     string
	GormTag  string
	IsPK     bool
	SkipAttr bool
}

type Entity struct {
	Name     string
	Fields   []Field
	FKFields map[string]bool
}

type Relation struct {
	From  string
	To    string
	Label string
}

func main() {
	inDir := flag.String("in", "modules/entities", "directory with entity files")
	outFile := flag.String("out", "docs/erd.mmd", "output Mermaid file")
	flag.Parse()

	entities, relations, err := buildERD(*inDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := writeMermaid(*outFile, entities, relations); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func buildERD(inDir string) (map[string]*Entity, []Relation, error) {
	fset := token.NewFileSet()
	pkEntities := map[string]bool{}
	allEntities := map[string]*Entity{}

	parsed, err := parser.ParseDir(fset, inDir, func(info os.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	for _, pkg := range parsed {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}
				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}
					entity := &Entity{
						Name:     typeSpec.Name.Name,
						Fields:   []Field{},
						FKFields: map[string]bool{},
					}
					for _, field := range structType.Fields.List {
						if len(field.Names) == 0 {
							continue
						}
						name := field.Names[0].Name
						typeStr := exprString(fset, field.Type)
						gormTag := getTagValue(field.Tag, "gorm")
						isPK := hasPrimaryKey(gormTag)
						skip := strings.Contains(gormTag, "-")
						entity.Fields = append(entity.Fields, Field{
							Name:     name,
							Type:     typeStr,
							GormTag:  gormTag,
							IsPK:     isPK,
							SkipAttr: skip,
						})
						if isPK {
							pkEntities[entity.Name] = true
						}
					}
					allEntities[entity.Name] = entity
				}
			}
		}
	}

	relationsMap := map[string]Relation{}

	for name, entity := range allEntities {
		if !pkEntities[name] {
			continue
		}
		for _, field := range entity.Fields {
			if !strings.Contains(strings.ToLower(field.GormTag), "foreignkey:") {
				continue
			}
			typeName, isSlice, ok := extractTypeName(field.Type)
			if !ok || !pkEntities[typeName] {
				continue
			}

			fkFields := parseTagList(field.GormTag, "foreignKey")
			refFields := parseTagList(field.GormTag, "references")
			label := joinKeyPairs(fkFields, refFields)

			var from, to string
			var fkOwner string
			if isSlice {
				from = name
				to = typeName
				fkOwner = typeName
			} else {
				from = typeName
				to = name
				fkOwner = name
			}

			if fkOwnerEntity, ok := allEntities[fkOwner]; ok {
				for _, fk := range fkFields {
					fkOwnerEntity.FKFields[fk] = true
				}
			}

			key := from + "|" + to + "|" + label
			if _, ok := relationsMap[key]; !ok {
				relationsMap[key] = Relation{From: from, To: to, Label: label}
			}
		}
	}

	relations := make([]Relation, 0, len(relationsMap))
	for _, rel := range relationsMap {
		relations = append(relations, rel)
	}

	sort.Slice(relations, func(i, j int) bool {
		if relations[i].From == relations[j].From {
			if relations[i].To == relations[j].To {
				return relations[i].Label < relations[j].Label
			}
			return relations[i].To < relations[j].To
		}
		return relations[i].From < relations[j].From
	})

	filtered := map[string]*Entity{}
	for name, entity := range allEntities {
		if pkEntities[name] {
			filtered[name] = entity
		}
	}

	return filtered, relations, nil
}

func writeMermaid(outFile string, entities map[string]*Entity, relations []Relation) error {
	if err := os.MkdirAll(filepath.Dir(outFile), 0o755); err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString("erDiagram\n")

	entityNames := make([]string, 0, len(entities))
	for name := range entities {
		entityNames = append(entityNames, name)
	}
	sort.Strings(entityNames)

	for _, name := range entityNames {
		entity := entities[name]
		buf.WriteString("  " + entity.Name + " {\n")
		for _, field := range entity.Fields {
			if field.SkipAttr {
				continue
			}
			if strings.Contains(strings.ToLower(field.GormTag), "foreignkey:") {
				continue
			}
			attrType := mapType(field.Type)
			buf.WriteString("    " + attrType + " " + field.Name)
			var flags []string
			if field.IsPK {
				flags = append(flags, "PK")
			}
			if entity.FKFields[field.Name] {
				flags = append(flags, "FK")
			}
			if len(flags) > 0 {
				buf.WriteString(" " + strings.Join(flags, ", "))
			}
			buf.WriteString("\n")
		}
		buf.WriteString("  }\n")
	}

	for _, rel := range relations {
		label := rel.Label
		if label == "" {
			label = "rel"
		}
		buf.WriteString("  " + rel.From + " ||--o{ " + rel.To + " : \"" + label + "\"\n")
	}

	return os.WriteFile(outFile, buf.Bytes(), 0o644)
}

func exprString(fset *token.FileSet, expr ast.Expr) string {
	var buf bytes.Buffer
	_ = printer.Fprint(&buf, fset, expr)
	return buf.String()
}

func getTagValue(tag *ast.BasicLit, key string) string {
	if tag == nil {
		return ""
	}
	unquoted, err := strconv.Unquote(tag.Value)
	if err != nil {
		return ""
	}
	return reflect.StructTag(unquoted).Get(key)
}

func hasPrimaryKey(gormTag string) bool {
	parts := strings.Split(gormTag, ";")
	for _, part := range parts {
		p := strings.TrimSpace(part)
		if strings.EqualFold(p, "primaryKey") || strings.EqualFold(p, "primarykey") {
			return true
		}
	}
	return false
}

func extractTypeName(typeStr string) (string, bool, bool) {
	trimmed := strings.TrimSpace(typeStr)
	isSlice := false
	for strings.HasPrefix(trimmed, "[]") {
		isSlice = true
		trimmed = strings.TrimPrefix(trimmed, "[]")
	}
	trimmed = strings.TrimPrefix(trimmed, "*")
	if trimmed == "" {
		return "", false, false
	}
	if strings.Contains(trimmed, ".") {
		return "", isSlice, false
	}
	return trimmed, isSlice, true
}

func parseTagList(gormTag, key string) []string {
	parts := strings.Split(gormTag, ";")
	for _, part := range parts {
		p := strings.TrimSpace(part)
		if strings.HasPrefix(strings.ToLower(p), strings.ToLower(key)+":") {
			value := strings.TrimSpace(p[len(key)+1:])
			if value == "" {
				return nil
			}
			items := strings.Split(value, ",")
			for i := range items {
				items[i] = strings.TrimSpace(items[i])
			}
			return items
		}
	}
	return nil
}

func joinKeyPairs(fks, refs []string) string {
	if len(fks) == 0 {
		return ""
	}
	if len(refs) == len(fks) && len(refs) > 0 {
		pairs := make([]string, len(fks))
		for i := range fks {
			pairs[i] = fks[i] + "->" + refs[i]
		}
		return strings.Join(pairs, ",")
	}
	return strings.Join(fks, ",")
}

func mapType(typeStr string) string {
	trimmed := strings.TrimSpace(typeStr)
	trimmed = strings.TrimPrefix(trimmed, "*")
	trimmed = strings.TrimPrefix(trimmed, "[]")
	trimmed = strings.TrimPrefix(trimmed, "*")

	switch trimmed {
	case "string":
		return "string"
	case "bool":
		return "boolean"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "int"
	case "float32", "float64":
		return "float"
	case "time.Time":
		return "datetime"
	case "datatypes.JSON":
		return "json"
	case "[]byte":
		return "blob"
	default:
		return "string"
	}
}
