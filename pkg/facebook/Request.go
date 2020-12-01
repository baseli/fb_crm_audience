package facebook

func buildGraphUrl(uri string, params map[string]string) string {
	if params == nil {
		return graphApiPrefix + graphVersion + uri
	}

	query := ""
	for k, v := range params {
		query += k + "=" + v + "&"
	}

	return graphApiPrefix + graphVersion + uri + "?" + query[:len(query) - 1]
}
