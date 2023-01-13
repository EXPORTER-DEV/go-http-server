package server

import "regexp"

type Matching struct {
	RequestedPath string
	HandlerPath   string
}

type Params map[string]string

type MatchingExecuteOptions struct {
	ParseParams bool
	IsRegexp    bool
}

func (m *Matching) ParseParams() Params {
	routeParamsRegexp := regexp.MustCompile(`(\:[a-zA-Z]+)`)

	routeParamsMatch := routeParamsRegexp.FindAllSubmatchIndex([]byte(m.HandlerPath), -1)

	if len(routeParamsMatch) > 0 {
		routePathRegexp := m.HandlerPath

		routeParamNames := []string{}

		for _, indexes := range routeParamsMatch {
			name := routePathRegexp[indexes[2]+1 : indexes[3]]

			routeParamNames = append(routeParamNames, name)

			routePathRegexp = routePathRegexp[:indexes[0]] + "(.*)" + routePathRegexp[indexes[1]:]
		}

		pathRegexp := regexp.MustCompile(routePathRegexp)

		pathMatch := pathRegexp.FindAllSubmatchIndex([]byte(m.RequestedPath), -1)

		if len(pathMatch) > 0 {
			params := Params{}

			for index, matchIndexes := range pathMatch {
				value := m.RequestedPath[matchIndexes[2]:matchIndexes[3]]

				params[routeParamNames[index]] = value
			}

			return params
		}
	}

	return nil
}

func (m *Matching) Execute(options MatchingExecuteOptions) (bool, Params) {
	if options.IsRegexp {
		r := regexp.MustCompile(m.HandlerPath)
		if r.Match([]byte(m.RequestedPath)) {
			return true, nil
		}
	}

	if options.ParseParams {
		params := m.ParseParams()

		if params != nil {
			return true, params
		}
	}

	if m.HandlerPath == m.RequestedPath {
		return true, nil
	}

	return false, nil
}
