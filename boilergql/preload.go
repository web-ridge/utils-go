package boilergql

import (
	"context"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/99designs/gqlgen/graphql"
)

type ColumnSetting struct {
	Name                  string
	RelationshipModelName string
	IDAvailable           bool // ID is available without preloading
}

func PreloadsContainMoreThanID(a []string, v string) bool {
	for _, av := range a {
		if strings.HasPrefix(av, v) &&
			av != v && // e.g. parentTable
			!strings.HasPrefix(av, v+".id") { // e.g parentTable.id
			return true
		}
	}
	return false
}

func PreloadsContain(a []string, v string) bool {
	for _, av := range a {
		if av == v {
			return true
		}
	}
	return false
}

func GetPreloadMods(ctx context.Context, preloadMap map[string]map[string]ColumnSetting, modelName string) (
	queryMods []qm.QueryMod) {
	return GetPreloadModsWithLevel(ctx, preloadMap, modelName, "")
}

func GetPreloadModsWithLevel(ctx context.Context, preloadMap map[string]map[string]ColumnSetting, modelName string,
	level string) (queryMods []qm.QueryMod) {
	jsonPreloads := GetPreloadsFromContext(ctx, level)
	// e.g. jsonPreloads: [user.organization.id, user.friends.organization]
	dbPreloads := getDatabasePreloads(jsonPreloads, preloadMap, modelName, 0, "")
	for _, dbPreload := range dbPreloads {
		queryMods = append(queryMods, qm.Load(dbPreload, qm.WithDeleted()))
	}
	return
}

func getDatabasePreloads(
	jsonPreloads []string,
	preloadMap map[string]map[string]ColumnSetting,
	modelName string,
	nested int,
	dbPreloadKey string,
) []string {
	// get column settings for current model
	columnSettings, hasColumnSettings := preloadMap[modelName]
	if !hasColumnSettings {
		return nil
	}

	var dbPreloads []string

	for _, jsonPreload := range jsonPreloads {
		// skip .id, .name, .whatever only pick root table for now
		if strings.Count(jsonPreload, ".") > 0 {
			continue
		}

		// get column setting for current preload
		columnSetting, hasColumnSetting := columnSettings[jsonPreload]
		if hasColumnSetting { //nolint:nestif
			dbKey := columnSetting.Name
			if dbPreloadKey != "" {
				dbKey = dbPreloadKey + "." + columnSetting.Name
			}

			// if root table has a foreign key available (inside the table) we don't need to preload the whole table
			// if the user only wanted the id of that table
			if columnSetting.IDAvailable {
				if PreloadsContainMoreThanID(jsonPreloads, jsonPreload) {
					dbPreloads = append(dbPreloads, dbKey)
				}
			} else {
				dbPreloads = append(dbPreloads, dbKey)
			}

			// get nested preloads for this relation
			if columnSetting.RelationshipModelName != "" {
				dbPreloads = append(dbPreloads,
					getDatabasePreloads(
						StripPreloads(jsonPreloads, jsonPreload),
						preloadMap,
						columnSetting.RelationshipModelName,
						nested+1,
						dbKey,
					)...)
			}
		}
	}
	return dbPreloads
}

func GetPreloadsFromContext(ctx context.Context, level string) []string {
	return StripPreloads(GetNestedPreloads(
		graphql.GetOperationContext(ctx),
		graphql.CollectFieldsCtx(ctx, nil),
		"",
	), level)
}

// e.g. sometimes input is deeper and we want
// createdFlowBlock.block.blockChoice => when we fetch block in database we want to strip flowBlock
func StripPreloads(preloads []string, prefix string) []string {
	if prefix == "" {
		return preloads
	}
	var newPreloads []string
	for _, preload := range preloads {
		if strings.HasPrefix(preload, prefix+".") {
			newPreloads = append(newPreloads, strings.TrimPrefix(preload, prefix+"."))
		}
	}
	return newPreloads
}

func GetNestedPreloads(ctx *graphql.OperationContext, fields []graphql.CollectedField, prefix string) (
	preloads []string) {
	for _, column := range fields {
		prefixColumn := GetPreloadString(prefix, column.Name)
		preloads = append(preloads, prefixColumn)
		preloads = append(preloads, GetNestedPreloads(ctx,
			graphql.CollectFields(ctx, column.Selections, nil), prefixColumn)...)
	}
	return
}

func GetPreloadString(prefix, name string) string {
	if len(prefix) > 0 {
		return prefix + "." + name
	}
	return name
}
