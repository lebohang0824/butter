package ast

import (
	"fmt"
	"strings"
)

type FilterBlockType string

const (
	FilterEndpoint FilterBlockType = "endpoint"
	FilterFeature  FilterBlockType = "feature"
	FilterListener FilterBlockType = "listener"
)

type Filter struct {
	BlockType FilterBlockType
	Name      string
}

func ParseFilters(raw string) ([]Filter, error) {
	if raw == "" {
		return nil, nil
	}

	var filters []Filter
	parts := strings.Split(raw, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		idx := strings.Index(part, ":")
		if idx == -1 {
			return nil, fmt.Errorf("invalid filter format %q — expected block_type:name (e.g. endpoint:ProcessOrder)", part)
		}

		blockType := FilterBlockType(part[:idx])
		name := strings.TrimSpace(part[idx+1:])

		if name == "" {
			return nil, fmt.Errorf("invalid filter %q — name cannot be empty", part)
		}

		switch blockType {
		case FilterEndpoint, FilterFeature, FilterListener:
		default:
			return nil, fmt.Errorf("invalid block type %q in filter — supported types: endpoint, feature, listener", blockType)
		}

		filters = append(filters, Filter{BlockType: blockType, Name: name})
	}

	return filters, nil
}

func FilterAppSpec(spec *AppSpec, onlyFilters, excludeFilters []Filter) (*AppSpec, error) {
	if len(onlyFilters) == 0 && len(excludeFilters) == 0 {
		return spec, nil
	}

	filtered := &AppSpec{
		App:         spec.App,
		Description: spec.Description,
		Version:     spec.Version,
	}

	excludeIndex := buildFilterIndex(excludeFilters)

	hasOnly := len(onlyFilters) > 0

	if hasOnly {
		for _, f := range onlyFilters {
			matched := false
			switch f.BlockType {
			case FilterFeature:
				for _, feat := range spec.Features {
					if feat.Name == f.Name {
						filtered.Features = append(filtered.Features, feat)
						matched = true
					}
				}
			case FilterEndpoint:
				for _, ep := range spec.Endpoints {
					if ep.Name == f.Name {
						filtered.Endpoints = append(filtered.Endpoints, ep)
						matched = true
					}
				}
			case FilterListener:
				for _, l := range spec.Listeners {
					if l.Name == f.Name {
						filtered.Listeners = append(filtered.Listeners, l)
						matched = true
					}
				}
			}
			if !matched {
				return nil, fmt.Errorf("--only target not found: %s:%s does not exist in spec", f.BlockType, f.Name)
			}
		}
	} else {
		filtered.Features = append(filtered.Features, spec.Features...)
		filtered.Endpoints = append(filtered.Endpoints, spec.Endpoints...)
		filtered.Listeners = append(filtered.Listeners, spec.Listeners...)
	}

	if len(excludeFilters) > 0 {
		for _, f := range excludeFilters {
			matched := false
			switch f.BlockType {
			case FilterFeature:
				for _, feat := range spec.Features {
					if feat.Name == f.Name {
						matched = true
						break
					}
				}
			case FilterEndpoint:
				for _, ep := range spec.Endpoints {
					if ep.Name == f.Name {
						matched = true
						break
					}
				}
			case FilterListener:
				for _, l := range spec.Listeners {
					if l.Name == f.Name {
						matched = true
						break
					}
				}
			}
			if !matched {
				return nil, fmt.Errorf("--exclude target not found: %s:%s does not exist in spec", f.BlockType, f.Name)
			}
		}
		filtered.Features = filterFeatures(filtered.Features, excludeIndex)
		filtered.Endpoints = filterEndpoints(filtered.Endpoints, excludeIndex)
		filtered.Listeners = filterListeners(filtered.Listeners, excludeIndex)
	}

	return filtered, nil
}

type filterIndex struct {
	features  map[string]bool
	endpoints map[string]bool
	listeners map[string]bool
}

func buildFilterIndex(filters []Filter) filterIndex {
	idx := filterIndex{
		features:  make(map[string]bool),
		endpoints: make(map[string]bool),
		listeners: make(map[string]bool),
	}
	for _, f := range filters {
		switch f.BlockType {
		case FilterFeature:
			idx.features[f.Name] = true
		case FilterEndpoint:
			idx.endpoints[f.Name] = true
		case FilterListener:
			idx.listeners[f.Name] = true
		}
	}
	return idx
}

func filterFeatures(features []FeatureSpec, exclude filterIndex) []FeatureSpec {
	var result []FeatureSpec
	for _, feat := range features {
		if !exclude.features[feat.Name] {
			result = append(result, feat)
		}
	}
	return result
}

func filterEndpoints(endpoints []EndpointSpec, exclude filterIndex) []EndpointSpec {
	var result []EndpointSpec
	for _, ep := range endpoints {
		if !exclude.endpoints[ep.Name] {
			result = append(result, ep)
		}
	}
	return result
}

func filterListeners(listeners []ListenerSpec, exclude filterIndex) []ListenerSpec {
	var result []ListenerSpec
	for _, l := range listeners {
		if !exclude.listeners[l.Name] {
			result = append(result, l)
		}
	}
	return result
}
