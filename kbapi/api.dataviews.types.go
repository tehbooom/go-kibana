package kbapi

type DataViewsObject struct {
	// AllowNoIndex Allows the data view saved object to exist before the data is available.
	AllowNoIndex *bool                           `json:"allowNoIndex,omitempty"`
	FieldAttrs   *map[string]DataViewsFieldAttrs `json:"fieldAttrs,omitempty"`

	// FieldFormats A map of field formats by field name.
	FieldFormats *map[string]interface{}     `json:"fieldFormats,omitempty"`
	Fields       *map[string]DataViewsFields `json:"fields,omitempty"`
	ID           *string                     `json:"id,omitempty"`

	// Name The data view name.
	Name *string `json:"name,omitempty"`

	// Namespaces An array of space identifiers for sharing the data view between multiple spaces.
	Namespaces      *[]string                            `json:"namespaces,omitempty"`
	RuntimeFieldMap *map[string]DataViewsRuntimeFieldMap `json:"runtimeFieldMap,omitempty"`

	// SourceFilters The array of field names you want to filter out in Discover.
	SourceFilters *[]DataViewsSourceFilter `json:"sourceFilters,omitempty"`

	// TimeFieldName The timestamp field name, which you use for time-based data views.
	TimeFieldName *string `json:"timeFieldName,omitempty"`

	// Title Comma-separated list of data streams, indices, and aliases that you want to search. Supports wildcards (`*`).
	Title string `json:"title"`

	// Type When set to `rollup`, identifies the rollup data views.
	Type *string `json:"type,omitempty"`

	// TypeMeta When you use rollup indices, contains the field list for the rollup data view API endpoints.
	TypeMeta *DataViewsTypeMeta `json:"typeMeta,omitempty"`
	Version  *string            `json:"version,omitempty"`
}

// DataViewsFieldattrs A map of field attributes by field name.
type DataViewsFieldAttrs struct {
	// Count Popularity count for the field.
	Count *int `json:"count,omitempty"`

	// CustomDescription Custom description for the field.
	CustomDescription *string `json:"customDescription,omitempty"`

	// CustomLabel Custom label for the field.
	CustomLabel *string `json:"customLabel,omitempty"`
}

// DataViewsRuntimefieldmap A map of runtime field definitions by field name.
type DataViewsRuntimeFieldMap struct {
	Script struct {
		// Source Script for the runtime field.
		Source *string `json:"source,omitempty"`
	} `json:"script"`

	// Type Mapping type of the runtime field.
	Type string `json:"type"`
}

// DataViewsSourcefilters The array of field names you want to filter out in Discover.
type DataViewsSourceFilter struct {
	Value string `json:"value"`
}

// DataViewsTypemeta When you use rollup indices, contains the field list for the rollup data view API endpoints.
type DataViewsTypeMeta struct {
	// Aggs A map of rollup restrictions by aggregation type and field name.
	Aggs map[string]interface{} `json:"aggs"`

	// Params Properties for retrieving rollup fields.
	Params map[string]interface{} `json:"params"`
}

type DataViewsFieldFormats struct {
	ID     string `json:"id"`
	Params struct {
		Pattern string `json:"pattern"`
	} `json:"params"`
}

type DataViewsFields struct {
	Count             int      `json:"count"`
	Name              string   `json:"name"`
	Type              string   `json:"type"`
	ESTypes           []string `json:"esTypes"`
	Scripted          bool     `json:"scripted"`
	Searchable        bool     `json:"searchable"`
	Aggregatable      bool     `json:"aggregatable"`
	ReadFromDocValues bool     `json:"readFromDocValues"`
	Format            struct {
		ID string `json:"id"`
	}
	ShortDotsEnable bool                     `json:"shortDotsEnable"`
	RuntimeField    DataViewsRuntimeFieldMap `json:"runtimeField,omitempty"`
}
