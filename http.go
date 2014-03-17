package main

func handlerFuncWith(stages ...func(w *http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, f := range stages {
			if err := f(&w, r); err != nil {
				break
			}
		}
	}
}

func statsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseWriter.Write([]byte(`OK`))
	}
}
