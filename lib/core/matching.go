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

	paramRegexp := "(.+)"

	routeParamsMatch := routeParamsRegexp.FindAllSubmatchIndex([]byte(m.HandlerPath), -1)

	if len(routeParamsMatch) > 0 {
		routePathRegexp := m.HandlerPath

		routeParamNames := []string{}

		diff := 0

		for _, indexes := range routeParamsMatch {
			name := routePathRegexp[indexes[2]+diff+1 : indexes[3]+diff]

			routeParamNames = append(routeParamNames, name)

			routePathRegexp = routePathRegexp[:indexes[0]+diff] + paramRegexp + routePathRegexp[indexes[1]+diff:]

			diff += len(paramRegexp) - (indexes[3] - indexes[2])
		}

		pathRegexp := regexp.MustCompile(`^` + routePathRegexp + `$`)

		pathMatch := pathRegexp.FindAllSubmatchIndex([]byte(m.RequestedPath), -1)

		if len(pathMatch) > 0 {
			params := Params{}

			for index := 2; index < len(pathMatch[0]); index += 2 {
				matchIndexes := pathMatch[0][index : index+2]
				value := m.RequestedPath[matchIndexes[0]:matchIndexes[1]]

				params[routeParamNames[(index/2)-1]] = value
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
