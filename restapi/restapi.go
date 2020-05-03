package restapi

import (
	"log"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/miun173/rss-reader/helper"
	"github.com/miun173/rss-reader/repository"
	"github.com/valyala/fasthttp"
)

// Server :nodoc:
type Server struct {
	router     *router.Router
	sourceRepo *repository.SourceRepository
}

// NewServer :nodoc:
func NewServer(sr *repository.SourceRepository) *Server {
	s := &Server{
		router:     router.New(),
		sourceRepo: sr,
	}

	return s
}

// Router :nodoc:
func (s *Server) Router() *router.Router {
	s.router.GET("/ping", ping)
	s.router.GET("/sources/{id}", s.findSourceByID)
	s.router.GET("/sources", s.findAllSources)

	return s.router
}

func (s *Server) findSourceByID(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	source, err := s.sourceRepo.FindByID(helper.StringToInt64(id))
	if err != nil {
		log.Println("error : ", err)
		return
	}

	if source == nil {
		writeError(ctx, http.StatusNotFound, "source not found")
		return
	}

	writeOK(ctx, source)
}

func (s *Server) findAllSources(ctx *fasthttp.RequestCtx) {
	sources, err := s.sourceRepo.FindAll(10, 0)
	if err != nil {
		log.Println("error : ", err)
		writeError(ctx, http.StatusNotFound, err.Error())
		return
	}

	writeOK(ctx, sources)
}

func (s *Server) findRSSItemsBySourceID(ctx *fasthttp.RequestCtx) {

}
