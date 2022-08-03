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
	router      *router.Router
	sourceRepo  *repository.SourceRepository
	rssItemRepo *repository.RSSItemRepository
}

// NewServer :nodoc:
func NewServer(sr *repository.SourceRepository,
	rir *repository.RSSItemRepository) *Server {
	s := &Server{
		router:      router.New(),
		sourceRepo:  sr,
		rssItemRepo: rir,
	}

	return s
}

// Router :nodoc:
func (s *Server) Router() *router.Router {
	s.router.GET("/ping", ping)
	s.router.GET("/sources/{id}", s.findSourceByID)
	s.router.GET("/sources", s.findAllSources)
	s.router.GET("/rss/sources/{sourceID}", s.findRSSItemsBySourceID)

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
	args := ctx.QueryArgs()
	size := helper.StringToInt(string(args.Peek("size")))
	page := helper.StringToInt(string(args.Peek("page")))

	sources, err := s.sourceRepo.FindAll(size, page)
	if err != nil {
		log.Println("error : ", err)
		writeError(ctx, http.StatusNotFound, err.Error())
		return
	}

	writeOK(ctx, sources)
}

func (s *Server) findRSSItemsBySourceID(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	size := helper.StringToInt(string(args.Peek("size")))
	page := helper.StringToInt(string(args.Peek("page")))
	sourceID := helper.StringToInt64(ctx.UserValue("sourceID").(string))

	sources, err := s.rssItemRepo.FindBySourceID(sourceID, size, page)
	if err != nil {
		log.Println("error : ", err.Error())
		writeError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	writeOK(ctx, sources)
}
