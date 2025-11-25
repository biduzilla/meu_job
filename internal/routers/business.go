package routers

import (
	"meu_job/internal/handlers"
	"meu_job/internal/middleware"
	"meu_job/internal/models"

	"github.com/go-chi/chi"
)

type businessRouter struct {
	business handlers.BusinessHandlerInterface
	m        middleware.MiddlewareInterface
}

func NewBusinessRouter(
	business handlers.BusinessHandlerInterface,
	m middleware.MiddlewareInterface,
) *businessRouter {
	return &businessRouter{
		business: business,
		m:        m,
	}
}

func (b *businessRouter) BusinessRoutes(r chi.Router) {
	r.Route("/business", func(r chi.Router) {
		r.Use(b.m.RequireActivatedUser)

		r.Get("/{id}", b.business.FindByID)
		r.Get("/", b.business.FindAll)
		adminOnly := b.m.RequirePermission([]models.Role{models.BUSINESS})

		r.With(adminOnly).Post("/", b.business.Save)
		r.With(adminOnly).Put("/", b.business.Update)
		r.With(adminOnly).Delete("/{id}", b.business.Delete)
	})
}
