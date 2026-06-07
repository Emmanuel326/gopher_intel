package fetcher 

var sources = map[string]Source{}


func Register (s Source) {
	sources[s.Name()]=s 
	
}

func All() map[string]Source{
	return sources
	
}

