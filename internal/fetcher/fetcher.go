package fetcher

import "time"

//message is the universal unit of informatin across all sources

type Message struct{
	ID  string
	Source string 
	Subject string
	Author string 
	Date time.Time 
	URL  string
	Body string 
	ThreadID string 
	Tags  []string 


}

//source is the interface every fetcher plugin must implement 
type Source interface{
	Name() string 
	Fetch() ([]Message, error)

}
